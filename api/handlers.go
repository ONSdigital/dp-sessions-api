package api

//go:generate moq -out mocks.go . Sessions

import (
	"encoding/json"
	"github.com/ONSdigital/dp-sessions-api/session"
	"github.com/ONSdigital/log.go/log"
	"net/http"
)


// Sessions interface for getting a new session
type Sessions interface {
	New(email string) (*session.Session, error)
}

// CreateSessionHandlerFunc returns a function that generates a session. Method = "POST"
func CreateSessionHandlerFunc(sessions Sessions) http.HandlerFunc {
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

		sessJson, err := sess.MarshalJSON()
		if err != nil {
			log.Event(ctx, "failed to marshal session", log.Error(err), log.ERROR)
			http.Error(w, "Failed to marshal session", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(sessJson)
	}
}
