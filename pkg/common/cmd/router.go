package cmd

import (
	"net/http"
	"time"

	m "github.com/end1essrage/efmob-tz/pkg/common/interfaces/http/middleware"
	"github.com/end1essrage/efmob-tz/pkg/common/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

var MiddlewareLogger = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Получаем RequestID из контекста
		requestID := middleware.GetReqID(r.Context())

		// Создаем обертку для ResponseWriter
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Выполняем запрос
		next.ServeHTTP(ww, r)

		// Логируем
		duration := time.Since(start)

		logger.Logger().Log("router", "middleware").WithFields(logrus.Fields{
			"request_id":  requestID,
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
			"status":      ww.Status(),
			"duration":    duration.String(),
			"duration_ms": duration.Milliseconds(),
			"user_agent":  r.UserAgent(),
			"time":        start.Format(time.RFC3339),
			"bytes":       ww.BytesWritten(),
		}).Info("HTTP request")
	})
}

// TODO metrics middleware
func CreateRouter() *chi.Mux {
	r := chi.NewRouter()

	// порядок важен
	// 1. Recovery middleware (обработка паник)
	r.Use(middleware.Recoverer)

	// 2. Request ID middleware
	r.Use(middleware.RequestID)

	// 3. Rate limiter (например, 100 запросов в минуту с burst 10)
	rateLimiter := m.NewRateLimiter(time.Minute, 100, 10)
	r.Use(m.RateLimitMiddleware(rateLimiter))

	// 4. Logging middleware
	r.Use(MiddlewareLogger)

	// 5. Timeout middleware (таймаут 30 секунд на запрос)
	r.Use(middleware.Timeout(30 * time.Second))

	// swagger docs
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Metrics
	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	return r
}
