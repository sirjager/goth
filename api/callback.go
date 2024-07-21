package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"

	"github.com/sirjager/goth/entity"
	repoerrors "github.com/sirjager/goth/repository/errors"
)

func (a *API) providerCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	redirect := chi.URLParam(r, "redirect")
	fmt.Println(redirect)

	req := r.WithContext(context.WithValue(r.Context(), "provider", provider))

	gothUser, err := gothic.CompleteUserAuth(w, req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// NOTE: saving user in database
	// IF EXISTS      : fetch from database using email, and return it
	// IF NOT EXISTS  : create and save user object, and return it
	var user *entity.User
	user, err = a.repo.UserCreate(r.Context(), GothUserToEntityUser(gothUser))
	if err != nil {
		// if its not, user already exits error, return it
		if !errors.Is(err, repoerrors.ErrUserAlreadyExists) {
			http.Error(w, err.Error(), 500)
			return
		}

		// if its user already exits error, get user and return it
		user, err = a.repo.UserReadByEmail(r.Context(), gothUser.Email)
	}

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// store user in cookie
	if err := a.StoreUserSession(w, r, gothUser); err != nil {
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	data := UserResponse{User: EntityToUser(user)}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
