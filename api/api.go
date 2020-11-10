package api

import (
	"context"
	"github.com/ONSdigital/dp-authorisation/auth"
	dpredis "github.com/ONSdigital/dp-redis"
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

func Setup(ctx context.Context, r *mux.Router, permissions AuthHandler, elasticacheClient *dpredis.Client) *API {
	api := &API{
		Router: r,
	}

	r.HandleFunc("/sessions", permissions.Require(create, CreateSessionHandlerFunc(elasticacheClient))).Methods("POST")
	r.HandleFunc("/sessions/{ID}", GetByIDSessionHandlerFunc(elasticacheClient, mux.Vars)).Methods("GET")
	r.HandleFunc("/sessions/{Email}", GetByEmailSessionHandlerFunc(elasticacheClient, mux.Vars)).Methods("GET")
	r.HandleFunc("/sessions", permissions.Require(delete, DeleteAllSessionsHandlerFunc(elasticacheClient))).Methods("DELETE")
	return api
}

func (*API) Close(ctx context.Context) error {
	// Close any dependencies
	log.Event(ctx, "graceful shutdown of api complete", log.INFO)
	return nil
}
