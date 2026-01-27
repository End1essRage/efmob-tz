package http

import (
	"github.com/go-chi/chi/v5"

	common "github.com/end1essrage/efmob-tz/pkg/common/cmd"
	"github.com/end1essrage/efmob-tz/pkg/subs/application/container"
)

type SubsHandler struct {
	container *container.Container
	env       common.ENV
}

func NewSubsHandler(env common.ENV, container *container.Container) *SubsHandler {
	return &SubsHandler{
		env:       env,
		container: container,
	}
}

// AddRoutes добавляет маршруты аутентификации
// @Summary Add subscriptions routes
func AddRoutes(r *chi.Mux, h *SubsHandler) {
	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", h.CreateSubscription)
		r.Get("/", h.ListSubscriptions)
		r.Get("/total", h.GetTotalCost)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetSubscription)
			r.Patch("/", h.UpdateSubscription)
			r.Delete("/", h.DeleteSubscription)
		})
	})
}
