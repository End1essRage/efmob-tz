package middleware

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/end1essrage/efmob-tz/pkg/common/interfaces/http/utils"
	l "github.com/end1essrage/efmob-tz/pkg/common/logger"
)

type RateLimiter struct {
	requests chan time.Time
	rate     time.Duration
	burst    int
	stop     chan struct{}
	running  atomic.Bool
}

func NewRateLimiter(rate time.Duration, limit int, burst int) *RateLimiter {
	if burst <= 0 {
		burst = limit
	}
	if burst > limit*2 {
		burst = limit * 2
	}

	rl := &RateLimiter{
		requests: make(chan time.Time, burst),
		rate:     rate / time.Duration(limit), // интервал между токенами
		burst:    burst,
		stop:     make(chan struct{}),
	}

	// Заполняем начальными токенами
	for i := 0; i < burst; i++ {
		rl.requests <- time.Now()
	}

	// Запускаем пополнение токенов
	rl.running.Store(true)
	go rl.refill()

	return rl
}

func (rl *RateLimiter) refill() {
	ticker := time.NewTicker(rl.rate)
	defer ticker.Stop()

	for {
		select {
		case <-rl.stop:
			rl.running.Store(false)
			return
		case <-ticker.C:
			select {
			case rl.requests <- time.Now():
				// Токен добавлен
			default:
				// Bucket полный
			}
		}
	}
}

func (rl *RateLimiter) Stop() {
	if rl.running.Load() {
		close(rl.stop)
	}
}

func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.requests:
		return true
	default:
		return false
	}
}

func RateLimitMiddleware(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !rl.Allow() {
				l.Logger().Log("middleware", "rate_limiter").Warn("Rate limit exceeded")

				// Добавляем заголовки согласно RFC 6585
				w.Header().Set("Retry-After", "60")
				utils.WriteJSON(w, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later.")
				return
			}

			// Добавляем информацию о лимитах в заголовки
			w.Header().Set("X-RateLimit-Limit", "100")
			w.Header().Set("X-RateLimit-Remaining", "99") // TODO: вычислять реальное значение
			w.Header().Set("X-RateLimit-Reset", "60")

			next.ServeHTTP(w, r)
		})
	}
}
