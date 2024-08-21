package mw

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirjager/gopkg/utils"

	"github.com/sirjager/goth/entity"
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
	if !ok || user == nil || user.Email.Value() == "" {
		panic("authenticated operation, must have valid user")
	}
	return user
}

// IsCurrentUserIdentity returns if identity params matches authorized user identity
// It matches email,user_id, "me" with /{identity} params.
func IsCurrentUserIdentity(r *http.Request) bool {
	user := UserOrPanic(r)
	identity := chi.URLParam(r, "identity")
	currentUserIdentities := []string{"me", user.Email.Value(), user.ID.Value().String()}
	return utils.ValueExist(identity, currentUserIdentities)
}
