package queries

import (
	"context"
	"time"

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
}

type TotalCostHandler struct {
	repo domain.SubscriptionStatsRepository
}

func NewTotalCostHandler(repo domain.SubscriptionStatsRepository) *TotalCostHandler {
	return &TotalCostHandler{repo: repo}
}

func (h *TotalCostHandler) Handle(ctx context.Context, q TotalCostQuery) (int, error) {
	// проверяем что в периодах заполнены обе границы
	if (q.StartFrom != nil && q.StartTo == nil) || (q.StartFrom == nil && q.StartTo != nil) {
		return 0, application.NewErrorValidationQuery("в фильтрах по периоду должны быть заполнены обе границы")
	}
	if (q.EndFrom != nil && q.EndTo == nil) || (q.EndFrom == nil && q.EndTo != nil) {
		return 0, application.NewErrorValidationQuery("в фильтрах по периоду должны быть заполнены обе границы")
	}

	startPeriod, endPeriod, err := application.Periods(q.StartFrom, q.StartTo, q.EndFrom, q.EndTo)
	if err != nil {
		return 0, err
	}

	query := domain.NewSubscriptionQuery(q.UserID, q.ServiceName, startPeriod, endPeriod)

	return h.repo.CalculateTotal(ctx, query)
}
