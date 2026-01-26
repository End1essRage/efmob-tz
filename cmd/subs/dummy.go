package main

import (
	"context"

	p "github.com/end1essrage/efmob-tz/pkg/common/persistance"
	"github.com/end1essrage/efmob-tz/pkg/subs/domain"
	"github.com/google/uuid"
)

type DummyRepo struct{}

func (*DummyRepo) Create(ctx context.Context, sub *domain.Subscription) error {
	return nil
}
func (*DummyRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	return nil, nil
}
func (*DummyRepo) Update(ctx context.Context, sub *domain.Subscription) error {
	return nil
}
func (*DummyRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (*DummyRepo) Find(ctx context.Context, q domain.SubscriptionQuery, p p.Pagination, s p.Sorting) ([]*domain.Subscription, error) {
	return nil, nil
}

func (*DummyRepo) CalculateTotal(ctx context.Context, q domain.SubscriptionQuery) (int, error) {
	return 0, nil
}
