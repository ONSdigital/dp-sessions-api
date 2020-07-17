package errors

import "errors"

var (
	SessionNotFound = errors.New("session not found")
	SessionExpired = errors.New("session has expired")
)
