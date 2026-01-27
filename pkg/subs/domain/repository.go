package domain

import (
	"context"

	p "github.com/end1essrage/efmob-tz/pkg/common/persistance"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *Subscription) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Subscription, error)
	Update(ctx context.Context, sub *Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	Find(ctx context.Context, q SubscriptionQuery, p *p.Pagination, s *p.Sorting) ([]*Subscription, error)
}
type SubscriptionStatsRepository interface {
	CalculateTotal(ctx context.Context, q SubscriptionQuery) (int, error)
}
