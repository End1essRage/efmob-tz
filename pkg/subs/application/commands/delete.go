package commands

import (
	"context"

	domain "github.com/end1essrage/efmob-tz/pkg/subs/domain"
	"github.com/google/uuid"
)

type DeleteSubscriptionCommand struct {
	ID uuid.UUID
}

type DeleteSubscriptionHandler struct {
	repo domain.SubscriptionRepository
}

func NewDeleteSubscriptionHandler(repo domain.SubscriptionRepository) *DeleteSubscriptionHandler {
	return &DeleteSubscriptionHandler{repo: repo}
}

func (h *DeleteSubscriptionHandler) Handle(ctx context.Context, cmd DeleteSubscriptionCommand) error {
	return h.repo.Delete(ctx, cmd.ID)
}
