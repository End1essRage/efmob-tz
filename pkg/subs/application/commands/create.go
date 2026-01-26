package commands

import (
	"context"
	"time"

	domain "github.com/end1essrage/efmob-tz/pkg/subs/domain"
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
	sub, err := domain.NewSubscription(
		uuid.Nil,
		cmd.UserID,
		cmd.ServiceName,
		cmd.Price,
		cmd.StartDate,
		cmd.EndDate,
	)
	if err != nil {
		// ошибка валидации при создании записи
		return nil, err
	}

	//TODO retry - на уровне реализации
	if _, err := h.repo.Create(ctx, sub); err != nil {
		// ошибка на стороне репозитория
		return nil, err
	}
	return sub, nil
}
