package container

import (
	cmd "github.com/end1essrage/efmob-tz/pkg/subs/application/commands"
	quer "github.com/end1essrage/efmob-tz/pkg/subs/application/queries"
	domain "github.com/end1essrage/efmob-tz/pkg/subs/domain"
)

type Container struct {
	CreateSubscriptionHandler *cmd.CreateSubscriptionHandler
	UpdateSubscriptionHandler *cmd.UpdateSubscriptionHandler
	DeleteSubscriptionHandler *cmd.DeleteSubscriptionHandler

	GetSubscriptionHandler   *quer.GetSubscriptionHandler
	ListSubscriptionsHandler *quer.ListSubscriptionsHandler
	TotalCostHandler         *quer.TotalCostHandler
}

func NewContainer(
	subRepo domain.SubscriptionRepository, // для queries
	subRepoTx domain.SubscriptionRepositoryWithTx, // для commands с транзакциями
	statsRepo domain.SubscriptionStatsRepository,
) *Container {
	return &Container{
		CreateSubscriptionHandler: cmd.NewCreateSubscriptionHandler(subRepoTx),
		UpdateSubscriptionHandler: cmd.NewUpdateSubscriptionHandler(subRepo),
		DeleteSubscriptionHandler: cmd.NewDeleteSubscriptionHandler(subRepoTx),

		GetSubscriptionHandler:   quer.NewGetSubscriptionHandler(subRepo),
		ListSubscriptionsHandler: quer.NewListSubscriptionsHandler(subRepo),
		TotalCostHandler:         quer.NewTotalCostHandler(statsRepo),
	}
}
