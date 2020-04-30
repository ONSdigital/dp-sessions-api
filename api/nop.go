package api

import (
	"github.com/ONSdigital/dp-sessions-api/session"
	"time"
)

const NOPSessionID = "123"

type NOPSessions struct{}

func (n *NOPSessions) New(email string) (*session.Session, error) {
	var sess = &session.Session{
		ID:    NOPSessionID,
		Email: email,
		Start: time.Now(),
	}
	return sess, nil
}
