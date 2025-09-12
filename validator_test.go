package fault_test

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/marcelofabianov/fault"
)

type RequestBody struct {
	Name  string `validate:"required"`
	Email string `validate:"email"`
	Age   int    `validate:"gte=18"`
}

func TestNewValidationErrorFromValidator(t *testing.T) {
	validate := validator.New()

	t.Run("Should convert multiple validator errors into a structured fault.Error", func(t *testing.T) {
		invalidReq := RequestBody{
			Name:  "",
			Email: "invalid-email",
			Age:   17,
		}

		errs := validate.Struct(invalidReq)
		require.NotNil(t, errs)

		faultErr := fault.NewValidationErrorFromValidator(errs.(validator.ValidationErrors))

		require.NotNil(t, faultErr)
		assert.Equal(t, fault.Invalid, faultErr.Code)
		assert.Equal(t, "Request validation failed", faultErr.Message)
		assert.ErrorIs(t, faultErr, fault.ErrValidation) // Mudança aqui
		require.Len(t, faultErr.Details, 3)

		detail1 := faultErr.Details[0]
		assert.Equal(t, fault.Invalid, detail1.Code)
		assert.Contains(t, detail1.Message, "validation failed on field 'Name'")
		assert.Equal(t, "Name", detail1.Context["field"])
		assert.Equal(t, "required", detail1.Context["tag"])

		detail2 := faultErr.Details[1]
		assert.Equal(t, fault.Invalid, detail2.Code)
		assert.Contains(t, detail2.Message, "validation failed on field 'Email'")
		assert.Equal(t, "Email", detail2.Context["field"])
		assert.Equal(t, "email", detail2.Context["tag"])

		detail3 := faultErr.Details[2]
		assert.Equal(t, fault.Invalid, detail3.Code)
		assert.Contains(t, detail3.Message, "validation failed on field 'Age'")
		assert.Equal(t, "Age", detail3.Context["field"])
		assert.Equal(t, "gte", detail3.Context["tag"])
	})
}
