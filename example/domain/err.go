package domain

import (
	"errors"

	"github.com/marcelofabianov/fault"
)

type AGGREGATE string

const (
	AGGREGATE_USER AGGREGATE = "user"
)

// Sentinelas para checagem de identidade com `errors.Is`.
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserInactive      = errors.New("user inactive")
)

// Construtores que criam erros `fault` ricos, encapsulando os sentinelas.

func NewUserNotFound(id string) error {
	return fault.Wrap(
		ErrUserNotFound,
		"user not found",
		fault.WithCode(fault.NotFound),
		fault.WithContext("aggregate", AGGREGATE_USER),
		fault.WithContext("user_id", id),
	)
}

func NewUserAlreadyExists(field, value string) error {
	return fault.Wrap(
		ErrUserAlreadyExists,
		"user already exists",
		fault.WithCode(fault.Conflict),
		fault.WithContext("aggregate", AGGREGATE_USER),
		fault.WithContext(field, value),
	)
}

func NewUserInactive(id string) error {
	return fault.Wrap(
		ErrUserInactive,
		"user inactive",
		fault.WithCode(fault.DomainViolation),
		fault.WithContext("aggregate", AGGREGATE_USER),
		fault.WithContext("user_id", id),
	)
}
