package http

import "github.com/end1essrage/efmob-tz/pkg/subs/domain"

var mapSubscriptionFromDomain = func(record *domain.Subscription) *Subscription {
	return &Subscription{
		ID:          record.ID(),
		UserID:      record.UserID(),
		ServiceName: record.ServiceName(),
		Price:       record.Price(),
		StartDate:   formatDate(record.StartDate()),
		EndDate:     formatOptionalDate(record.EndDate()),
	}
}
