package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	SubscriptionsCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "efmob",
			Subsystem: "subs",
			Name:      "created_total",
			Help:      "Total created subscriptions",
		},
	)
)

func Register() {
	prometheus.MustRegister(
		SubscriptionsCreatedTotal,
	)
}
