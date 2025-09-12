package fault

import (
	"errors"
	"net/http"
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

func TestIsCode(t *testing.T) {
	t.Run("Check IsCode on a wrapped error chain", func(t *testing.T) {
		originalErr := New("validation failed", WithCode(Invalid))
		wrappedErr := Wrap(originalErr, "could not process request")

		// Debugging steps
		t.Logf("wrappedErr type: %T, value: %+v", wrappedErr, wrappedErr)
		t.Logf("wrappedErr.Unwrap() type: %T, value: %+v", wrappedErr.Unwrap(), wrappedErr.Unwrap())
		t.Logf("IsCode result: %t", IsCode(wrappedErr, Invalid))

		assert.True(t, IsCode(wrappedErr, Invalid)) // Linha 83
		assert.False(t, IsCode(wrappedErr, NotFound))
	})
}

func TestGetHTTPStatusCode(t *testing.T) {
	t.Run("Check all known fault codes", func(t *testing.T) {
		testCases := map[Code]int{
			Invalid:         http.StatusBadRequest,
			Conflict:        http.StatusConflict,
			NotFound:        http.StatusNotFound,
			Unauthorized:    http.StatusUnauthorized,
			Forbidden:       http.StatusForbidden,
			DomainViolation: http.StatusUnprocessableEntity,
			InfraError:      http.StatusBadGateway,
			Internal:        http.StatusInternalServerError,
		}

		for code, expectedStatus := range testCases {
			t.Run(string(code), func(t *testing.T) {
				status := GetHTTPStatusCode(code)
				assert.Equal(t, expectedStatus, status, "Expected status for code %s", code)
			})
		}
	})

	t.Run("Check a code with no mapping", func(t *testing.T) {
		const NewCode Code = "new_test_code"
		status := GetHTTPStatusCode(NewCode)
		assert.Equal(t, http.StatusInternalServerError, status)
	})
}

func TestIsSpecificCodes(t *testing.T) {
	testCases := []struct {
		name     string
		errCode  Code
		isFunc   func(error) bool
		expected bool
	}{
		{"IsDomainViolation", DomainViolation, IsDomainViolation, true},
		{"IsInfraError", InfraError, IsInfraError, true},
		{"IsNotFound", NotFound, IsNotFound, true},
		{"IsUnauthorized", Unauthorized, IsUnauthorized, true},
		{"IsForbidden", Forbidden, IsForbidden, true},
		{"IsConflict", Conflict, IsConflict, true},
		{"IsInvalid", Invalid, IsInvalid, true},
		{"IsInternal", Internal, IsInternal, true},
		{"IsDomainViolation (negative)", NotFound, IsDomainViolation, false},
		{"IsInternal (negative)", Invalid, IsInternal, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := New("test error", WithCode(tc.errCode))
			assert.Equal(t, tc.expected, tc.isFunc(err))
		})
	}
}
