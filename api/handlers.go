package api

//go:generate moq -out mocksessions_test.go . Sessions
//go:generate moq -out mockcache_test.go . Cache

import (
	"context"
	"encoding/json"
	"fmt"
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
	Set(s *session.Session)
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
			log.Event(ctx, "failed to unmarshal request body", log.Error(err), log.ERROR)
			http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
			return
		}

		if len(c.Email) == 0 {
			log.Event(ctx, "missing email field in json", log.ERROR)
			http.Error(w, "Missing email field in json", http.StatusBadRequest)
			return
		}

		sess, err := sessions.New(c.Email)
		if err != nil {
			log.Event(ctx, "failed to create session", log.Error(err), log.ERROR)
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}

		log.Event(ctx, "session created", log.Data{"session": sess}, log.INFO)
		s := &session.Session{
			ID:    sess.ID,
			Email: sess.Email,
			Start: sess.Start,
		}

		cache.Set(s)
		log.Event(ctx, "session added to cache", log.INFO)

		sessJSON, err := sess.MarshalJSON()
		if err != nil {
			failedToMarshalError(ctx, err, w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(sessJSON)
	}
}

func GetByIDSessionHandlerFunc(cache Cache, getVarsFunc GetVarsFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ID := getVarsFunc(r)["ID"]

		log.Event(ctx, fmt.Sprintf("get session by ID: %s", ID), log.INFO)

		result, err := cache.GetByID(ID)
		if err != nil {
			log.Event(ctx,"unable to get session by id", log.Error(err), log.ERROR)
			http.Error(w, "Unable to get session by id", http.StatusBadRequest)
			return
		}

		resultJSON, err := result.MarshalJSON()
		if err != nil {
			failedToMarshalError(ctx, err, w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resultJSON)
	}
}

func failedToMarshalError(ctx context.Context, err error, w http.ResponseWriter) {
	log.Event(ctx, "failed to marshal session", log.Error(err), log.ERROR)
	http.Error(w, "Failed to marshal session", http.StatusInternalServerError)
}
