package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/sirjager/gopkg/cache"
	"github.com/sirjager/gopkg/tokens"
	"golang.org/x/sync/errgroup"

	"github.com/sirjager/goth/config"
	mw "github.com/sirjager/goth/middlewares"
	"github.com/sirjager/goth/repository"
	"github.com/sirjager/goth/worker"
)

type Server struct {
	logr     zerolog.Logger
	cache    cache.Cache
	toknb    tokens.TokenBuilder // tokens was too common and was conflicting in code, so using toknb for now.
	router   *chi.Mux
	repo     repository.Repo
	validate *validator.Validate
	config   *config.Config
	tasks    worker.TaskDistributor
	adapters *mw.Adapters
}

func NewServer(a *mw.Adapters) *Server {
	server := &Server{
		logr:     a.Logr,
		config:   a.Config,
		repo:     a.Repo,
		validate: a.Validate,
		toknb:    a.Toknb,
		cache:    a.Cache,
		tasks:    a.Tasks,
		adapters: a,
	}
	server.MountHandlers()
	return server
}

func (server *Server) StartServer(address string, ctx context.Context, wg *errgroup.Group) {
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
