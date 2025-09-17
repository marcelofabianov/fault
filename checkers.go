package fault

import "errors"

func AsFault(err error) (*Error, bool) {
	var fErr *Error
	ok := errors.As(err, &fErr)
	return fErr, ok
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
