package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
		r.Get("/{provider}", a.authProvider)
		r.Get("/{provider}/callback", a.providerCallback)
		r.Get("/{provider}/logout", a.logoutProvider)
	})

	// authenticated routes
	c.Group(func(r chi.Router) {
		r.Use(a.RequiresAuth)
		r.Get("/users", a.getUser)
	})
}
