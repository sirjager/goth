package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"

	"github.com/sirjager/goth/modules"
)

type Server struct {
	router *chi.Mux
	*modules.Modules
}

func NewServer(modules *modules.Modules) *Server {
	server := &Server{Modules: modules}
	server.MountHandlers()
	return server
}

func (server *Server) StartServer(address string, ctx context.Context, wg *errgroup.Group) {
	httpServer := &http.Server{Handler: server.router, Addr: address}
	wg.Go(func() error {
		server.Logger().Info().Msgf("started http server at %s", address)
		if err := httpServer.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				server.Logger().Error().Err(err).Msg("failed to start http server")
				return err
			}
		}
		return nil
	})

	wg.Go(func() error {
		<-ctx.Done()
		server.Logger().Info().Msg("gracefully shutting down http server")
		// NOTE: here we can limit maximum time for graceful shutdown
		// but for now we do not need it, we can use context.Background()
		if err := httpServer.Shutdown(context.Background()); err != nil {
			server.Logger().Error().Err(err).Msg("failed to shutdown http server")
			return err
		}
		server.Logger().Info().Msg("http server gracefully shutdown complete")

		return nil
	})
}
