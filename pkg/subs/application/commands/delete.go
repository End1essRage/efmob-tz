package commands

import (
	"context"

	"github.com/end1essrage/efmob-tz/pkg/common/logger"
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
	log := logger.Logger().WithFields(logger.LogOptions{
		Pkg:  "DeleteSubscriptionHandler",
		Func: "Handle",
		Ctx:  ctx,
	})

	err := h.repo.Delete(ctx, cmd.ID)
	if err != nil {
		log.Errorf("deleting error: %v", err)
	}
	return err
}
