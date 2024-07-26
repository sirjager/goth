package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/sirjager/goth/config"
	"github.com/sirjager/goth/repository"
)

type API struct {
	logr     zerolog.Logger
	router   *chi.Mux
	repo     *repository.Repo
	validate *validator.Validate
	config   config.Config
}

func NewServer(repo *repository.Repo, logr zerolog.Logger, config config.Config) *API {
	validator := validator.New(validator.WithRequiredStructEnabled())
	server := &API{
		logr:      logr,
		config:    config,
		repo:      repo,
		validate: validator,
	}
	server.setupRouter()
	return server
}

func (server *API) StartServer(address string, ctx context.Context, wg *errgroup.Group) {
	httpServer := &http.Server{Handler: server.router, Addr: address}
	wg.Go(func() error {
		server.logr.Info().Msgf("started http server at %s", address)
		if err := httpServer.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				server.logr.Error().Err(err).Msg("failed to start http server")
				return err
			}
		}
		return nil
	})

	wg.Go(func() error {
		<-ctx.Done()
		server.logr.Info().Msg("gracefully shutting down http server")
		// NOTE: here we can limit maximum time for graceful shutdown
		// but for now we do not need it, we can use context.Background()
		if err := httpServer.Shutdown(context.Background()); err != nil {
			server.logr.Error().Err(err).Msg("failed to shutdown http server")
			return err
		}
		server.logr.Info().Msg("http server gracefully shutdown complete")

		return nil
	})
}
