package session

import (
"encoding/json"

"github.com/google/uuid"
"github.com/pkg/errors"

"time"
)

const (
	DateTimeFMT = "2006-01-02T15:04:05.000Z"
)

var (
	EmailEmptyErr        = errors.New("error creating session email required but was empty")
	StartEmptyErr        = errors.New("error unmarshalling session start field required but was missing/empty")
	LastAccessedEmptyErr = errors.New("error unmarshalling session last accessed field required but was missing/empty")
)

// Session defines the structure required for a session
type Session struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Start        time.Time `json:"start"`
	LastAccessed time.Time `json:"last_accessed"`
}

// NewSessionDetails is the create HTTP request body required to creating new session
type NewSessionDetails struct {
	Email string `json:"email"`
}

type jsonModel struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Start        string `json:"start"`
	LastAccessed string `json:"last_accessed"`
}

//New construct a new fully populated session object for the provided email. Returns session.EmailEmptyErr if the email
//is empty/blank, returns an error if a new session ID could not be generated.
func New(email string) (*Session, error) {
	if len(email) == 0 {
		return nil, EmailEmptyErr
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.WithMessage(err, "error generating new session ID")
	}

	var createdAt time.Time
	createdAt, err = FormatTime(time.Now().UTC())
	if err != nil {
		return nil, errors.WithMessage(err, "error formatting session timestamp value")
	}

	return &Session{
		ID:           id.String(),
		Email:        email,
		Start:        createdAt,
		LastAccessed: createdAt,
	}, nil
}

// MarshalJSON is a custom JSON marshaller for Session objects. Handles marshalling time.Time fields into the expected date time format
func (s *Session) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonModel{
		ID:           s.ID,
		Email:        s.Email,
		Start:        s.Start.Format(DateTimeFMT),
		LastAccessed: s.LastAccessed.Format(DateTimeFMT),
	})
}

func (s *Session) UnmarshalJSON(data []byte) error {
	var raw jsonModel
	var err error

	if err = json.Unmarshal(data, &raw); err != nil {
		return errors.WithMessage(err, "failed to unmarshal session JSON")
	}

	if len(raw.Start) == 0 {
		return StartEmptyErr
	}

	if len(raw.LastAccessed) == 0 {
		return LastAccessedEmptyErr
	}

	var startT time.Time
	startT, err = time.Parse(DateTimeFMT, raw.Start)
	if err != nil {
		return errors.WithMessage(err, "error parsing session.Start as time.Time value")
	}

	var lastAccessedT time.Time
	lastAccessedT, err = time.Parse(DateTimeFMT, raw.LastAccessed)
	if err != nil {
		return errors.WithMessage(err, "error parsing session.LastAccessed as time.Time value")
	}

	s.ID = raw.ID
	s.Email = raw.Email
	s.Start = startT
	s.LastAccessed = lastAccessedT
	return nil
}

func FormatTime(t time.Time) (time.Time, error) {
	// Format time t with the desired layout then parse it back to a time.Time object.
	return time.Parse(DateTimeFMT, t.Format(DateTimeFMT))
}
