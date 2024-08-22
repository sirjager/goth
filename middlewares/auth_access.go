package mw

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
	"github.com/sirjager/gopkg/cache"
	"github.com/sirjager/gopkg/tokens"

	"github.com/sirjager/goth/entity"
	"github.com/sirjager/goth/payload"
	"github.com/sirjager/goth/repository"
	"github.com/sirjager/goth/vo"
)

// RequiresAccessToken authenticates the request and adds the user to the context
func RequiresAccessToken(
	repo repository.Repo,
	builder tokens.TokenBuilder,
	sessions cache.Cache,
	logr zerolog.Logger,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			var err error
			var user *entity.User
			user, err = authenticateUsingOAuth(r, repo)
			if err != nil || user == nil {

				accessToken := extractAuthToken(r, CookieAccessToken)
				var accessData payload.AccessToken
				if _, err = builder.VerifyToken(accessToken, &accessData); err != nil {
					logr.Error().Msg("access token verification failed")
					http.Error(w, unauthorized, http.StatusUnauthorized)
					return
				}

				var accessPayload tokens.Payload
				accessKey := payload.SessionAccessKey(accessData.UserID, accessData.SessionID)
				if err = sessions.Get(ctx, accessKey, &accessPayload); err != nil {
					if !errors.Is(err, cache.ErrNoRecord) {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					logr.Error().Msg("session not stored")
					http.Error(w, unauthorized, http.StatusUnauthorized)
					return
				}

				var _storedData payload.AccessToken
				if err = builder.ReadPayload(&accessPayload, &_storedData); err != nil {
					logr.Error().Msg("failed to read stored payload")
					http.Error(w, unauthorized, http.StatusUnauthorized)
					return
				}

				userID, err := vo.NewIDFrom(accessData.UserID)
				if err != nil {
					logr.Error().Msg("invalid user id")
					http.Error(w, unauthorized, http.StatusUnauthorized)
					return
				}

				if !_storedData.IsEqual(accessData) {
					logr.Error().Msg("incoming access payload != stored access payload")
					http.Error(w, unauthorized, http.StatusUnauthorized)
					return
				}

				res := repo.UserGetByID(ctx, userID)
				if res.Error != nil {
					if res.StatusCode == http.StatusNotFound {
						logr.Error().Msg("user not found")
						http.Error(w, unauthorized, http.StatusUnauthorized)
						return
					}
					http.Error(w, res.Error.Error(), res.StatusCode)
					return
				}
				user = res.User
			}
			if user == nil {
				logr.Error().Msg("user object is nil")
				http.Error(w, unauthorized, http.StatusUnauthorized)
				return
			}
			ctx = context.WithValue(ctx, ContextKeyUser, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractAuthToken(r *http.Request, cookieName string) (token string) {
	if cookie, cookieErr := r.Cookie(cookieName); cookieErr == nil {
		token = cookie.Value
	}
	// If access token is already in cookies return it
	if token != "" {
		return
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return
	}
	values := strings.Split(authHeader, " ")
	if len(values) != 2 {
		return
	}
	if values[0] != "Bearer" {
		return
	}
	token = values[1]
	return
}
