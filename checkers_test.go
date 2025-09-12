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
}
