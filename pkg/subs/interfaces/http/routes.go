package http

import (
	"github.com/go-chi/chi/v5"

	common "github.com/end1essrage/efmob-tz/pkg/common/cmd"
)

type SubsHandler struct {
	//container      *container.Container
	//authMiddleware *middleware.AuthMiddleware
	env common.ENV
}

func NewSubsHandler(env common.ENV) *SubsHandler {
	return &SubsHandler{
		env: env,
	}
}

// AddRoutes добавляет маршруты аутентификации
// @Summary Add subscriptions routes
func AddRoutes(r *chi.Mux, h *SubsHandler) {
	// Регистрация
	r.Get("/test", h.Test)
}
