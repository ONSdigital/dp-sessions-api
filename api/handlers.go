package api

//go:generate moq -out mocksessions_test.go . Sessions
//go:generate moq -out mockcache_test.go . Cache

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-sessions-api/session"
	"github.com/ONSdigital/log.go/log"
)


// Sessions interface for getting a new session
type Sessions interface {
	New(email string) (*session.Session, error)
}

// Cache interface for storing and retrieving sessions
type Cache interface {
	Set(s *session.Session) error
	GetByID(ID string) (*session.Session, error)
}

// GetVarsFunc is a helper function that returns a map of request variables and parameters
type GetVarsFunc func(r *http.Request) map[string]string

// CreateSessionHandlerFunc returns a function that generates a session. Method = "POST"
func CreateSessionHandlerFunc(sessions Sessions, cache Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var c session.NewSessionDetails
		err := json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
			writeErrorResponse(ctx, w, "failed to unmarshal request body", err, http.StatusBadRequest)
			return
		}

		if len(c.Email) == 0 {
			writeErrorResponse(ctx, w, "missing email field in json", err, http.StatusBadRequest)
			return
		}

		sess, err := sessions.New(c.Email)
		if err != nil {
			writeErrorResponse(ctx, w, "failed to create session", err, http.StatusInternalServerError)
			return
		}

		s := &session.Session{
			ID:    sess.ID,
			Email: sess.Email,
			Start: sess.Start,
		}

		if err := cache.Set(s); err != nil {
			writeErrorResponse(ctx, w, "unable to add session to cache", err, http.StatusInternalServerError)
			return
		}

		sessJSON, err := sess.MarshalJSON()
		if err != nil {
			writeErrorResponse(ctx, w, "failed to marshal session", err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(sessJSON)
	}
}

// GetByIDSessionHandlerFunc returns a function that retrieves a session by ID from the cache
func GetByIDSessionHandlerFunc(cache Cache, getVarsFunc GetVarsFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ID := getVarsFunc(r)["ID"]

		result, err := cache.GetByID(ID)
		if err != nil {
			writeErrorResponse(ctx, w, "unable to get session by id", err, http.StatusInternalServerError)
			return
		}

		if result == nil {
			writeErrorResponse(ctx, w,"session not found", err, http.StatusNotFound)
			return
		}

		resultJSON, err := result.MarshalJSON()
		if err != nil {
			writeErrorResponse(ctx, w, "failed to marshal session", err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resultJSON)
	}
}

func writeErrorResponse(ctx context.Context, w http.ResponseWriter, msg string, err error, status int) {
	if err != nil {
		log.Event(ctx, msg, log.Error(err), log.ERROR)
	} else {
		log.Event(ctx, msg, log.ERROR)
	}
	http.Error(w, msg, status)
}
