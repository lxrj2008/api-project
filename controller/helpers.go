package controller

import (
	"github.com/go-playground/validator/v10"

	"github.com/example/go-api/utils"
)

// parseValidationErrors builds a readable map from validator errors.
func parseValidationErrors(err error) map[string]string {
	if err == nil {
		return nil
	}

	details := make(map[string]string)
	switch verr := err.(type) {
	case validator.ValidationErrors:
		for _, fe := range verr {
			field := fe.Field()
			details[field] = fe.Error()
		}
	default:
		details["error"] = err.Error()
	}

	return details
}

// NewBindingError converts validation errors to AppError.
func NewBindingError(err error) *utils.AppError {
	return utils.Clone(utils.ErrBadRequest, parseValidationErrors(err), err)
}
