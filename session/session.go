package session

import (
	"encoding/json"
	"time"
)

const (
	dateTimeFMT = "2006-01-02T15:04:05.000Z"
)

// Session defines the structure required for a session
type Session struct {
	ID    string    `json:"id"`
	Email string    `json:"email"`
	Start time.Time `json:"start"`
}

// IDGenerator interface for creating new IDs
type IDGenerator interface {
	NewID() string
}

type jsonModel struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Start string `json:"start"`
}

// MarshalJSON used to marshal Session object for outgoing requests
func (sess *Session) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonModel{
		ID:    sess.ID,
		Email: sess.Email,
		Start: sess.Start.Format(dateTimeFMT),
	})
}
