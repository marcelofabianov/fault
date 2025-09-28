package domain

import "github.com/marcelofabianov/wisp"

type User struct {
	ID     wisp.UUID
	Name   wisp.NonEmptyString
	Email  wisp.NonEmptyString
	Active bool
}
