package fault

type ErrorResponse struct {
	StatusCode int             `json:"-"`
	Message    string          `json:"message"`
	Code       string          `json:"code,omitempty"`
	Context    map[string]any  `json:"context,omitempty"`
	Details    []ErrorResponse `json:"details,omitempty"`
}

func ToResponse(err *Error) ErrorResponse {
	return toResponse(err)
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
