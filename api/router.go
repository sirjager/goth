package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	mw "github.com/sirjager/goth/middlewares"
)

func (a *API) setupRouter() {
	c := chi.NewRouter()
	defer func() { a.router = c }()

	c.Use(middleware.Recoverer)

	c.Get("/", a.SysWelcome)
	c.Get("/health", a.SysHealth)
	c.Get("/swagger", a.SysDocs)

	// NOTE: Authentication routes
	c.Route("/auth", func(r chi.Router) {
		r.Get("/{provider}", a.AuthProvider)
		r.Get("/{provider}/callback", a.AuthCallback)
		r.Get("/logout/{provider}", a.AuthLogout)
	})

	// NOTE: Authenticated routes
	c.Group(func(authenticated chi.Router) {
		authenticated.Use(mw.RequiresAuth(a.repo))
		authenticated.Use(mw.AquireRoles(a.repo))
		authenticated.Use(mw.AquirePermissions(a.repo))

		authenticated.Get("/users", a.UsersGet)
		authenticated.Get("/users/{identity}", a.UserGet)
		authenticated.Patch("/users/{identity}", a.UserUpdate)
	})
}
