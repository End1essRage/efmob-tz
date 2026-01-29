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
	repo domain.SubscriptionRepositoryWithTx
}

func NewDeleteSubscriptionHandler(repo domain.SubscriptionRepositoryWithTx) *DeleteSubscriptionHandler {
	return &DeleteSubscriptionHandler{repo: repo}
}

func (h *DeleteSubscriptionHandler) Handle(ctx context.Context, cmd DeleteSubscriptionCommand) error {
	log := logger.Logger().WithFields(logger.LogOptions{
		Pkg:  "DeleteSubscriptionHandler",
		Func: "Handle",
		Ctx:  ctx,
	})

	err := h.repo.RunInTransaction(ctx, func(tx domain.TxSubscriptionRepository) error {
		// создаём подписку
		if err := tx.Delete(ctx, cmd.ID); err != nil {
			return err
		}

		// создаем событие
		event := domain.SubDeletedEvent{
			Id: cmd.ID,
		}
		if err := tx.CreateEvent(ctx, event); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Errorf("deleting error: %v", err)
	}
	return err
}
