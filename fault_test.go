package fault

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestError_CreationAndOptions(t *testing.T) {
	originalErr := errors.New("original error")

	t.Run("Create error with all options", func(t *testing.T) {
		detail := New("detail message")
		err := New(
			"user message",
			WithWrappedErr(originalErr),
			WithCode(Invalid),
			WithContext("key1", "value1"),
			WithDetails(detail),
		)

		require.NotNil(t, err)
		assert.Equal(t, "user message", err.Message)
		assert.Equal(t, originalErr, err.Unwrap())
		assert.Equal(t, Invalid, err.Code)
		assert.Equal(t, map[string]any{"key1": "value1"}, err.Context)
		require.Len(t, err.Details, 1)
		assert.Equal(t, detail, err.Details[0])
	})

	t.Run("Create error with Wrap helper", func(t *testing.T) {
		err := Wrap(originalErr, "wrapped message", WithCode(Internal))
		assert.Equal(t, "wrapped message", err.Message)
		assert.Equal(t, Internal, err.Code)
		assert.ErrorIs(t, err, originalErr)
	})
}

func TestError_ErrorMethod(t *testing.T) {
	assert.Equal(t, "message", New("message").Error())
	assert.Equal(t, "message: original", Wrap(errors.New("original"), "message").Error())
}

func TestHelperConstructors(t *testing.T) {
	t.Run("NewValidationError", func(t *testing.T) {
		err := NewValidationError(
			errors.New("db error"),
			"invalid email",
			map[string]any{"field": "email"},
		)
		assert.Equal(t, Invalid, err.Code)
		assert.Equal(t, "invalid email", err.Message)
		assert.Equal(t, "email", err.Context["field"])
		assert.ErrorContains(t, err, "db error")
	})

	t.Run("NewInternalError", func(t *testing.T) {
		err := NewInternalError(
			errors.New("timeout"),
			map[string]any{"service": "payment"},
		)
		assert.Equal(t, Internal, err.Code)
		assert.Equal(t, "An unexpected internal error occurred.", err.Message)
		assert.Equal(t, "payment", err.Context["service"])
		assert.ErrorContains(t, err, "timeout")
	})
}
