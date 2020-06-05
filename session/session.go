package session

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

const (
	dateTimeFMT = "2006-01-02T15:04:05.000Z"
)

// Session defines the structure required for a session
type Session struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Start        time.Time `json:"start"`
	LastAccessed time.Time `json:"lastAccessed"`
}

// NewSessionDetails is the structure of the request needed to create a session
type NewSessionDetails struct {
	Email string `json:"email"`
}

// IDGenerator interface for creating new IDs
type IDGenerator interface {
	NewID() string
}

type jsonModel struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Start        string `json:"start"`
	LastAccessed string `json:"last_accessed"`
}

// NewSession creates a new session
func NewSession() *Session {
	return &Session{
		ID:           "",
		Email:        "",
		Start:        time.Time{},
		LastAccessed: time.Time{},
	}
}

// New creates a new session using email parameter
func (sess *Session) New(email string) (*Session, error) {
	id, err := newID()
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:           id,
		Email:        email,
		Start:        time.Now(),
		LastAccessed: time.Now(),
	}, nil
}

// MarshalJSON used to marshal Session object for outgoing requests
func (sess *Session) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonModel{
		ID:           sess.ID,
		Email:        sess.Email,
		Start:        sess.Start.Format(dateTimeFMT),
		LastAccessed: sess.LastAccessed.Format(dateTimeFMT),
	})
}

func newID() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
