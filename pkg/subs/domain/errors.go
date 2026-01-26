package domain

import "errors"

var (
	ErrInvalidServiceName = errors.New("service name cannot be empty")
	ErrInvalidPrice       = errors.New("price must be positive")
	ErrInvalidDates       = errors.New("end date must be after start date")
	ErrInvalidPeriod      = errors.New("invalid period")
)
