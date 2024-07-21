package mw

import (
	"context"
	"errors"
	"net/http"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"

	"github.com/sirjager/goth/repository"
	repoerrors "github.com/sirjager/goth/repository/errors"
)

type contextType string

const UserContext contextType = "user"

const SessionName = "session"

// StoreUserSession stores the user in the cookies
func StoreUserSession(w http.ResponseWriter, r *http.Request, user goth.User) error {
	session, _ := gothic.Store.Get(r, SessionName)
	session.Values["user"] = user

	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}

func RequiresAuthententicated(repo repository.Repo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// extracts session request cookie and fetch session from store
			session, err := gothic.Store.Get(r, SessionName)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// extract goth user object from session values
			userValues, okValues := session.Values["user"]
			if !okValues || userValues == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
 
			// type cast goth user object
			u, okUser := userValues.(goth.User)
			if !okUser {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// fetch user from database
			user, err := repo.UserReadByEmail(r.Context(), u.Email)
			if err != nil {
				if errors.Is(err, repoerrors.ErrUserNotFound) {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), UserContext, &user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
