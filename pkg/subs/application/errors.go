package application

import (
	"errors"
	"net/http"
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
	/*
		// ErrTokenNotFound, ErrTokenInvalid, ErrTokenRevoked  - 401
		case errors.Is(err, domain.ErrTokenInvalid):
			return &AppError{Err: err, HTTPStatus: http.StatusUnauthorized, Code: "UNAUTHORIZED"}
		case errors.Is(err, domain.ErrTokenRevoked):
			return &AppError{Err: err, HTTPStatus: http.StatusUnauthorized, Code: "UNAUTHORIZED"}
		case errors.Is(err, domain.ErrTokenNotFound):
			return &AppError{Err: err, HTTPStatus: http.StatusUnauthorized, Code: "UNAUTHORIZED"}

		// ErrTokenExpired - 401
		case errors.Is(err, domain.ErrTokenExpired):
			return &AppError{Err: err, HTTPStatus: http.StatusUnauthorized, Code: "TOKEN_EXPIRED"}
		// ErrPasswordReuse - 400
		case errors.Is(err, ErrPasswordReuse):
			return &AppError{Err: err, HTTPStatus: http.StatusBadRequest, Code: "PASSWORD_REUSE"}
		// ErrInvalidPassword(validation) - 400
		case errors.Is(err, domain.ErrInvalidPassword):
			return &AppError{Err: err, HTTPStatus: http.StatusBadRequest, Code: "BAD_PASSWORD"}
		// ErrInvalidEmail(validation) - 400
		case errors.Is(err, domain.ErrInvalidEmail):
			return &AppError{Err: err, HTTPStatus: http.StatusBadRequest, Code: "BAD_EMAIL"}
		// ErrInvalidCredentials(login/password) - 401
		case errors.Is(err, domain.ErrInvalidCredentials):
			return &AppError{Err: err, HTTPStatus: http.StatusUnauthorized, Code: "INVALID_CREDENTIALS"}
		// ErrUserNotFound - 404
		case errors.Is(err, domain.ErrUserNotFound):
			return &AppError{Err: err, HTTPStatus: http.StatusNotFound, Code: "USER_NOT_FOUND"}
		// ErrEmailNotVerified - 403
		case errors.Is(err, domain.ErrEmailNotVerified):
			return &AppError{Err: err, HTTPStatus: http.StatusForbidden, Code: "EMAIL_NOT_VERIFIED"}
		// ErrUserAlreadyExists - 409
		case errors.Is(err, domain.ErrUserAlreadyExists):
			return &AppError{Err: err, HTTPStatus: http.StatusConflict, Code: "USER_ALREADY_EXISTS"}
	*/
	// Other - 500
	default:
		return &AppError{Err: err, HTTPStatus: http.StatusInternalServerError, Code: "INTERNAL_ERROR"}
	}
}
