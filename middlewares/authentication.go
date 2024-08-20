package mw

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirjager/gopkg/cache"
	"github.com/sirjager/gopkg/tokens"
	"github.com/sirjager/gopkg/utils"

	"github.com/sirjager/goth/entity"
	"github.com/sirjager/goth/repository"
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
	AuthProviderOAuth   AuthProvder = "oauth"
	AuthProviderTokens  AuthProvder = "tokens"
	ContextAuthProvider AuthProvder = "ctx_auth_provider"
	ContextKeyUser      contextType = "ctx_authenticated_user"
)

// UserOrPanic assert that the user is authenticated, set by RequiresAuth middleware
func UserOrPanic(r *http.Request) *entity.User {
	user, ok := r.Context().Value(ContextKeyUser).(*entity.User)
	if !ok || user == nil {
		panic("authenticated operation, must have valid user")
	}
	return user
}

// RequiresAuth authenticates the request and adds the user to the context
func RequiresAuth(
	repo repository.Repository,
	tokens tokens.TokenBuilder,
	cache cache.Cache,
	allowRefreshToken ...bool,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var err error
			var user *entity.User
			provider := AuthProviderOAuth // provider used to authenticate
			user, err = authenticateUsingOAuth(r, repo)
			if err != nil {
				provider = AuthProviderTokens
				user, err = authenticateUsingTokens(r, repo, cache, tokens, allowRefreshToken...)
			}

			if err != nil {
				http.Error(w, unauthorized, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyUser, user)
			ctx = context.WithValue(ctx, ContextAuthProvider, provider)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// IsAuthenticated authenticates the request and returns `user` and `boolean` value.
// `true if logged in` else `false if not logged in`
func IsAuthenticated(
	r *http.Request,
	repo repository.Repository,
	tokens tokens.TokenBuilder,
	cache cache.Cache,
) (*entity.User, bool) {
	var err error
	var user *entity.User
	user, err = authenticateUsingOAuth(r, repo)
	if err != nil {
		user, err = authenticateUsingTokens(r, repo, cache, tokens)
	}
	if err != nil || user == nil {
		return nil, false
	}
	return user, true
}

// IsCurrentUserIdentity returns if identity params matches authorized user identity
// It matches email,user_id, "me" with /{identity} params.
func IsCurrentUserIdentity(r *http.Request) bool {
	user := UserOrPanic(r)
	identity := chi.URLParam(r, "identity")
	currentUserIdentities := []string{"me", user.Email.Value(), user.ID.Value().String()}
	return utils.ValueExist(identity, currentUserIdentities)
}
