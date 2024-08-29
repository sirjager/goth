package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rakyll/statik/fs"

	mw "github.com/sirjager/goth/middlewares"
	"github.com/sirjager/goth/statik/docs"
)

func (s *Server) MountHandlers() {
	mux := chi.NewRouter()
	defer func() { s.router = mux }()

	mux.Use(middleware.RealIP)
	mux.Use(mw.UseCors())
	mux.Use(mw.Logger(*s.Logger(), s.Config().LoggerLogfile))
	mux.Use(middleware.Compress(5))
	mux.Use(middleware.Recoverer)
	mux.Use(mw.RequestID())

	mux.Route("/api", func(r chi.Router) {
		r.Get("/", s.apiWelcome)
		r.Get("/health", s.apiHealth)
		r.Route("/docs", s.docsHandler)
		r.Route("/auth", s.authHandlers)
		r.Route("/admin", s.adminHandlers)
	})
}

func (s *Server) docsHandler(r chi.Router) {
	r.Get("/", s.swaggerDocs)
	docsFS, err := fs.NewWithNamespace(docs.Docs)
	if err != nil {
		s.Logger().Fatal().Err(err).Msg("can not statik file server")
	}
	swaggerHandler := http.StripPrefix("/api/docs/", http.FileServer(docsFS))
	r.Handle("/swagger.json", swaggerHandler)
}

func (s *Server) authHandlers(r chi.Router) {
	r.Get("/signin", s.authUserSignIn)
	r.Post("/signup", s.authUserSignUp)
	r.Get("/verify", s.authUserVerify)
	r.Post("/reset", s.authUserResetPassword)

	r.With(mw.RequiresAccessToken(s.App), mw.RequiresVerified()).
		Get("/user", s.authUserFetch)

	r.With(mw.RequiresAccessToken(s.App), mw.RequiresVerified()).
		Patch("/user", s.authUserUpdate)

	r.With(mw.RequiresAccessToken(s.App), mw.RequiresVerified()).
		Get("/delete", s.authUserDelete)

	r.With(mw.RequiresRefreshToken(s.App), mw.RequiresVerified()).
		Get("/refresh", s.authUserRefreshToken)

	r.Get("/signout/{provider}", s.authUserSignOut)
	r.Get("/{provider}", s.oauthProvider)
	r.Get("/{provider}/callback", s.oauthCallback)
}

func (s *Server) adminHandlers(r chi.Router) {
	r.Use(mw.RequiresAccessToken(s.App))
	r.Use(mw.RequiresVerified())
	r.Use(mw.RequiresMaster())
	r.Get("/users", s.adminFetchUsers)
	r.Patch("/users/{identity}", s.adminUpdateUser)
}
