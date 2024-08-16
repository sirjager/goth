package mw

import (
	"net/http"
)

// RequiresMaster authenticates the request and adds the user to the context, and user has to be a master,
// This middlewares should be used after RequiresAuth
func RequiresMaster() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := UserOrPanic(r)
			if !user.Verified {
				http.Error(w, emailNotVerified, http.StatusForbidden)
				return
			}
			if !user.Master {
				http.Error(w, insufficientPermissions, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
