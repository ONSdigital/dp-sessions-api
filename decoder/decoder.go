package decoder

import (
	"encoding/json"
	"net/http"
)

type validator interface {
	ValidateNewSessionDetails() error
}

// Decode is a custom decoder validator
func Decode(r *http.Request, v validator) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return v.ValidateNewSessionDetails()
}
