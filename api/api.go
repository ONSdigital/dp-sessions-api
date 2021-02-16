package api

import (
	"context"

	"github.com/ONSdigital/dp-authorisation/auth"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

var (
	create = auth.Permissions{Create: true}
	delete = auth.Permissions{Delete: true}
)

//API provides a struct to wrap the api around
type API struct {
	Router *mux.Router
}

func Setup(ctx context.Context, r *mux.Router, permissions AuthHandler, cache Cache) *API {
	api := &API{
		Router: r,
	}

	r.HandleFunc("/sessions", permissions.Require(create, CreateSessionHandlerFunc(cache))).Methods("POST")
	r.HandleFunc("/sessions/{ID}", GetByIDSessionHandlerFunc(cache, mux.Vars)).Methods("GET")
	r.HandleFunc("/sessions/{Email}", GetByEmailSessionHandlerFunc(cache, mux.Vars)).Methods("GET")
	r.HandleFunc("/sessions", permissions.Require(delete, DeleteAllSessionsHandlerFunc(cache))).Methods("DELETE")
	return api
}

func (*API) Close(ctx context.Context) error {
	// Close any dependencies
	log.Event(ctx, "graceful shutdown of api complete", log.INFO)
	return nil
}
