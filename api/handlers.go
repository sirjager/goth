package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rakyll/statik/fs"

	mw "github.com/sirjager/goth/middlewares"
	_ "github.com/sirjager/goth/statik"
)

func (s *Server) MountHandlers() {
	c := chi.NewRouter()
	defer func() { s.router = c }()

	c.Use(middleware.RealIP)
	c.Use(mw.UseCors())
	c.Use(mw.Logger(s.Modules))
	c.Use(middleware.Compress(5))
	c.Use(middleware.Recoverer)
	c.Use(mw.RequestID(s.Modules))

	c.Get("/", s.Welcome)
	c.Get("/health", s.Health)
	c.Get(docsPath, s.SwaggerDocs)

	// NOTE: Authentication routes
	c.Route("/auth", func(r chi.Router) {
		r.Get("/signin", s.Signin)
		r.Post("/signup", s.Signup)
		r.Get("/verify", s.VerifyEmail)
		r.Post("/reset", s.Reset)

		r.With(mw.RequiresAccessToken(s.Modules), mw.RequiresVerified()).
			Get("/user", s.AuthUser)

		r.With(mw.RequiresAccessToken(s.Modules), mw.RequiresVerified()).
			Get("/delete", s.Delete)

		r.With(mw.RequiresRefreshToken(s.Modules), mw.RequiresVerified()).
			Get("/refresh", s.RefreshToken)

		r.Get("/signout/{provider}", s.Signout)
		r.Get("/{provider}", s.OAuthProvider)
		r.Get("/{provider}/callback", s.OAuthCallback)
	})

	c.Route("/users", func(r chi.Router) {
		r.Use(mw.RequiresAccessToken(s.Modules))
		r.Use(mw.RequiresVerified())
		r.Use(mw.RequiresMaster())
		r.Get("/", s.UsersGet)
	})

	c.Route("/users/{identity}", func(r chi.Router) {
		r.Use(mw.RequiresAccessToken(s.Modules))
		r.Use(mw.RequiresVerified())
		r.Use(mw.RequiresPermissions())

		r.Get("/", s.UserGet)
		r.Patch("/", s.UserUpdate)
	})

	// file server for swagger documentations, it serves only swagger.json
	// that will be used by scalar docs
	statikFS, err := fs.New()
	if err != nil {
		s.Logger().Fatal().Err(err).Msg("can not statik file server")
	}
	swaggerHandler := http.StripPrefix("/docs/", http.FileServer(statikFS))
	c.Get("/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		swaggerHandler.ServeHTTP(w, r)
	})
}
