package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func HTTPMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		routePattern := "unknown"
		if rc := chi.RouteContext(r.Context()); rc != nil {
			routePattern = rc.RoutePattern()
		}

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(ww.Status())

		HTTPRequestsTotal.WithLabelValues(
			r.Method,
			routePattern,
			status,
		).Inc()

		HTTPRequestDuration.WithLabelValues(
			r.Method,
			routePattern,
		).Observe(duration)
	})
}
