package cmd

import (
	"net/http"
	"time"

	"github.com/end1essrage/efmob-tz/pkg/common/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

var MiddlewareLogger func(http.Handler) http.Handler = func(next http.Handler) http.Handler {
	l := logger.Logger().Log("router", "middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем обертку для ResponseWriter для получения статуса
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Выполняем следующий обработчик
		next.ServeHTTP(ww, r)

		// Логируем информацию о запросе
		duration := time.Since(start)

		fields := make(logrus.Fields)
		fields["router"] = "middleware"
		fields["method"] = r.Method
		fields["path"] = r.URL.Path
		fields["remote_addr"] = r.RemoteAddr
		fields["status"] = ww.Status()
		fields["duration"] = duration.String()
		fields["duration_ms"] = duration.Milliseconds()
		fields["user_agent"] = r.UserAgent()
		fields["time"] = start.Format(time.RFC3339)

		l.WithFields(fields).Info()
	})
}

// TODO metrics middleware
func CreateRouter() *chi.Mux {
	r := chi.NewRouter()

	// порядок важен
	r.Use(MiddlewareLogger)

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
