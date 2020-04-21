package session

import (
	"encoding/json"
	"time"
)

type Session struct {
	ID    string    `json:"id"`
	Email string    `json:"email"`
	Start time.Time `json:"start"`
}

const (
	dateTimeFMT = "2006-01-02T15:04:05.000Z"
)

type IDGenerator interface {
	NewID() string
}

type jsonModel struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Start string `json:"start"`
}

func (sess *Session) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonModel{
		ID:    sess.ID,
		Email: sess.Email,
		Start: sess.Start.Format(dateTimeFMT),
	})
}
