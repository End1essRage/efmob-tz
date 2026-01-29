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
	repo domain.SubscriptionRepository
}

func NewCreateSubscriptionHandler(repo domain.SubscriptionRepository) *CreateSubscriptionHandler {
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

	if _, err := h.repo.Create(ctx, sub); err != nil {
		log.Errorf("creating error: %v", err)
		return nil, err
	}

	subsMetrics.SubscriptionsCreatedTotal.Inc()

	return sub, nil
}
