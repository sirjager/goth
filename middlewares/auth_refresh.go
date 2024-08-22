package mw

import (
	"context"
	"errors"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/sirjager/gopkg/cache"
	"github.com/sirjager/gopkg/tokens"

	"github.com/sirjager/goth/payload"
	"github.com/sirjager/goth/repository"
	"github.com/sirjager/goth/vo"
)

// RequiresAccessToken authenticates the request and adds the user to the context
func RequiresRefreshToken(
	repo repository.Repo,
	builder tokens.TokenBuilder,
	sessions cache.Cache,
	logr zerolog.Logger,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			var err error
			token := extractAuthToken(r, CookieRefreshToken)
			var incoming payload.RefreshToken
			if _, err = builder.VerifyToken(token, &incoming); err != nil {
				logr.Error().Msg("refresh token validation failed")
				http.Error(w, unauthorized, http.StatusUnauthorized)
				return
			}

			var tokenPayload tokens.Payload
			tokenKey := payload.SessionRefreshKey(incoming.UserID, incoming.SessionID)
			if err = sessions.Get(ctx, tokenKey, &tokenPayload); err != nil {
				if !errors.Is(err, cache.ErrNoRecord) {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				logr.Error().Msg("session not stored")
				http.Error(w, unauthorized, http.StatusUnauthorized)
				return
			}

			var stored payload.RefreshToken
			if err = builder.ReadPayload(&tokenPayload, &stored); err != nil {
				logr.Error().Msg("failed to read stored payload")
				http.Error(w, unauthorized, http.StatusUnauthorized)
				return
			}

			userID, err := vo.NewIDFrom(incoming.UserID)
			if err != nil {
				logr.Error().Msg("invalid user id")
				http.Error(w, unauthorized, http.StatusUnauthorized)
				return
			}

			if !stored.IsEqual(incoming) {
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
			if res.User == nil {
				logr.Error().Msg("user object is nil")
				http.Error(w, unauthorized, http.StatusUnauthorized)
				return
			}
			ctx = context.WithValue(ctx, ContextKeyUser, res.User)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
