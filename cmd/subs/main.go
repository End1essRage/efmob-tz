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
	subs_http "github.com/end1essrage/efmob-tz/pkg/subs/interfaces/http"
	"github.com/go-chi/chi/v5"

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

	logger.Infof("запуск %s сервиса", cfg.Env)

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
	logger := l.Logger().Log("main", "main")

	// бд
	repo := &DummyRepo{}

	logger.Info("бд инициализирована")

	//контейнер бизнес логики
	container := container.NewContainer(repo, repo)
	logger.Info("di контейнер собран")

	//создаем хендлер
	h := subs_http.NewSubsHandler(
		common.ENV(os.Getenv("ENV")),
		container,
	)
	logger.Info("хендлеры инициализированы")

	//заполняем роуты
	r := common.CreateRouter()

	subs_http.AddRoutes(r, h)
	logger.Info("роуты созданы")

	cleanup := func() {
		//очищаем ресурсы
	}

	return r, cleanup
}
