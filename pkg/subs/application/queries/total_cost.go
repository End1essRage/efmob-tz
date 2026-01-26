package queries

import (
	"context"

	domain "github.com/end1essrage/efmob-tz/pkg/subs/domain"
)

type TotalCostQuery struct {
	Query domain.SubscriptionQuery
}

type TotalCostHandler struct {
	repo domain.SubscriptionStatsRepository
}

func NewTotalCostHandler(repo domain.SubscriptionStatsRepository) *TotalCostHandler {
	return &TotalCostHandler{repo: repo}
}

func (h *TotalCostHandler) Handle(ctx context.Context, q TotalCostQuery) (int, error) {
	return h.repo.CalculateTotal(ctx, q.Query)
}
