package queries

import (
	"context"
	"time"

	"github.com/end1essrage/efmob-tz/pkg/common/logger"
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

	Pagination p.Pagination
	Sorting    *p.Sorting
}

type ListSubscriptionsHandler struct {
	repo domain.SubscriptionRepository
}

func NewListSubscriptionsHandler(repo domain.SubscriptionRepository) *ListSubscriptionsHandler {
	return &ListSubscriptionsHandler{repo: repo}
}

func (h *ListSubscriptionsHandler) Handle(ctx context.Context, q ListSubscriptionsQuery) ([]*domain.Subscription, error) {
	log := logger.Logger().WithFields(logger.LogOptions{
		Pkg:  "ListSubscriptionsHandler",
		Func: "Handle",
		Ctx:  ctx,
	})

	startPeriod, endPeriod, err := application.Periods(q.StartFrom, q.StartTo, q.EndFrom, q.EndTo)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	query := domain.NewSubscriptionQuery(q.UserID, q.ServiceName, startPeriod, endPeriod, q.WithNilEnd)

	r, err := h.repo.Find(ctx, query, q.Pagination, q.Sorting)
	if err != nil {
		log.Error(err)
	}

	return r, err
}
