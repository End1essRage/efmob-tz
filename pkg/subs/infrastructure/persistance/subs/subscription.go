package subs

import (
	"time"

	domain "github.com/end1essrage/efmob-tz/pkg/subs/domain"
	"github.com/google/uuid"
)

type SubscriptionModel struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index"`
	ServiceName string     `gorm:"type:text;not null;index"`
	Price       int        `gorm:"not null"`
	StartDate   time.Time  `gorm:"not null;index"`
	EndDate     *time.Time `gorm:"index"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`

	Version int `gorm:"not null;default:1"`
}

func (SubscriptionModel) TableName() string {
	return "subscriptions"
}

func (m *SubscriptionModel) ToDomain() *domain.Subscription {
	sub, err := domain.NewSubscriptionWithVersion(
		m.ID,
		m.UserID,
		m.ServiceName,
		m.Price,
		m.StartDate,
		m.EndDate,
		m.CreatedAt,
		m.UpdatedAt,
		m.Version,
	)
	if err != nil {
		// Лучше вернуть ошибку, но для совместимости:
		return nil
	}
	return sub
}

// Конвертация из домена
func FromDomain(sub *domain.Subscription) *SubscriptionModel {
	return &SubscriptionModel{
		ID:          sub.ID(),
		UserID:      sub.UserID(),
		ServiceName: sub.ServiceName(),
		Price:       sub.Price(),
		StartDate:   sub.StartDate(),
		EndDate:     sub.EndDate(),
		CreatedAt:   sub.CreatedAt(),
		UpdatedAt:   sub.UpdatedAt(),
		Version:     sub.Version(),
	}
}
