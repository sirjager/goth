package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	mw "github.com/sirjager/goth/middlewares"
)

func (a *API) setupRouter() {
	c := chi.NewRouter()
	defer func() { a.router = c }()

	// c.Use(mw.Logger(a.logr, a.config.Logger))
	c.Use(mw.UseCors())
	c.Use(mw.RequestID())
	c.Use(middleware.Compress(5))
	c.Use(middleware.RealIP)
	c.Use(middleware.Recoverer)

	c.Get("/", a.SysWelcome)
	c.Get("/health", a.SysHealth)
	c.Get("/swagger", a.SysDocs)

	// NOTE: Authentication routes
	c.Route("/auth", func(r chi.Router) {
		r.With(mw.RequiresAuth(a.repo, a.tokens, a.cache)).Get("/user", a.AuthUser)
		r.Get("/signin", a.Signin)
		r.Post("/signup", a.Signup)
		r.Get("/signout/{provider}", a.Signout)

		r.Get("/{provider}", a.AuthProvider)
		r.Get("/{provider}/callback", a.AuthCallback)
	})

	c.Route("/users", func(r chi.Router) {
		r.Use(mw.RequiresAuth(a.repo, a.tokens, a.cache))
		r.Use(mw.RequiresVerified())
		r.Use(mw.RequiresMaster())
		r.Get("/", a.UsersGet)
	})

	c.Route("/users/{identity}", func(r chi.Router) {
		r.Use(mw.RequiresAuth(a.repo, a.tokens, a.cache))
		r.Use(mw.RequiresVerified())
		r.Use(mw.RequiresPermissions())

		r.Get("/", a.UserGet)
		r.Patch("/", a.UserUpdate)
	})
}
