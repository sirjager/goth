package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"

	mw "github.com/sirjager/goth/middlewares"
)

// @Summary		Login
// @Description	Authenticates a user with a specified provider
// @Tags			Auth
// @Produce		json
// @Param			provider	path		string			true	"Provider Name"
// @Success		200			{object}	UserResponse	"User Response"
// @Router			/auth/{provider} [get]
func (a *API) AuthProvider(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	if user, loggedIn := mw.IsAuthenticated(r, a.repo); loggedIn {
		response := UserResponse{User: EntityToUser(user)}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	req := r.WithContext(context.WithValue(r.Context(), "provider", provider))
	gothic.BeginAuthHandler(w, req)
}
