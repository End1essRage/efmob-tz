package application

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/end1essrage/efmob-tz/pkg/subs/domain"
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

	switch {
	// Subscription domain errors
	//400
	case errors.Is(err, domain.ErrInvalidServiceName):
		return &AppError{Err: err, HTTPStatus: http.StatusBadRequest, Code: "INVALID_SERVICE_NAME"}
	case errors.Is(err, domain.ErrInvalidPrice):
		return &AppError{Err: err, HTTPStatus: http.StatusBadRequest, Code: "INVALID_PRICE"}
	case errors.Is(err, domain.ErrInvalidDates):
		return &AppError{Err: err, HTTPStatus: http.StatusBadRequest, Code: "INVALID_DATES"}
	case errors.Is(err, domain.ErrInvalidPeriod):
		return &AppError{Err: err, HTTPStatus: http.StatusBadRequest, Code: "INVALID_PERIOD"}
	//404
	case errors.Is(err, domain.ErrSubscriptionNotFound):
		return &AppError{Err: err, HTTPStatus: http.StatusNotFound, Code: "NOT_FOUND"}
	// Default - 500 Internal Server Error
	default:
		return &AppError{Err: err, HTTPStatus: http.StatusInternalServerError, Code: "INTERNAL_ERROR"}
	}
}
