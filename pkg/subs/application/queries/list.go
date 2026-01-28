package queries

import (
	"context"
	"time"

	p "github.com/end1essrage/efmob-tz/pkg/common/persistance"
	"github.com/end1essrage/efmob-tz/pkg/subs/application"
	domain "github.com/end1essrage/efmob-tz/pkg/subs/domain"
	"github.com/google/uuid"
)

type ListSubscriptionsQuery struct {
	UserID      *uuid.UUID
	ServiceName *string
	StartFrom   *time.Time
	StartTo     *time.Time
	EndFrom     *time.Time
	EndTo       *time.Time
	WithNilEnd  *bool

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
	startPeriod, endPeriod, err := application.Periods(q.StartFrom, q.StartTo, q.EndFrom, q.EndTo)
	if err != nil {
		return nil, err
	}

	query := domain.NewSubscriptionQuery(q.UserID, q.ServiceName, startPeriod, endPeriod, q.WithNilEnd)

	return h.repo.Find(ctx, query, q.Pagination, q.Sorting)
}
