package errors

import "errors"

var (
	SessionNotFound = errors.New("unable to get session")
	SessionExpired = errors.New("session has expired")
)
