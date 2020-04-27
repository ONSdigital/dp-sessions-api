package api

import (
	"fmt"
	"github.com/ONSdigital/dp-sessions-api/decoder"
	. "github.com/ONSdigital/dp-sessions-api/session"
	"github.com/ONSdigital/log.go/log"
	"net/http"
)

// CreateSessionHandlerFunc returns a function that generates a session. Method = "POST"
func CreateSessionHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var c CreateSessionEntity
		if err := decoder.Decode(r, &c); err != nil {
			log.Event(ctx, "missing a required field", log.Error(err), log.ERROR)
			http.Error(w, "Missing a required field", http.StatusBadRequest)
			return
		}

		sess := CreateNewSession(c.Email)

		log.Event(ctx, "session created", log.Data{"session": sess}, log.INFO)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Location", fmt.Sprintf("/session/%s", sess.ID))
		w.WriteHeader(http.StatusCreated)
	}
}
