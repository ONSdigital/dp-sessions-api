package api

//go:generate moq -out mock/mockauth.go -pkg mock . AuthHandler
//go:generate moq -out mock/mocksessions.go -pkg mock . Sessions
//go:generate moq -out mock/mockcache.go -pkg mock . Cache

import (
	"github.com/ONSdigital/dp-authorisation/auth"
	"github.com/ONSdigital/dp-sessions-api/session"
	"net/http"
)

// AuthHandler interface for adding auth to endpoints
type AuthHandler interface {
	Require(required auth.Permissions, handler http.HandlerFunc) http.HandlerFunc
}

// Sessions interface for getting a new session
type Sessions interface {
	New(email string) (*session.Session, error)
}

// Cache interface for storing and retrieving sessions
type Cache interface {
	Set(s *session.Session) error
	GetByID(ID string) (*session.Session, error)
	DeleteAll() error
}
