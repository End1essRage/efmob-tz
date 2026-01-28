package commands

import (
	"context"
	"time"

	"github.com/end1essrage/efmob-tz/pkg/common/logger"
	"github.com/end1essrage/efmob-tz/pkg/subs/application"
	domain "github.com/end1essrage/efmob-tz/pkg/subs/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UpdateSubscriptionCommand struct {
	ID             uuid.UUID
	Price          *int
	StartDate      *time.Time
	EndDate        *time.Time
	SetEndDateNull bool
}

type UpdateSubscriptionHandler struct {
	repo domain.SubscriptionRepository
}

func NewUpdateSubscriptionHandler(repo domain.SubscriptionRepository) *UpdateSubscriptionHandler {
	return &UpdateSubscriptionHandler{repo: repo}
}

// бизнес валидация
func (h *UpdateSubscriptionHandler) Validate(sub *domain.Subscription, cmd UpdateSubscriptionCommand) error {
	// если меняется цена, то нельзя сдвигать срок начала на дату раньше
	if cmd.Price != nil {
		if cmd.StartDate != nil && cmd.StartDate.Before(sub.StartDate()) {
			return application.NewErrorValidationCommand("если меняется цена, то нельзя сдвигать срок начала на дату раньше")
		}
	}
	return nil
}

func (h *UpdateSubscriptionHandler) Handle(ctx context.Context, cmd UpdateSubscriptionCommand) (*domain.Subscription, error) {
	log := logger.Logger().WithFields(logger.LogOptions{
		Pkg:  "UpdateSubscriptionHandler",
		Func: "Handle",
		Ctx:  ctx,
	})

	sub, err := h.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		log.Errorf("getting error: %v", err)
		return nil, err
	}

	// обогащаем лог
	log = log.WithField("entity_id", sub.ID())

	log.Info("запись успешно найдена")

	if err := h.Validate(sub, cmd); err != nil {
		log.Errorf("validation error: %v", err)
		return nil, err
	}

	if cmd.Price != nil {
		oldPrice := sub.Price()

		if err := sub.ChangePrice(*cmd.Price); err != nil {
			log.Errorf("price changing error: %v", err)
			return nil, err
		}

		log.WithFields(logrus.Fields{
			"updated": []logrus.Fields{
				{
					"field_name": "price",
					"old_value":  oldPrice,
					"new_value":  sub.Price(),
				},
			},
		},
		).Info("цена изменилась")
	}

	if cmd.StartDate != nil {
		log.Info("дата начала изменена")
		sub.ChangeStartDate(*cmd.StartDate)
	}

	if cmd.EndDate != nil {
		log.Info("дата окончания изменена")
		sub.ChangeEndDate(*cmd.EndDate)
	} else {
		if cmd.SetEndDateNull {
			log.Info("дата окончания обнулена")
			sub.NilEndDate()
		}
	}

	if err := h.repo.Update(ctx, sub); err != nil {
		log.Errorf("updating error: %v", err)
		return nil, err
	}

	return sub, nil
}
