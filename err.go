package fault

import "fmt"

type Code string

const (
	Conflict        Code = "conflict"
	Invalid         Code = "invalid_input"
	NotFound        Code = "not_found"
	Internal        Code = "internal_error"
	Unauthorized    Code = "unauthorized"
	Forbidden       Code = "forbidden"
	DomainViolation Code = "domain_violation"
	InfraError      Code = "infra_error"
)

type Error struct {
	Err     error
	Message string
	Code    Code
	Context map[string]any
	Details []*Error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}
