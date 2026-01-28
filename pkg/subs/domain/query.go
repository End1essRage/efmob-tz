package domain

import (
	"time"

	"github.com/google/uuid"
)

type Period struct {
	from *time.Time
	to   *time.Time
}

func NewPeriod(from, to *time.Time) (*Period, error) {
	if from != nil && to != nil {
		if to.Before(*from) {
			return nil, ErrInvalidPeriod
		}
	}

	return &Period{
		from: from,
		to:   to,
	}, nil
}

func (p Period) From() *time.Time { return p.from }
func (p Period) To() *time.Time   { return p.to }

type SubscriptionQuery struct {
	userID      *uuid.UUID
	serviceName *string
	startPeriod *Period
	endPeriod   *Period
	endIsNil    *bool
}

func NewSubscriptionQuery(
	userID *uuid.UUID,
	serviceName *string,
	startPeriod *Period,
	endPeriod *Period,
	endIsNil *bool,
) SubscriptionQuery {
	return SubscriptionQuery{
		userID:      userID,
		serviceName: serviceName,
		startPeriod: startPeriod,
		endPeriod:   endPeriod,
		endIsNil:    endIsNil,
	}
}

func (q SubscriptionQuery) UserID() *uuid.UUID   { return q.userID }
func (q SubscriptionQuery) ServiceName() *string { return q.serviceName }
func (q SubscriptionQuery) StartPeriod() *Period { return q.startPeriod }
func (q SubscriptionQuery) EndPeriod() *Period   { return q.endPeriod }
func (q SubscriptionQuery) EndIsNil() *bool      { return q.endIsNil }
