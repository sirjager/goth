package mw

import (
	"net/http"

	"github.com/sirjager/gopkg/httpx"
)

// RequiresMaster authenticates the request and adds the user to the context, and user has to be a master,
// This middlewares should be used after RequiresAuth
func RequiresMaster() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := UserOrPanic(r)
			if !user.Verified {
				httpx.Error(w, ErrEmailNotVerified, http.StatusForbidden)
				return
			}
			if !user.Master {
				httpx.Error(w, ErrInsufficientPermissions, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
