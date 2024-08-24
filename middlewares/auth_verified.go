package mw

import (
	"net/http"

	"github.com/sirjager/gopkg/httpx"
)

// RequiresVerified check if authenticated user if verified or not. if not it rejects the requsets
// This middlewares should be used after RequiresAuth
func RequiresVerified() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := UserOrPanic(r)
			if !user.Verified {
				httpx.Error(w, ErrEmailNotVerified, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
