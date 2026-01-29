package queries

import (
	"context"
	"time"

	"github.com/end1essrage/efmob-tz/pkg/common/logger"
	"github.com/end1essrage/efmob-tz/pkg/subs/application"
	domain "github.com/end1essrage/efmob-tz/pkg/subs/domain"
	"github.com/google/uuid"
)

type TotalCostQuery struct {
	UserID      *uuid.UUID
	ServiceName *string
	StartFrom   *time.Time
	StartTo     *time.Time
	EndFrom     *time.Time
	EndTo       *time.Time
	WithNilEnd  *bool
}

type TotalCostHandler struct {
	repo domain.SubscriptionStatsRepository
}

func NewTotalCostHandler(repo domain.SubscriptionStatsRepository) *TotalCostHandler {
	return &TotalCostHandler{repo: repo}
}

func (h *TotalCostHandler) Handle(ctx context.Context, q TotalCostQuery) (int, error) {
	log := logger.Logger().WithFields(logger.LogOptions{
		Pkg:  "TotalCostHandler",
		Func: "Handle",
		Ctx:  ctx,
	})

	startPeriod, endPeriod, err := application.Periods(q.StartFrom, q.StartTo, q.EndFrom, q.EndTo)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	query := domain.NewSubscriptionQuery(q.UserID, q.ServiceName, startPeriod, endPeriod, q.WithNilEnd)

	r, err := h.repo.CalculateTotalCost(ctx, query)
	if err != nil {
		log.Error(err)
	}

	return r, err
}
