package api

//go:generate moq -out mockauth_test.go . AuthHandler

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-authorisation/auth"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
	"net/http"
)

//API provides a struct to wrap the api around
type API struct {
	Router *mux.Router
}

type AuthHandler interface {
	Require(required auth.Permissions, handler http.HandlerFunc) http.HandlerFunc
}

func Setup(ctx context.Context, r *mux.Router, permissions AuthHandler) *API {
	api := &API{
		Router: r,
	}

	nopSess := &NOPSessions{}
	nopCache := &NOPCache{}
	create := auth.Permissions{Create: true}

	r.HandleFunc("/session", permissions.Require(create, CreateSessionHandlerFunc(nopSess, nopCache))).Methods("POST")
	r.HandleFunc("/session/{ID}", GetByIDSessionHandlerFunc(nopCache, mux.Vars)).Methods("GET")
	return api
}

func (*API) Close(ctx context.Context) error {
	// Close any dependencies
	log.Event(ctx, "graceful shutdown of api complete", log.INFO)
	return nil
}
