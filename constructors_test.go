package fault_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/marcelofabianov/fault"
)

func TestConstructors(t *testing.T) {
	originalErr := errors.New("original error")

	t.Run("New should create a new error with options", func(t *testing.T) {
		detail := fault.New("detail message")
		err := fault.New(
			"user message",
			fault.WithWrappedErr(originalErr),
			fault.WithCode(fault.Invalid),
			fault.WithContext("key1", "value1"),
			fault.WithDetails(detail),
		)

		require.NotNil(t, err)
		assert.Equal(t, "user message", err.Message)
		assert.Equal(t, originalErr, err.Unwrap())
		assert.Equal(t, fault.Invalid, err.Code)
		assert.Equal(t, map[string]any{"key1": "value1"}, err.Context)
		require.Len(t, err.Details, 1)
		assert.Equal(t, detail, err.Details[0])
	})

	t.Run("Wrap should create a new error wrapping an existing one", func(t *testing.T) {
		err := fault.Wrap(originalErr, "wrapped message", fault.WithCode(fault.Internal))
		assert.Equal(t, "wrapped message", err.Message)
		assert.Equal(t, fault.Internal, err.Code)
		assert.ErrorIs(t, err, originalErr)
	})

	t.Run("NewValidationError should create a new error with Invalid code", func(t *testing.T) {
		err := fault.NewValidationError(
			errors.New("db error"),
			"invalid email",
			map[string]any{"field": "email"},
		)
		assert.Equal(t, fault.Invalid, err.Code)
		assert.Equal(t, "invalid email", err.Message)
		assert.Equal(t, "email", err.Context["field"])
		assert.ErrorContains(t, err, "db error")
	})

	t.Run("NewInternalError should create a new error with Internal code", func(t *testing.T) {
		err := fault.NewInternalError(
			errors.New("timeout"),
			map[string]any{"service": "payment"},
		)
		assert.Equal(t, fault.Internal, err.Code)
		assert.Equal(t, "An unexpected internal error occurred.", err.Message)
		assert.Equal(t, "payment", err.Context["service"])
		assert.ErrorContains(t, err, "timeout")
	})
}
