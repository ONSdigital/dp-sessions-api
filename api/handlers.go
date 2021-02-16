package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/ONSdigital/dp-sessions-api/cache"
	"github.com/ONSdigital/dp-sessions-api/session"
	"github.com/ONSdigital/log.go/log"
	"github.com/pkg/errors"
)

const (
	internalServerErr    = "internal server error"
	sessionNotFoundErr   = "session not found"
	marshallSessionErr   = "failed to marshal session to JSON"
	unmarshallSessionErr = "failed to unmarshal session JSON"
	sessionEmailEmptyErr = "session.Email required but was empty"
	createSessionErr     = "error creating new session"
	addSessionToCacheErr = "error adding new session to cache"
)

var (
	sessionNilErr = errors.New("expected session object but was nil")
)

// GetVarsFunc is a helper function that returns a map of request variables and parameters
type GetVarsFunc func(r *http.Request) map[string]string

// CreateSessionHandlerFunc returns HTTP HandlerFunc for handling POST requests to create sessions.
func CreateSessionHandlerFunc(sessionCache Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		email, getEmailErr := getEmailForNewSession(r.Body)
		if getEmailErr != nil {
			writeErrorResponse(ctx, w, sessionEmailEmptyErr, getEmailErr, http.StatusBadRequest)
			return
		}

		s, newSessErr := session.New(email)
		if newSessErr != nil {
			writeErrorResponse(ctx, w, createSessionErr, newSessErr, http.StatusInternalServerError)
			return
		}

		if cacheSessErr := sessionCache.SetSession(s); cacheSessErr != nil {
			writeErrorResponse(ctx, w, addSessionToCacheErr, cacheSessErr, http.StatusInternalServerError)
			return
		}

		log.Event(ctx, "session was successfully added to cache", log.INFO, log.Data{"email": s.Email})

		sessionJSON, marshalErr := s.MarshalJSON()
		if marshalErr != nil {
			writeErrorResponse(ctx, w, marshallSessionErr, marshalErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(sessionJSON)
	}
}

func getEmailForNewSession(r io.Reader) (string, error) {
	var details session.NewSessionDetails
	if err := json.NewDecoder(r).Decode(&details); err != nil {
		return "", errors.WithMessage(err, unmarshallSessionErr)
	}

	if len(details.Email) == 0 {
		return "", errors.New(sessionEmailEmptyErr)
	}

	return details.Email, nil
}

// GetByIDSessionHandlerFunc returns a HTTP HandlerFunc that attempts to retrieve an existing session by ID from the cache
func GetByIDSessionHandlerFunc(sessionCache Cache, getVarsFunc GetVarsFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ID := getVarsFunc(r)["ID"]

		s, getSessErr := sessionCache.GetByID(ID)
		if getSessErr != nil {
			if getSessErr == cache.ErrSessionNotFound {
				writeErrorResponse(ctx, w, sessionNotFoundErr, getSessErr, http.StatusNotFound)
				return
			}

			writeErrorResponse(ctx, w, internalServerErr, getSessErr, http.StatusInternalServerError)
			return
		}

		if s == nil {
			writeErrorResponse(ctx, w, internalServerErr, sessionNilErr, http.StatusInternalServerError)
			return
		}

		sessionJSON, marshalErr := json.Marshal(s)
		if marshalErr != nil {
			writeErrorResponse(ctx, w, marshallSessionErr, marshalErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(sessionJSON)
	}
}

// GetByEmailSessionHandlerFunc returns a HTTP HandlerFunc that attempts to retrieve an existing session by Email from the cache
func GetByEmailSessionHandlerFunc(sessionCache Cache, getVarsFunc GetVarsFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		email := getVarsFunc(r)["Email"]

		s, getSessErr := sessionCache.GetByEmail(email)
		if getSessErr != nil {
			if getSessErr == cache.ErrSessionNotFound {
				writeErrorResponse(ctx, w, sessionNotFoundErr, getSessErr, http.StatusNotFound)
				return
			}

			writeErrorResponse(ctx, w, internalServerErr, getSessErr, http.StatusInternalServerError)
			return
		}

		if s == nil {
			writeErrorResponse(ctx, w, internalServerErr, sessionNilErr, http.StatusInternalServerError)
			return
		}

		sessionJSON, marshalErr := json.Marshal(s)
		if marshalErr != nil {
			writeErrorResponse(ctx, w, internalServerErr, marshalErr, http.StatusInternalServerError)
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
	log.Event(ctx, err.Error(), log.ERROR, log.Error(err))
	http.Error(w, msg, status)
}
