package queries

import (
	"context"

	p "github.com/end1essrage/efmob-tz/pkg/common/persistance"
	domain "github.com/end1essrage/efmob-tz/pkg/subs/domain"
)

type ListSubscriptionsQuery struct {
	Query      domain.SubscriptionQuery
	Pagination *p.Pagination
	Sorting    *p.Sorting
}

type ListSubscriptionsHandler struct {
	repo domain.SubscriptionRepository
}

func NewListSubscriptionsHandler(repo domain.SubscriptionRepository) *ListSubscriptionsHandler {
	return &ListSubscriptionsHandler{repo: repo}
}

func (h *ListSubscriptionsHandler) Handle(ctx context.Context, q ListSubscriptionsQuery) ([]*domain.Subscription, error) {
	return h.repo.Find(ctx, q.Query, q.Pagination, q.Sorting)
}
