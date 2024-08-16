package mw

import (
	"net/http"

	"github.com/go-chi/cors"
)

func UseCors(opts ...cors.Options) func(next http.Handler) http.Handler {
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any major browsers
	}
	if len(opts) == 1 {
		corsOptions = opts[0]
	}
	return cors.Handler(corsOptions)
}
