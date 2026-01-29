package commands

import (
	"context"
	"time"

	"github.com/end1essrage/efmob-tz/pkg/common/logger"
	domain "github.com/end1essrage/efmob-tz/pkg/subs/domain"
	subsMetrics "github.com/end1essrage/efmob-tz/pkg/subs/metrics"
	"github.com/google/uuid"
)

type CreateSubscriptionCommand struct {
	UserID      uuid.UUID
	ServiceName string
	Price       int
	StartDate   time.Time
	EndDate     *time.Time
}

type CreateSubscriptionHandler struct {
	repo domain.SubscriptionRepositoryWithTx
}

func NewCreateSubscriptionHandler(repo domain.SubscriptionRepositoryWithTx) *CreateSubscriptionHandler {
	return &CreateSubscriptionHandler{repo: repo}
}

func (h *CreateSubscriptionHandler) Handle(ctx context.Context, cmd CreateSubscriptionCommand) (*domain.Subscription, error) {
	log := logger.Logger().WithFields(logger.LogOptions{
		Pkg:  "CreateSubscriptionHandler",
		Func: "Handle",
		Ctx:  ctx,
	})

	sub, err := domain.NewSubscription(
		uuid.Nil,
		cmd.UserID,
		cmd.ServiceName,
		cmd.Price,
		cmd.StartDate,
		cmd.EndDate,
	)
	if err != nil {
		log.Errorf("entity validation error: %v", err)
		return nil, err
	}

	err = h.repo.RunInTransaction(ctx, func(tx domain.TxSubscriptionRepository) error {
		// создаём подписку
		uid, err := tx.Create(ctx, sub)
		if err != nil {
			return err
		}

		// создаем событие
		event := domain.SubCreatedEvent{
			Id:     uid,
			UserID: sub.UserID(),
		}
		if err := tx.CreateEvent(ctx, event); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Errorf("creating error: %v", err)
		return nil, err
	}

	subsMetrics.SubscriptionsCreatedTotal.Inc()

	return sub, nil
}
