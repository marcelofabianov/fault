package httputil

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/marcelofabianov/fault"
)

func TestHTTPStatus(t *testing.T) {
	testCases := []struct {
		name           string
		code           fault.Code
		expectedStatus int
	}{
		{"Conflict", fault.Conflict, http.StatusConflict},
		{"Invalid Input", fault.Invalid, http.StatusBadRequest},
		{"Not Found", fault.NotFound, http.StatusNotFound},
		{"Internal Error", fault.Internal, http.StatusInternalServerError},
		{"Unauthorized", fault.Unauthorized, http.StatusUnauthorized},
		{"Forbidden", fault.Forbidden, http.StatusForbidden},
		{"Unknown Code", fault.Code("SOME_NEW_CODE"), http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status := HTTPStatus(tc.code)
			assert.Equal(t, tc.expectedStatus, status)
		})
	}
}

func TestToResponse(t *testing.T) {
	t.Run("converts a simple error correctly", func(t *testing.T) {
		err := fault.New(
			"Could not process request",
			fault.WithCode(fault.Internal),
			fault.WithContext("request_id", "abc-123"),
		)
		response := ToResponse(err)

		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
		assert.Equal(t, "Could not process request", response.Message)
		assert.Equal(t, string(fault.Internal), response.Code)
		assert.Equal(t, map[string]any{"request_id": "abc-123"}, response.Context)
	})

	t.Run("converts an error with details correctly", func(t *testing.T) {
		detail1 := fault.New("must be a valid email", fault.WithCode(fault.Invalid), fault.WithContext("field", "email"))
		parentErr := fault.New(
			"One or more fields are invalid",
			fault.WithCode(fault.Invalid),
			fault.WithContext("form", "registration"),
			fault.WithDetails(detail1),
		)
		response := ToResponse(parentErr)

		require.Len(t, response.Details, 1)
		assert.Equal(t, "must be a valid email", response.Details[0].Message)
	})
}
