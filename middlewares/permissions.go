package mw

import (
	"net/http"
)

// RequiresPermissions ensures that user requesting resource has permissions to access it.
// User with master role can access this route
// User with same identity can access this route.
func RequiresPermissions() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowRequest := false
			user := UserOrPanic(r)
			if IsCurrentUserIdentity(r) {
				allowRequest = true
			}
			if user.Master {
				allowRequest = true
			}
			if !allowRequest {
				http.Error(w, insufficientPermissions, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
