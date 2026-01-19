package utils

import (
	"errors"
	"net/http"
)

// AppError represents a domain level error with HTTP mapping.
type AppError struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	HTTPStatus int         `json:"-"`
	Err        error       `json:"-"`
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// Unwrap exposes the wrapped error for errors.Is/As.
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError builds a new application error.
func NewAppError(httpStatus, code int, message string, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Err:        err,
	}
}

// Predefined error helpers.
var (
	ErrBadRequest   = NewAppError(http.StatusBadRequest, 400001, "Bad Request", nil)
	ErrUnauthorized = NewAppError(http.StatusUnauthorized, 401001, "Unauthorized", nil)
	ErrForbidden    = NewAppError(http.StatusForbidden, 403001, "Forbidden", nil)
	ErrNotFound     = NewAppError(http.StatusNotFound, 404001, "Resource Not Found", nil)
	ErrInternal     = NewAppError(http.StatusInternalServerError, 500001, "Internal Server Error", nil)
)

// Clone clones predefined errors while overriding details.
func Clone(base *AppError, details interface{}, err error) *AppError {
	if base == nil {
		return &AppError{HTTPStatus: http.StatusInternalServerError, Code: 500001, Message: "Internal Server Error", Err: err}
	}
	return &AppError{
		Code:       base.Code,
		Message:    base.Message,
		Details:    details,
		HTTPStatus: base.HTTPStatus,
		Err:        err,
	}
}

// IsAppError determines whether the error chain contains AppError.
func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
