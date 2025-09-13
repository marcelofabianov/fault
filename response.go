package fault

import (
	"errors"
)

type ErrorResponse struct {
	StatusCode int             `json:"-"`
	Message    string          `json:"message"`
	Code       string          `json:"code,omitempty"`
	Context    map[string]any  `json:"context,omitempty"`
	Details    []ErrorResponse `json:"details,omitempty"`
}

func ToResponse(err error) ErrorResponse {
	var fErr *Error

	if !errors.As(err, &fErr) {
		fErr = Wrap(err, "An unexpected internal error occurred.", WithCode(Internal))
	}

	return toResponse(fErr)
}

func toResponse(err *Error) ErrorResponse {
	resp := ErrorResponse{
		StatusCode: GetHTTPStatusCode(err.Code),
		Message:    err.Message,
		Code:       string(err.Code),
		Context:    err.Context,
	}
	for _, detail := range err.Details {
		resp.Details = append(resp.Details, toResponse(detail))
	}
	return resp
}
