package application

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/end1essrage/efmob-tz/pkg/subs/domain"
	infra "github.com/end1essrage/efmob-tz/pkg/subs/infrastructure/persistance/subs"
)

type ErrorValidationCommand struct {
	Msg string
}

func NewErrorValidationCommand(msg string) *ErrorValidationCommand {
	return &ErrorValidationCommand{Msg: msg}
}
func (e ErrorValidationCommand) Error() string {
	return fmt.Sprintf("недоступное действие: %s", e.Msg)
}

type ErrorValidationQuery struct {
	Msg string
}

func NewErrorValidationQuery(msg string) *ErrorValidationQuery {
	return &ErrorValidationQuery{Msg: msg}
}
func (e ErrorValidationQuery) Error() string {
	return fmt.Sprintf("ошибка в qury: %s", e.Msg)
}

type AppError struct {
	Err        error
	HTTPStatus int
	Code       string
}

// маппер ошибок
func MapError(err error) *AppError {
	var valCmdErr *ErrorValidationCommand
	if errors.As(err, &valCmdErr) {
		return &AppError{
			Err:        err,
			HTTPStatus: http.StatusBadRequest,
			Code:       "INVALID_COMMAND",
		}
	}

	var valQueryErr *ErrorValidationQuery
	if errors.As(err, &valQueryErr) {
		return &AppError{
			Err:        err,
			HTTPStatus: http.StatusBadRequest,
			Code:       "INVALID_QUERY",
		}
	}

	var retriesExcErr *infra.ErrorRetriesExceeded
	if errors.As(err, &retriesExcErr) {
		return &AppError{
			Err:        err,
			HTTPStatus: http.StatusInternalServerError,
			Code:       "INTERNAL_ERROR",
		}
	}

	switch {
	// Infra errors
	case errors.Is(err, infra.ErrConcurrentModification):
		return &AppError{Err: err, HTTPStatus: http.StatusConflict, Code: "CONCURRENT_MODIFICATION"}
	case errors.Is(err, infra.ErrInvalidSortingField):
		return &AppError{Err: err, HTTPStatus: http.StatusBadRequest, Code: "INVALID_SORTING_FIELD"}
	// Subscription domain errors
	case errors.Is(err, domain.ErrInvalidServiceName):
		return &AppError{Err: err, HTTPStatus: http.StatusBadRequest, Code: "INVALID_SERVICE_NAME"}
	case errors.Is(err, domain.ErrInvalidPrice):
		return &AppError{Err: err, HTTPStatus: http.StatusBadRequest, Code: "INVALID_PRICE"}
	case errors.Is(err, domain.ErrInvalidDates):
		return &AppError{Err: err, HTTPStatus: http.StatusBadRequest, Code: "INVALID_DATES"}
	case errors.Is(err, domain.ErrInvalidPeriod):
		return &AppError{Err: err, HTTPStatus: http.StatusBadRequest, Code: "INVALID_PERIOD"}
	case errors.Is(err, domain.ErrSubscriptionNotFound):
		return &AppError{Err: err, HTTPStatus: http.StatusNotFound, Code: "NOT_FOUND"}
	// Default - 500 Internal Server Error
	default:
		return &AppError{Err: err, HTTPStatus: http.StatusInternalServerError, Code: "INTERNAL_ERROR"}
	}
}
