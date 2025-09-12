package fault

import (
	"errors"
	"fmt"
	"net/http"
)

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

type Option func(*Error)

func New(message string, opts ...Option) *Error {
	err := &Error{Message: message}
	for _, opt := range opts {
		opt(err)
	}
	return err
}

func Wrap(err error, message string, opts ...Option) *Error {
	opts = append(opts, WithWrappedErr(err))
	return New(message, opts...)
}

func WithWrappedErr(err error) Option {
	return func(e *Error) {
		e.Err = err
	}
}

func WithCode(code Code) Option {
	return func(e *Error) {
		e.Code = code
	}
}

func WithContext(key string, value any) Option {
	return func(e *Error) {
		if e.Context == nil {
			e.Context = make(map[string]any)
		}
		e.Context[key] = value
	}
}

func WithDetails(details ...*Error) Option {
	return func(e *Error) {
		e.Details = append(e.Details, details...)
	}
}

func NewValidationError(err error, message string, context map[string]any) *Error {
	opts := []Option{WithCode(Invalid)}
	if err != nil {
		opts = append(opts, WithWrappedErr(err))
	}
	for k, v := range context {
		opts = append(opts, WithContext(k, v))
	}
	return New(message, opts...)
}

func NewInternalError(err error, context map[string]any) *Error {
	opts := []Option{WithCode(Internal)}
	if err != nil {
		opts = append(opts, WithWrappedErr(err))
	}
	for k, v := range context {
		opts = append(opts, WithContext(k, v))
	}
	return New("An unexpected internal error occurred.", opts...)
}

func IsCode(err error, code Code) bool {
	for err != nil {
		if fErr, ok := err.(*Error); ok {
			if fErr.Code == code {
				return true
			}
		}
		err = errors.Unwrap(err)
	}
	return false
}

func IsDomainViolation(err error) bool {
	return IsCode(err, DomainViolation)
}

func IsInfraError(err error) bool {
	return IsCode(err, InfraError)
}

func IsNotFound(err error) bool {
	return IsCode(err, NotFound)
}

func IsUnauthorized(err error) bool {
	return IsCode(err, Unauthorized)
}

func IsForbidden(err error) bool {
	return IsCode(err, Forbidden)
}

func IsConflict(err error) bool {
	return IsCode(err, Conflict)
}

func IsInvalid(err error) bool {
	return IsCode(err, Invalid)
}

func IsInternal(err error) bool {
	return IsCode(err, Internal)
}

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
