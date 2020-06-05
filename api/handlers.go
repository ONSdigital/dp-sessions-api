package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	. "github.com/ONSdigital/dp-sessions-api/errors"
	"github.com/ONSdigital/dp-sessions-api/session"
	"github.com/ONSdigital/log.go/log"
)

// GetVarsFunc is a helper function that returns a map of request variables and parameters
type GetVarsFunc func(r *http.Request) map[string]string

// CreateSessionHandlerFunc returns a function that generates a session. Method = "POST"
func CreateSessionHandlerFunc(sess Session, cache Cache) http.HandlerFunc {
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

		sess, err := sess.New(c.Email)
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

		log.Event(ctx, "session added to cache", log.INFO)

		sessJSON, err := s.MarshalJSON()
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

		sess, err := cache.GetByID(ID)
		if err != nil {
			writeErrorResponse(ctx, w, err.Error(), err, getErrorStatus(err))
			return
		}

		if sess == nil {
			writeErrorResponse(ctx, w, "session not found", err, http.StatusNotFound)
			return
		}

		sessJSON, err := sess.MarshalJSON()
		if err != nil {
			writeErrorResponse(ctx, w, "failed to marshal session", err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(sessJSON)
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

func getErrorStatus(err error) int {
	var status int
	switch {
	case errors.Is(err, SessionNotFound):
		status = http.StatusNotFound
	case errors.Is(err, SessionExpired):
		status = http.StatusNotFound
	default:
		status = http.StatusInternalServerError
	}
	return status
}
