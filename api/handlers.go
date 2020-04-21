package api

import (
	"encoding/json"
	"github.com/ONSdigital/dp-sessions-api/session"
	"github.com/ONSdigital/log.go/log"
	"net/http"
)

func CreateNewSessionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var sess session.Session
		err := json.NewDecoder(r.Body).Decode(&sess)
		if err != nil {
			log.Event(ctx, "unable to decode session", log.Error(err), log.ERROR)
			http.Error(w, "Failed to decode session json", http.StatusBadRequest)
			return
		}

		jsonString, err := sess.MarshalJSON()
		if err != nil {
			log.Event(ctx, "marshalling response failed", log.Error(err), log.ERROR)
			http.Error(w, "Failed to marshall json response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonString)
	}
}
