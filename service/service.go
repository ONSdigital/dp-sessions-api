package service

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/zebedee"
	"github.com/ONSdigital/dp-authorisation/auth"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	rchttp "github.com/ONSdigital/dp-net/http"
	"github.com/ONSdigital/dp-sessions-api/api"
	"github.com/ONSdigital/dp-sessions-api/config"
	"github.com/ONSdigital/go-ns/server"
	"github.com/ONSdigital/log.go/log"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Service struct {
	Config      *config.Config
	server      *server.Server
	Router      *mux.Router
	API         *api.API
	HealthCheck *healthcheck.HealthCheck
}

// Run the service
func Run(buildTime, gitCommit, version string, svcErrors chan error) (*Service, error) {
	ctx := context.Background()
	log.Event(ctx, "running service", log.INFO)

	cfg, err := config.Get()
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve service configuration")
	}
	log.Event(ctx, "got service configuration", log.Data{"config": cfg}, log.INFO)

	r := mux.NewRouter()

	s := server.New(cfg.BindAddr, r)
	s.HandleOSSignals = false

	permissions := getAuthorisationHandlers(cfg)

	if cfg.EnableTrainingFlag {
		a := api.Setup(ctx, r, permissions)
		spew.Dump(a)
	}

	versionInfo, err := healthcheck.NewVersionInfo(
		buildTime,
		gitCommit,
		version,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse version information")
	}

	hc := healthcheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	zebedeeClient := zebedee.New(cfg.ZebedeeURL)

	if err := registerCheckers(ctx, &hc, zebedeeClient); err != nil {
		return nil, errors.Wrap(err, "unable to register checkers")
	}
	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)

	hc.Start(ctx)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			svcErrors <- errors.Wrap(err, "failure in http listen and serve")
		}
	}()

	return &Service{
		Config:      cfg,
		Router:      r,
		API:         a,
		HealthCheck: &hc,
		server:      s,
	}, nil
}

// Gracefully shutdown the service
func (svc *Service) Close(ctx context.Context) {
	timeout := svc.Config.GracefulShutdownTimeout
	log.Event(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout}, log.INFO)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// stop any incoming requests before closing any outbound connections
	if err := svc.server.Shutdown(ctx); err != nil {
		log.Event(ctx, "failed to shutdown http server", log.Error(err), log.ERROR)
	}

	if err := svc.API.Close(ctx); err != nil {
		log.Event(ctx, "error closing API", log.Error(err), log.ERROR)
	}

	log.Event(ctx, "graceful shutdown complete", log.INFO)
}

func registerCheckers(ctx context.Context, hc *healthcheck.HealthCheck, zebedeeClient *zebedee.Client) (err error) {
	hasErrors := false

	if err = hc.AddCheck("Zebedee", zebedeeClient.Checker); err != nil {
		hasErrors = true
		log.Event(ctx, "error adding check for zebedeee", log.ERROR, log.Error(err))
	}

	if hasErrors {
		return errors.New("Error(s) registering checkers for healthcheck")
	}
	return nil
}

func getAuthorisationHandlers(cfg *config.Config) api.AuthHandler {
	auth.LoggerNamespace("dp-sessions-api-auth")

	authClient := auth.NewPermissionsClient(rchttp.NewClient())
	authVerifier := auth.DefaultPermissionsVerifier()

	// for checking caller permissions when we only have a user/service token
	permissions := auth.NewHandler(
		auth.NewPermissionsRequestBuilder(cfg.ZebedeeURL),
		authClient,
		authVerifier,
	)

	return permissions
}
