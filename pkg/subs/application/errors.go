package application

import (
	"errors"
	"net/http"

	"github.com/end1essrage/efmob-tz/pkg/subs/domain"
)

var (
	ErrPasswordReuse = errors.New("password_reuse")
)

type AppError struct {
	Err        error
	HTTPStatus int
	Code       string
}

// маппер ошибок
func MapDomainError(err error) *AppError {
	switch {
	// Subscription domain errors
	case errors.Is(err, domain.ErrInvalidServiceName):
		return &AppError{Err: err, HTTPStatus: http.StatusBadRequest, Code: "INVALID_SERVICE_NAME"}
	case errors.Is(err, domain.ErrInvalidPrice):
		return &AppError{Err: err, HTTPStatus: http.StatusBadRequest, Code: "INVALID_PRICE"}
	case errors.Is(err, domain.ErrInvalidDates):
		return &AppError{Err: err, HTTPStatus: http.StatusBadRequest, Code: "INVALID_DATES"}
	case errors.Is(err, domain.ErrInvalidPeriod):
		return &AppError{Err: err, HTTPStatus: http.StatusBadRequest, Code: "INVALID_PERIOD"}
	// Default - 500 Internal Server Error
	default:
		return &AppError{Err: err, HTTPStatus: http.StatusInternalServerError, Code: "INTERNAL_ERROR"}
	}

}
