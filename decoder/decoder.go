package decoder

import (
	"encoding/json"
	"net/http"
)

type ok interface {
	OK() error
}

// Decode is a custom decoder validator
func Decode(r *http.Request, v ok) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return v.OK()
}
