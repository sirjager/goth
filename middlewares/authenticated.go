package mw

import (
	"context"
	"net/http"

	"github.com/markbates/goth/gothic"

	"github.com/sirjager/goth/entity"
	"github.com/sirjager/goth/repository"
	"github.com/sirjager/goth/vo"
)

const (
	emailNotVerified        = "Email not verified"
	insufficientPermissions = "Insufficient permissions"
	unauthorized            = "Unauthorized"
)

type contextType string

const SessionCookieName = "gothic_session"

const ContextKeyUser contextType = "ctx_authenticated_user"


// UserOrPanic assert that the user is authenticated, set by RequiresAuth middleware
func UserOrPanic(r *http.Request) *entity.User {
	user, ok := r.Context().Value(ContextKeyUser).(*entity.User)
	if !ok || user == nil {
		panic("authenticated operation, must have valid user")
	}
	return user
}

// RequiresAuth authenticates the request and adds the user to the context
func RequiresAuth(repo *repository.Repo) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, loggedIn := IsAuthenticated(r, repo)
			if !loggedIn {
				http.Error(w, unauthorized, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ContextKeyUser, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

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

// RequiresVerified check if authenticated user if verified or not. if not it rejects the requsets
// This middlewares should be used after RequiresAuth
func RequiresVerified() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := UserOrPanic(r)
			if !user.Verified {
				http.Error(w, emailNotVerified, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// IsAuthenticated authenticates the request and returns `user` and `boolean` value.
// `true if logged in` else `false if not logged in`
func IsAuthenticated(r *http.Request, repo *repository.Repo) (*entity.User, bool) {
	sess, err := gothic.Store.Get(r, SessionCookieName)
	if err != nil {
		return nil, false
	}
	value, ok := sess.Values["email"].(string)
	if !ok || value == "" {
		return nil, false
	}
	email, err := vo.NewEmail(value)
	if err != nil {
		return nil, false
	}
	result := repo.UserReadByEmail(r.Context(), email)
	if result.Error != nil {
		return nil, false
	}
	return result.User, true
}
