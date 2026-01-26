package queries

import (
	"context"

	domain "github.com/end1essrage/efmob-tz/pkg/subs/domain"
	"github.com/google/uuid"
)

type GetSubscriptionQuery struct {
	ID uuid.UUID
}

type GetSubscriptionHandler struct {
	repo domain.SubscriptionRepository
}

func NewGetSubscriptionHandler(repo domain.SubscriptionRepository) *GetSubscriptionHandler {
	return &GetSubscriptionHandler{repo: repo}
}

func (h *GetSubscriptionHandler) Handle(ctx context.Context, q GetSubscriptionQuery) (*domain.Subscription, error) {
	return h.repo.GetByID(ctx, q.ID)
}
