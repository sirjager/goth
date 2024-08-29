package mw

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sirjager/gopkg/httpx"

	"github.com/sirjager/goth/entity"
	"github.com/sirjager/goth/core"
)

type (
	AuthProvder string
	contextType string
)

const (
	CookieSessionID     = "sessionId"
	CookieAccessToken   = "accessToken"
	CookieRefreshToken  = "refreshToken"
	CookieGothicSession = CookieAccessToken
)

const (
	AuthProviderOAuth       AuthProvder = "oauth"
	AuthProviderTokens      AuthProvder = "tokens"
	ContextAuthProvider     AuthProvder = "ctx_auth_provider"
	ContextKeyUser          contextType = "ctx_authenticated_user"
	ContextKeyAccessPayload contextType = "ctx_access_payload"
)

// RequiresAccessToken authenticates the request and adds the user to the context
func RequiresAccessToken(app *core.App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user, status, err := getAuthenticatedUser(r, app, CookieAccessToken)
			if err != nil {
				httpx.Error(w, err, status)
				return
			}
			ctx = context.WithValue(ctx, ContextKeyUser, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequiresRefreshToken authenticates the request and adds the user to the context
func RequiresRefreshToken(app *core.App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user, status, err := getAuthenticatedUser(r, app, CookieRefreshToken)
			if err != nil {
				httpx.Error(w, err, status)
				return
			}
			ctx = context.WithValue(ctx, ContextKeyUser, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserOrPanic assert that the user is authenticated, set by RequiresAuth middleware
func UserOrPanic(r *http.Request) *entity.User {
	user, ok := r.Context().Value(ContextKeyUser).(*entity.User)
	if !ok || user == nil || user.Email.Value() == "" {
		panic("authenticated operation, must have valid user")
	}
	return user
}

// AdminOrPanic assert that the user is a admin/master
func AdminOrPanic(r *http.Request) *entity.User {
	user := UserOrPanic(r)
	if !user.Master {
		panic("admin operation, must be a admin user")
	}
	return user
}

// IsCurrentUserIdentity returns if identity params matches authorized user identity
// It matches email,user_id, "me" with /{identity} params.
func IsCurrentUserIdentity(r *http.Request) bool {
	user := UserOrPanic(r)                  // authenticated user
	identity := chi.URLParam(r, "identity") //  /some-protected-path/<identity>
	currentUserIdentities := []string{"me", user.Email.Value(), user.ID.Value().String()}
	return valueExist(identity, currentUserIdentities)
}

// Generic function to check if a value exists in a slice
func valueExist[T comparable](find T, in []T) bool {
	for _, v := range in {
		if v == find {
			return true
		}
	}
	return false
}
