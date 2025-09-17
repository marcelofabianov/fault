package fault_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/marcelofabianov/fault"
)

func TestCheckers(t *testing.T) {
	t.Run("IsCode should check a single error", func(t *testing.T) {
		err := fault.New("resource not found", fault.WithCode(fault.NotFound))
		assert.True(t, fault.IsCode(err, fault.NotFound))
		assert.False(t, fault.IsCode(err, fault.Invalid))
	})

	t.Run("IsCode should check a wrapped error chain", func(t *testing.T) {
		originalErr := fault.New("validation failed", fault.WithCode(fault.Invalid))
		wrappedErr := fault.Wrap(originalErr, "could not process request")

		assert.True(t, fault.IsCode(wrappedErr, fault.Invalid))
		assert.False(t, fault.IsCode(wrappedErr, fault.NotFound))
	})

	t.Run("IsCode should return false for generic errors", func(t *testing.T) {
		genericErr := errors.New("a simple error")
		assert.False(t, fault.IsCode(genericErr, fault.NotFound))
	})

	t.Run("Specific Is* functions should work correctly", func(t *testing.T) {
		testCases := []struct {
			name     string
			code     fault.Code
			isFunc   func(error) bool
			expected bool
		}{
			{"IsDomainViolation", fault.DomainViolation, fault.IsDomainViolation, true},
			{"IsInfraError", fault.InfraError, fault.IsInfraError, true},
			{"IsNotFound", fault.NotFound, fault.IsNotFound, true},
			{"IsUnauthorized", fault.Unauthorized, fault.IsUnauthorized, true},
			{"IsForbidden", fault.Forbidden, fault.IsForbidden, true},
			{"IsConflict", fault.Conflict, fault.IsConflict, true},
			{"IsInvalid", fault.Invalid, fault.IsInvalid, true},
			{"IsInternal", fault.Internal, fault.IsInternal, true},
			{"IsDomainViolation (negative)", fault.NotFound, fault.IsDomainViolation, false},
			{"IsInternal (negative)", fault.Invalid, fault.IsInternal, false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := fault.New("test error", fault.WithCode(tc.code))
				assert.Equal(t, tc.expected, tc.isFunc(err))
			})
		}
	})

	t.Run("Specific Is* functions should work with wrapped errors", func(t *testing.T) {
		wrappedError := fault.Wrap(fault.New("original error", fault.WithCode(fault.Invalid)), "wrapped message")
		assert.True(t, fault.IsInvalid(wrappedError), "should be true for a wrapped 'Invalid' error")
		assert.False(t, fault.IsNotFound(wrappedError), "should be false for a wrapped 'NotFound' error")
	})

	t.Run("AsFault should work correctly", func(t *testing.T) {
		t.Run("should return true and the error for a direct fault.Error", func(t *testing.T) {
			originalErr := fault.New("direct error", fault.WithCode(fault.Internal))
			fErr, ok := fault.AsFault(originalErr)

			assert.True(t, ok)
			assert.NotNil(t, fErr)
			assert.Same(t, originalErr, fErr)
			assert.Equal(t, "direct error", fErr.Message)
		})

		t.Run("should return true and the error for a wrapped fault.Error", func(t *testing.T) {
			originalErr := fault.New("db connection failed", fault.WithCode(fault.InfraError))
			wrappedErr := fault.Wrap(originalErr, "could not fetch user")
			fErr, ok := fault.AsFault(wrappedErr)

			assert.True(t, ok)
			assert.NotNil(t, fErr)
			assert.Same(t, wrappedErr, fErr)
			assert.Equal(t, "could not fetch user", fErr.Message)
		})

		t.Run("should return false for a generic error", func(t *testing.T) {
			genericErr := errors.New("generic error")
			fErr, ok := fault.AsFault(genericErr)

			assert.False(t, ok)
			assert.Nil(t, fErr)
		})

		t.Run("should return false for a nil error", func(t *testing.T) {
			fErr, ok := fault.AsFault(nil)

			assert.False(t, ok)
			assert.Nil(t, fErr)
		})
	})
}
