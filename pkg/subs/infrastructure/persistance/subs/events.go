package subs

import (
	"context"
	"time"

	domain "github.com/end1essrage/efmob-tz/pkg/subs/domain"
	"github.com/google/uuid"
)

type EventModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Type      string    `gorm:"type:text;not null"`
	Payload   []byte
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (r *GormSubscriptionRepo) CreateEvent(ctx context.Context, event domain.Event) error {
	payload, err := event.MarshalJSON()
	if err != nil {
		return err
	}

	model := EventModel{
		ID:        uuid.New(),
		Type:      event.Type(),
		Payload:   payload,
		CreatedAt: time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	return nil
}
