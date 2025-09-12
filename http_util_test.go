package fault_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/marcelofabianov/fault"
)

func TestHTTPUtil(t *testing.T) {
	t.Run("GetHTTPStatusCode should return correct HTTP status for a given code", func(t *testing.T) {
		testCases := map[fault.Code]int{
			fault.Invalid:         http.StatusBadRequest,
			fault.Conflict:        http.StatusConflict,
			fault.NotFound:        http.StatusNotFound,
			fault.Unauthorized:    http.StatusUnauthorized,
			fault.Forbidden:       http.StatusForbidden,
			fault.DomainViolation: http.StatusUnprocessableEntity,
			fault.InfraError:      http.StatusBadGateway,
			fault.Internal:        http.StatusInternalServerError,
		}

		for code, expectedStatus := range testCases {
			t.Run(string(code), func(t *testing.T) {
				status := fault.GetHTTPStatusCode(code)
				assert.Equal(t, expectedStatus, status, "Expected status for code %s", code)
			})
		}
	})

	t.Run("GetHTTPStatusCode should return InternalServerError for unknown codes", func(t *testing.T) {
		const newCode fault.Code = "new_test_code"
		status := fault.GetHTTPStatusCode(newCode)
		assert.Equal(t, http.StatusInternalServerError, status)
	})
}
