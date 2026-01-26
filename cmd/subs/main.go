package main

import (
	"context"
	_ "embed"
	"net/http"
	"os"
	"time"

	common "github.com/end1essrage/efmob-tz/pkg/common/cmd"
	l "github.com/end1essrage/efmob-tz/pkg/common/logger"
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
	//создаем инстанс логгера
	logger := l.New(os.Getenv("SERVICE_NAME"), true, common.ENV(os.Getenv("ENV")) != common.ENV_PROD).Log("main", "main")

	// проверяем энвы
	if os.Getenv("ENV") == "" {
		logger.Fatal("ENV is empty")
	}
	if os.Getenv("PORT") == "" {
		logger.Fatal("PORT is empty")
	}
	if os.Getenv("SERVICE_NAME") == "" {
		logger.Fatal("SERVICE_NAME is empty")
	}

	logger.Infof("запуск %s сервиса", os.Getenv("SERVICE_NAME"))

	//корневой контекст
	ctx := common.Context()

	//создаем роутер
	r, cleanup := createSubsMicroservice(ctx)
	//создаем http сервер
	server := &http.Server{
		Addr:              ":" + os.Getenv("PORT"),
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
	logger.Infof("остановка %s серёвиса", os.Getenv("SERVICE_NAME"))

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("http shutdown error: %v", err)
	}

	cleanup()
}

func createSubsMicroservice(ctx context.Context) (*chi.Mux, func()) {
	logger := l.Logger().Log("main", "main")

	/*
		publicKey, err := crypto.LoadPublicKeyFromPEM(os.Getenv("JWT_PUBLIC_KEY_PATH"))
		if err != nil {
			logger.Fatal(err)
		}
	*/

	// бд
	//usersRepo := memory.NewUserRepository()
	//tokenRepo := memory.NewTokenRepository()
	logger.Info("бд инициализирована")

	// helpers

	/*
		// клиент для доступа к блэклисту
		kvClient, err := kv.NewClient(kv.Config{
			Addr:     os.Getenv("TOKEN_BLACKLIST_ADDR"),
			Password: os.Getenv("TOKEN_BLACKLIST_PWD"),
			DB:       0,

			DialTimeout:  3 * time.Second,
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
		})
		if err != nil {
			logger.Warn("KV redis unavailable", err)
		}

		// проверятель блэклиста
		tokenBlackListChecker := common_blacklist.NewRedisTokenBlacklistChecker(kvClient)

		// валидатор токенов
		tokenValidator := common_jwt.NewJwtTokenValidator(
			publicKey,
			os.Getenv("JWT_ISSUER"),
			tokenBlackListChecker)

		profileRepo := memory.NewProfileRepository()

	*/

	/*
		//контейнер бизнес логики
		container := container.NewContainer(
			profileRepo,
		)
		logger.Info("di контейнер собран")
	*/
	/*
		// консьюмер событий
		broker, err := broker.NewRedisStreamClient(broker.Config{
			Addr:     os.Getenv("BROKER_ADDR"),
			Password: os.Getenv("BROKER_PWD"),
			DB:       0, //?

			DialTimeout:  3 * time.Second,
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
		})
		if err != nil {
			logger.Warn("BROKER redis unavailable", err)
		}

		event_subscriber := events.NewRedisUserEventSubscriber(broker, container.CreateProfileFromAuthHandler)

		go func(consumer string) {
			if err := event_subscriber.Start(ctx, consumer); err != nil {
				logger.Error("event subscriber stopped", err)
			}
		}(os.Getenv("SERVICE_NAME"))
	*/

	//создаем хендлер
	h := subs_http.NewSubsHandler(
		common.ENV(os.Getenv("ENV")),
	)
	logger.Info("хендлеры инициализированы")

	//заполняем роуты
	r := common.CreateRouter()

	subs_http.AddRoutes(r, h)
	logger.Info("роуты созданы")

	cleanup := func() {
		//очищаем ресурсы

	}
	// подписываемся на прерывание контекста

	return r, cleanup
}
