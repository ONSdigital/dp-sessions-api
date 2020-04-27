package apierrors

import "errors"

var (
	ErrMissingField = errors.New("missing a required field")
)
