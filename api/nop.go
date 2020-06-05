package api

import (
	"github.com/ONSdigital/dp-sessions-api/session"
	"time"
)

const (
	NOPSessionID = "123"
	NOPStart     = "Mon Jan _2 15:04:05 2006"
)

// NOPSessions no-op struct
type NOPSession struct{}

// NOPCache no-op struct
type NOPCache struct{}

// New creates a new session using an email address
func (n *NOPSession) New(email string) (*session.Session, error) {
	startP := parseTime()
	var sess = &session.Session{
		ID:    NOPSessionID,
		Email: email,
		Start: startP,
	}
	return sess, nil
}

// Set stores a session into the cache
func (n *NOPCache) Set(s *session.Session) error {
	return nil
}

// GetByID retrieves a session from the cache by ID
func (n *NOPCache) GetByID(ID string) (*session.Session, error) {
	startP := parseTime()
	return &session.Session{
		ID:    ID,
		Email: "test@test.com",
		Start: startP,
	}, nil
}

func parseTime() time.Time {
	startP, err := time.Parse(time.ANSIC, NOPStart)
	if err != nil {
		return time.Time{}
	}
	return startP
}
