package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	id          uuid.UUID
	userID      uuid.UUID
	serviceName string
	price       int
	startDate   time.Time
	endDate     *time.Time
	createdAt   time.Time
	updatedAt   time.Time

	version int // оптимистичная блокировка
}

func NewSubscriptionWithVersion(
	id uuid.UUID,
	userID uuid.UUID,
	serviceName string,
	price int,
	startDate time.Time,
	endDate *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
	version int,
) (*Subscription, error) {
	// Проверяем только бизнес-правила
	if userID == uuid.Nil {
		return nil, errors.New("user id is required")
	}
	if strings.TrimSpace(serviceName) == "" {
		return nil, ErrInvalidServiceName
	}
	if price <= 0 {
		return nil, ErrInvalidPrice
	}

	startDate = normalizeMonth(startDate)
	if endDate != nil {
		normalizedEnd := normalizeMonth(*endDate)
		if !normalizedEnd.After(startDate) {
			return nil, ErrInvalidDates
		}
		endDate = &normalizedEnd
	}

	return &Subscription{
		id:          id,
		userID:      userID,
		serviceName: serviceName,
		price:       price,
		startDate:   startDate,
		endDate:     endDate,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		version:     version,
	}, nil
}

func NewSubscription(
	id uuid.UUID,
	userID uuid.UUID,
	serviceName string,
	price int,
	startDate time.Time,
	endDate *time.Time,
) (*Subscription, error) {
	if id == uuid.Nil {
		id = uuid.New()
	}

	if userID == uuid.Nil {
		return nil, errors.New("user id is required")
	}

	if strings.TrimSpace(serviceName) == "" {
		return nil, ErrInvalidServiceName
	}

	if price <= 0 {
		return nil, ErrInvalidPrice
	}

	startDate = normalizeMonth(startDate)

	if endDate != nil {
		normalizedEnd := normalizeMonth(*endDate)
		if !normalizedEnd.After(startDate) {
			return nil, ErrInvalidDates
		}
		endDate = &normalizedEnd
	}

	now := time.Now()

	return &Subscription{
		id:          id,
		userID:      userID,
		serviceName: serviceName,
		price:       price,
		startDate:   startDate,
		endDate:     endDate,
		createdAt:   now,
		updatedAt:   now,
		version:     1, // Начальная версия
	}, nil
}

func (s Subscription) ID() uuid.UUID     { return s.id }
func (s Subscription) UserID() uuid.UUID { return s.userID }

func (s Subscription) ServiceName() string { return s.serviceName }
func (s Subscription) Price() int          { return s.price }

func (s Subscription) StartDate() time.Time { return s.startDate }
func (s Subscription) EndDate() *time.Time  { return s.endDate }
func (s Subscription) CreatedAt() time.Time { return s.createdAt }
func (s Subscription) UpdatedAt() time.Time { return s.updatedAt }

func (s Subscription) Version() int { return s.version }
func (s Subscription) IsActive(at time.Time) bool {
	if at.Before(s.startDate) {
		return false
	}
	if s.endDate == nil {
		return true
	}
	return at.Before(*s.endDate) || at.Equal(*s.endDate)
}

func (s *Subscription) ChangePrice(price int) error {
	if price <= 0 {
		return ErrInvalidPrice
	}

	s.price = price
	s.updatedAt = time.Now()
	return nil
}

func (s *Subscription) ChangeStartDate(start time.Time) {
	s.startDate = normalizeMonth(start)
	s.updatedAt = time.Now()
}

func (s *Subscription) ChangeEndDate(end time.Time) {
	e := normalizeMonth(end)
	s.endDate = &e
	s.updatedAt = time.Now()
}

func (s *Subscription) NilEndDate() {
	s.endDate = nil
	s.updatedAt = time.Now()
}

func normalizeMonth(t time.Time) time.Time {
	return time.Date(
		t.Year(),
		t.Month(),
		1,
		0, 0, 0, 0,
		time.UTC,
	)
}
