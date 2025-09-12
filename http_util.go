package fault

import "net/http"

var httpStatusCodes = map[Code]int{
	Invalid:         http.StatusBadRequest,
	Conflict:        http.StatusConflict,
	NotFound:        http.StatusNotFound,
	Unauthorized:    http.StatusUnauthorized,
	Forbidden:       http.StatusForbidden,
	DomainViolation: http.StatusUnprocessableEntity,
	InfraError:      http.StatusBadGateway,
	Internal:        http.StatusInternalServerError,
}

func GetHTTPStatusCode(code Code) int {
	if statusCode, ok := httpStatusCodes[code]; ok {
		return statusCode
	}
	return http.StatusInternalServerError
}
