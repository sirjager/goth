package mw

import (
	"context"
	"net/http"

	"github.com/markbates/goth/gothic"

	"github.com/sirjager/goth/entity"
	"github.com/sirjager/goth/repository"
	"github.com/sirjager/goth/vo"
)

type contextType string

const SessionCookieName = "gothic_session"

const ContextKeyUser contextType = "ctx_authenticated_user"

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
	result := repo.UserReadByEmail(r.Context(), email.Value())
	if result.Error != nil {
		return nil, false
	}
	return result.User, true
}

func RequiresAuth(repo *repository.Repo) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, loggedIn := IsAuthenticated(r, repo)
			if !loggedIn {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ContextKeyUser, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
