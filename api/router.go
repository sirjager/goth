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

	c.Get("/", a.welcome)
	c.Get("/health", a.health)
	c.Get("/swagger", a.scalarDocs)

	// authentication routes
	c.Route("/auth", func(r chi.Router) {
		r.Get("/{provider}", a.AuthProvider)
		r.Get("/{provider}/callback", a.AuthCallback)
		r.Get("/logout/{provider}", a.AuthLogout)
	})

	// NOTE: Authenticated Routes
	c.Group(func(r chi.Router) {
		r.Use(mw.RequiresAuth(a.repo))

		r.Get("/users", a.UsersGet)
		r.Get("/users/{identity}", a.UserGet)
		r.Patch("/users/{identity}", a.UserUpdate)
	})
}
