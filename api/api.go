package api

import (
	"context"
	"github.com/ONSdigital/dp-authorisation/auth"
	dpredis "github.com/ONSdigital/dp-redis"
	"github.com/ONSdigital/dp-sessions-api/config"
	"github.com/ONSdigital/dp-sessions-api/session"
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

func Setup(ctx context.Context, r *mux.Router, permissions AuthHandler) *API {
	api := &API{
		Router: r,
	}

	s := session.NewSession()

	cfg, err := config.Get()
	if err != nil {
		return nil
	}
	
	cache, err := dpredis.NewClient(dpredis.Config{
		Addr:     cfg.ElasticacheAddr,
		Password: cfg.ElasticachePassword,
		Database: cfg.ElasticacheDatabase,
		TTL:      cfg.ElasticacheTTL,
	})
	if err != nil {
		return nil
	}

	r.HandleFunc("/sessions", permissions.Require(create, CreateSessionHandlerFunc(s, cache))).Methods("POST")
	r.HandleFunc("/sessions/{ID}", GetByIDSessionHandlerFunc(cache, mux.Vars)).Methods("GET")
	r.HandleFunc("/sessions", permissions.Require(delete, DeleteAllSessionsHandlerFunc(cache))).Methods("DELETE")
	return api
}

func (*API) Close(ctx context.Context) error {
	// Close any dependencies
	log.Event(ctx, "graceful shutdown of api complete", log.INFO)
	return nil
}
