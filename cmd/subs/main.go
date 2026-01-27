package main

import (
	"context"
	_ "embed"
	"net/http"
	"os"
	"time"

	common "github.com/end1essrage/efmob-tz/pkg/common/cmd"
	l "github.com/end1essrage/efmob-tz/pkg/common/logger"
	"github.com/end1essrage/efmob-tz/pkg/subs/application/container"
	"github.com/end1essrage/efmob-tz/pkg/subs/domain"
	subs_repo "github.com/end1essrage/efmob-tz/pkg/subs/infrastructure/persistance/subs"
	subs_http "github.com/end1essrage/efmob-tz/pkg/subs/interfaces/http"
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

	//корневой контекст
	ctx := common.Context()

	//создаем роутер
	r, cleanup := createSubsMicroservice(cfg)
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

func createSubsMicroservice(cfg *Config) (*chi.Mux, func()) {
	log := l.Logger().Log("main", "createSubsMicroservice")

	// бд
	var repo domain.SubscriptionRepository
	var statsRepo domain.SubscriptionStatsRepository
	var cleanupDB func()

	if common.ENV(cfg.Env) == common.ENV_DEV {
		memRepo := subs_repo.NewInMemorySubscriptionRepo()
		repo = memRepo
		statsRepo = memRepo
	} else {
		dsn := cfg.PostgresDSN
		gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent), // отключаем логирование
		})
		if err != nil {
			log.Fatalf("failed to connect to postgres: %v", err)
		}

		// только в ТЕСТ
		if common.ENV(cfg.Env) == common.ENV_TEST {
			// Авто-миграция таблицы subscriptions
			err = gormDB.AutoMigrate(&subs_repo.SubscriptionModel{})
			if err != nil {
				log.Fatalf("failed to auto-migrate subscriptions: %v", err)
			}
		}

		pgRepo := subs_repo.NewGormSubscriptionRepo(gormDB)
		repo = pgRepo
		statsRepo = pgRepo

		// Если нужно закрывать соединение при shutdown
		sqlDB, err := gormDB.DB()
		if err == nil {
			cleanupDB = func() {
				_ = sqlDB.Close()
			}
		}
	}

	di := container.NewContainer(repo, statsRepo)

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

	cleanup := func() {
		log.Info("очистка зависимостей")
		if cleanupDB != nil {
			cleanupDB()
		}
	}

	return r, cleanup
}
