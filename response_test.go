package fault_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/marcelofabianov/fault"
)

func TestResponse(t *testing.T) {
	t.Run("ToResponse should convert a simple error correctly", func(t *testing.T) {
		err := fault.New(
			"User not found",
			fault.WithCode(fault.NotFound),
			fault.WithContext("user_id", "123"),
		)
		response := fault.ToResponse(err)

		assert.Equal(t, "User not found", response.Message)
		assert.Equal(t, "not_found", response.Code)
		assert.Equal(t, 404, response.StatusCode)
		assert.Equal(t, map[string]any{"user_id": "123"}, response.Context)
		assert.Empty(t, response.Details)
	})

	t.Run("ToResponse should convert a nested error correctly", func(t *testing.T) {
		detailErr := fault.New(
			"Invalid field",
			fault.WithCode(fault.Invalid),
			fault.WithContext("field", "email"),
		)
		parentErr := fault.New(
			"Validation failed",
			fault.WithCode(fault.DomainViolation),
			fault.WithDetails(detailErr),
		)

		response := fault.ToResponse(parentErr)

		assert.Equal(t, "Validation failed", response.Message)
		assert.Equal(t, "domain_violation", response.Code)
		assert.Equal(t, 422, response.StatusCode)
		assert.Empty(t, response.Context)
		assert.Len(t, response.Details, 1)

		detailResponse := response.Details[0]
		assert.Equal(t, "Invalid field", detailResponse.Message)
		assert.Equal(t, "invalid_input", detailResponse.Code)
		assert.Equal(t, 400, detailResponse.StatusCode)
		assert.Equal(t, map[string]any{"field": "email"}, detailResponse.Context)
	})

	t.Run("ToResponse should handle an error without code or context", func(t *testing.T) {
		err := fault.New("Internal server error")
		response := fault.ToResponse(err)

		assert.Equal(t, "Internal server error", response.Message)
		assert.Empty(t, response.Code)
		assert.Equal(t, 500, response.StatusCode)
		assert.Empty(t, response.Context)
		assert.Empty(t, response.Details)
	})
}
