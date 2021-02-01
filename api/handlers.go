package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-sessions-api/cache"
	"github.com/ONSdigital/dp-sessions-api/session"
	"github.com/ONSdigital/log.go/log"
)

const (
	internalServerErr  = "internal server error"
	sessionNotFoundErr = "session not found"
)

// GetVarsFunc is a helper function that returns a map of request variables and parameters
type GetVarsFunc func(r *http.Request) map[string]string

// CreateSessionHandlerFunc returns a function that generates a session. Method = "POST"
func CreateSessionHandlerFunc(sessionCache Cache) http.HandlerFunc {
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

		s, err := session.NewSession().Update(c.Email)
		if err != nil {
			writeErrorResponse(ctx, w, "failed to create session", err, http.StatusInternalServerError)
			return
		}

		if err := sessionCache.SetSession(s); err != nil {
			writeErrorResponse(ctx, w, "unable to add session to cache", err, http.StatusInternalServerError)
			return
		}

		log.Event(ctx, "session added to cache", log.INFO)

		sessionJSON, err := s.MarshalJSON()
		if err != nil {
			writeErrorResponse(ctx, w, "failed to marshal session", err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(sessionJSON)
	}
}

// GetByIDSessionHandlerFunc returns a function that retrieves a session by ID from the cache
func GetByIDSessionHandlerFunc(sessionCache Cache, getVarsFunc GetVarsFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ID := getVarsFunc(r)["ID"]

		s, err := sessionCache.GetByID(ID)
		if err != nil {
			if err == cache.ErrSessionNotFound {
				writeErrorResponse(ctx, w, sessionNotFoundErr, err, http.StatusNotFound)
				return
			}

			writeErrorResponse(ctx, w, internalServerErr, err, http.StatusInternalServerError)
			return
		}

		if s == nil {
			writeErrorResponse(ctx, w, internalServerErr, err, http.StatusInternalServerError)
			return
		}

		sessionJSON, err := s.MarshalJSON()
		if err != nil {
			writeErrorResponse(ctx, w, internalServerErr, err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(sessionJSON)
	}
}

// GetByEmailSessionHandlerFunc returns a function that retrieves a session by ID from the cache
func GetByEmailSessionHandlerFunc(sessionCache Cache, getVarsFunc GetVarsFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		email := getVarsFunc(r)["Email"]

		s, err := sessionCache.GetByEmail(email)
		if err != nil {
			if err == cache.ErrSessionNotFound {
				writeErrorResponse(ctx, w, sessionNotFoundErr, err, http.StatusNotFound)
				return
			}

			writeErrorResponse(ctx, w, internalServerErr, err, http.StatusInternalServerError)
			return
		}

		if s == nil {
			writeErrorResponse(ctx, w, internalServerErr, err, http.StatusInternalServerError)
			return
		}

		sessionJSON, err := s.MarshalJSON()
		if err != nil {
			writeErrorResponse(ctx, w, internalServerErr, err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(sessionJSON)
	}
}

func DeleteAllSessionsHandlerFunc(cache Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if err := cache.DeleteAll(); err != nil {
			writeErrorResponse(ctx, w, "no sessions to delete", err, http.StatusNotFound)
			return
		}

		log.Event(ctx, "all sessions deleted", log.INFO)

		w.WriteHeader(http.StatusOK)
	}
}

func writeErrorResponse(ctx context.Context, w http.ResponseWriter, msg string, err error, status int) {
	if err != nil {
		log.Event(ctx, err.Error(), log.Error(err), log.ERROR)
	} else {
		log.Event(ctx, msg, log.ERROR)
	}

	http.Error(w, msg, status)
}
