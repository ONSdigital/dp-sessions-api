package session

import (
	"encoding/json"
	errs "github.com/ONSdigital/dp-sessions-api/apierrors"
	"github.com/google/uuid"
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

// CreateSessionEntity is the structure of the request needed to create a session
type CreateSessionEntity struct {
	Email string `json:"email"`
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

// CreateNewSession creates a new session using email parameter
func CreateNewSession(email string) Session {
	return Session{
		ID: newID(),
		Email: email,
		Start: time.Now(),
	}
}

// OK checks CreateSessionEntity is valid
func (c *CreateSessionEntity) OK() error {
	if len(c.Email) == 0 {
		return errs.ErrMissingField
	}
	return nil
}

func newID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		return ""
	}
	return id.String()
}
