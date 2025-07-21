package httputil

import (
	"net/http"

	"github.com/marcelofabianov/fault"
)

type ErrorResponse struct {
	StatusCode int             `json:"-"`
	Message    string          `json:"message"`
	Code       string          `json:"code,omitempty"`
	Context    map[string]any  `json:"context,omitempty"`
	Details    []ErrorResponse `json:"details,omitempty"`
}

func ToResponse(err *fault.Error) ErrorResponse {
	resp := ErrorResponse{
		StatusCode: HTTPStatus(err.Code),
		Message:    err.Message,
		Code:       string(err.Code),
		Context:    err.Context,
	}
	for _, detail := range err.Details {
		resp.Details = append(resp.Details, ToResponse(detail))
	}
	return resp
}

func HTTPStatus(code fault.Code) int {
	switch code {
	case fault.Conflict:
		return http.StatusConflict
	case fault.Invalid:
		return http.StatusBadRequest
	case fault.NotFound:
		return http.StatusNotFound
	case fault.Unauthorized:
		return http.StatusUnauthorized
	case fault.Forbidden:
		return http.StatusForbidden
	case fault.Internal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
