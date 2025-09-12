package fault_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/marcelofabianov/fault"
)

func TestError(t *testing.T) {
	t.Run("Should implement standard error interface", func(t *testing.T) {
		err := &fault.Error{
			Message: "test message",
		}
		assert.Implements(t, (*error)(nil), err)
	})

	t.Run("Error method should return message only when no wrapped error", func(t *testing.T) {
		err := &fault.Error{
			Message: "test message",
		}
		assert.Equal(t, "test message", err.Error())
	})

	t.Run("Error method should return message and wrapped error", func(t *testing.T) {
		wrappedErr := errors.New("original error")
		err := &fault.Error{
			Message: "test message",
			Err:     wrappedErr,
		}
		assert.Equal(t, "test message: original error", err.Error())
	})

	t.Run("Unwrap method should return the wrapped error", func(t *testing.T) {
		wrappedErr := errors.New("original error")
		err := &fault.Error{
			Err: wrappedErr,
		}
		assert.Equal(t, wrappedErr, err.Unwrap())
	})
}
