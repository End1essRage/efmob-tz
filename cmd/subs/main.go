package main

import (
	"context"
	_ "embed"
	"net/http"
	"os"
	"sync"
	"time"

	common "github.com/end1essrage/efmob-tz/pkg/common/cmd"
	l "github.com/end1essrage/efmob-tz/pkg/common/logger"
	common_metrics "github.com/end1essrage/efmob-tz/pkg/common/metrics"
	common_test "github.com/end1essrage/efmob-tz/pkg/common/testing"
	"github.com/end1essrage/efmob-tz/pkg/subs/application/container"
	subs_repo "github.com/end1essrage/efmob-tz/pkg/subs/infrastructure/persistance/subs"
	"github.com/end1essrage/efmob-tz/pkg/subs/infrastructure/publisher"
	subs_http "github.com/end1essrage/efmob-tz/pkg/subs/interfaces/http"
	subs_metrics "github.com/end1essrage/efmob-tz/pkg/subs/metrics"
	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/end1essrage/efmob-tz/cmd/subs/docs"
)

// main.go
// @title Subs Service API
// @description This is the subscriptions service API
// @version 1.0.0
// @host localhost
// @BasePath /subs
// @schemes http
// @tag.name subs
// @tag.description Subscriptions control
func main() {
	// зaгружаем энвы
	cfg := LoadConfig()

	//создаем инстанс логгера
	logger := l.New(cfg.ServiceName, true, common.ENV(cfg.Env) != common.ENV_PROD).Log("main", "main")

	logger.Infof("запуск %s сервиса", cfg.ServiceName)

	// регистрация метрик
	common_metrics.Register()
	subs_metrics.Register()

	//корневой контекст
	ctx := common.Context()

	//создаем роутер
	r, cleanup := createSubsMicroservice(ctx, cfg)
	//создаем http сервер
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second}

	//запускаем http сервер
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatalf("http server error: %v", err)
		}
	}()

	//ждем сигнала остановки
	<-ctx.Done()
	logger.Infof("остановка %s серёвиса", cfg.ServiceName)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("http shutdown error: %v", err)
	}

	// очищаем ресурсы
	cleanup()
}

func createSubsMicroservice(ctx context.Context, cfg *Config) (*chi.Mux, func()) {
	log := l.Logger().Log("main", "createSubsMicroservice")

	// бд
	cleanupStack := make([]func(), 0)

	pushCleanup := func(fn func()) {
		cleanupStack = append(cleanupStack, fn)
	}

	popAllCleanup := func() {
		// чистим в обратном порядке
		for i := len(cleanupStack) - 1; i >= 0; i-- {
			cleanupStack[i]()
		}
	}

	dsn := cfg.PostgresDSN
	// DEV - создаем тест контейнер
	if common.ENV(cfg.Env) == common.ENV_DEV {
		container, cs, err := common_test.SetupPostgresContainer(ctx)
		if err != nil {
			log.Fatalf("failed to start postgres container: %v", err)
		}

		dsn = cs

		// регистрируем очистку контейнера
		pushCleanup(func() {
			if err := container.Terminate(context.Background()); err != nil {
				log.Errorf("ошибка остановки постгрес контейнера: %v", err)
			}
		})
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // отключаем логирование
	})
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	pgRepo := subs_repo.NewGormSubscriptionRepo(gormDB)

	// выключаем миграцию в проде
	if common.ENV(cfg.Env) != common.ENV_PROD {
		// Авто-миграция
		if err := pgRepo.Migrate(); err != nil {
			log.Fatalf("failed to auto-migrate: %v", err)
		}
	}

	sqlDB, err := gormDB.DB()
	if err == nil {
		// регистрируем очистку подключения
		pushCleanup(func() {
			_ = sqlDB.Close()
		})
	}

	di := container.NewContainer(pgRepo, pgRepo, pgRepo)

	log.Info("di контейнер собран")

	//создаем хендлер
	h := subs_http.NewSubsHandler(
		common.ENV(os.Getenv("ENV")),
		di,
	)
	log.Info("хендлеры инициализированы")

	//заполняем роуты
	r := common.CreateRouter()

	subs_http.AddRoutes(r, h)
	log.Info("роуты созданы")

	// создаем и запускаем EventWorker
	publisher := publisher.NewMockPublisher()
	worker := subs_repo.NewEventWorker(gormDB, publisher, 5*time.Second, 100)

	workerCtx, workerCancel := context.WithCancel(ctx)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		worker.Run(workerCtx)
	}()

	pushCleanup(func() {
		workerCancel()
		wg.Wait()
	})

	log.Info("EventWorker запущен")

	return r, popAllCleanup
}
