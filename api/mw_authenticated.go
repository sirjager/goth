package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"

	repoerrors "github.com/sirjager/goth/repository/errors"
)

type contextType string
const userContext contextType = "user"

func (a *API) RequiresAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := gothic.Store.Get(r, SessionName)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if session.Values["user"] == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		u := session.Values["user"].(goth.User)

		// fetch user from database
		user, err := a.repo.UserReadByEmail(r.Context(), u.Email)
		if err != nil {
			if errors.Is(err, repoerrors.ErrUserNotFound) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), userContext, &user)

		//
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
