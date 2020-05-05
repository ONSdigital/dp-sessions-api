package api

import (
	"github.com/ONSdigital/dp-sessions-api/session"
	"time"
)

const (
	NOPSessionID = "123"
	NOPStart     = "Mon Jan _2 15:04:05 2006"
)

type NOPSessions struct{}

type NOPCache struct{}

func (n *NOPSessions) New(email string) (*session.Session, error) {
	startP := parseTime()
	var sess = &session.Session{
		ID:    NOPSessionID,
		Email: email,
		Start: startP,
	}
	return sess, nil
}

func (n *NOPCache) Set(s *session.Session) {}

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
