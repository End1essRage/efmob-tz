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
	// проверяем что в периодах заполнены обе границы
	if (q.StartFrom != nil && q.StartTo == nil) || (q.StartFrom == nil && q.StartTo != nil) {
		return nil, application.NewErrorValidationQuery("в фильтрах по периоду должны быть заполнены обе границы")
	}
	if (q.EndFrom != nil && q.EndTo == nil) || (q.EndFrom == nil && q.EndTo != nil) {
		return nil, application.NewErrorValidationQuery("в фильтрах по периоду должны быть заполнены обе границы")
	}

	startPeriod, endPeriod, err := application.Periods(q.StartFrom, q.StartTo, q.EndFrom, q.EndTo)
	if err != nil {
		return nil, err
	}

	query := domain.NewSubscriptionQuery(q.UserID, q.ServiceName, startPeriod, endPeriod)

	return h.repo.Find(ctx, query, q.Pagination, q.Sorting)
}
