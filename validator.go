package fault

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func NewValidationErrorFromValidator(errs validator.ValidationErrors) *Error {
	details := make([]*Error, 0, len(errs))

	for _, fieldErr := range errs {
		details = append(details, New(
			fmt.Sprintf("validation failed on field '%s'", fieldErr.Field()),
			WithCode(Invalid),
			WithContext("field", fieldErr.Field()),
			WithContext("tag", fieldErr.Tag()),
			WithContext("param", fieldErr.Param()),
		))
	}

	return New(
		"Request validation failed",
		WithWrappedErr(ErrValidation),
		WithCode(Invalid),
		WithDetails(details...),
	)
}
