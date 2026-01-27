package domain

import "errors"

var (
	ErrInvalidServiceName   = errors.New("service name cannot be empty")
	ErrInvalidPrice         = errors.New("price must be positive")
	ErrInvalidDates         = errors.New("end date must be after start date")
	ErrInvalidDateFormat    = errors.New("invalid date format, should be mm-yyyy")
	ErrInvalidPeriod        = errors.New("invalid period")
	ErrSubscriptionNotFound = errors.New("subscription not found")
)
