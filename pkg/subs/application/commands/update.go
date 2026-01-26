package commands

import (
	"context"
	"time"

	domain "github.com/end1essrage/efmob-tz/pkg/subs/domain"
	"github.com/google/uuid"
)

type UpdateSubscriptionCommand struct {
	ID        uuid.UUID
	Price     int
	StartDate time.Time
	EndDate   *time.Time
}

type UpdateSubscriptionHandler struct {
	repo domain.SubscriptionRepository
}

func NewUpdateSubscriptionHandler(repo domain.SubscriptionRepository) *UpdateSubscriptionHandler {
	return &UpdateSubscriptionHandler{repo: repo}
}

func (h *UpdateSubscriptionHandler) Handle(ctx context.Context, cmd UpdateSubscriptionCommand) (*domain.Subscription, error) {
	sub, err := h.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		// ошибка поиска
		// либо нету, либо сервис не доступен
		return nil, err
	}

	if err := sub.ChangePrice(cmd.Price); err != nil {
		// ошибка валидации команды
		return nil, err
	}
	if err := sub.ChangePeriod(cmd.StartDate, cmd.EndDate); err != nil {
		// ошибка валидации команды
		return nil, err
	}

	if err := h.repo.Update(ctx, sub); err != nil {
		// ошибка сохранения - на стороне сервиса
		return nil, err
	}

	return sub, nil
}
