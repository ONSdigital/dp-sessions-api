package api

import (
	"context"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

//API provides a struct to wrap the api around
type API struct {
	Router *mux.Router
}

func Setup(ctx context.Context, r *mux.Router) *API {
	api := &API{
		Router: r,
	}

	nopSess := &NOPSessions{}
	nopCache := &NOPCache{}

	r.HandleFunc("/session", CreateSessionHandlerFunc(nopSess, nopCache)).Methods("POST")
	r.HandleFunc("/session/{ID}", GetByIDSessionHandlerFunc(nopCache, mux.Vars)).Methods("GET")
	return api
}

func (*API) Close(ctx context.Context) error {
	// Close any dependencies
	log.Event(ctx, "graceful shutdown of api complete", log.INFO)
	return nil
}
