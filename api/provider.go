package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"
)

// @Summary		Login
// @Description	Authenticates a user with a specified provider
// @Tags			Auth
// @Produce		json
// @Param			provider	path		string			true	"Provider Name"
// @Success		200			{object}	UserResponse	"User Response"
// @Router			/auth/{provider} [get]
func (a *API) authProvider(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	req := r.WithContext(context.WithValue(r.Context(), "provider", provider))

	gothuser, err := gothic.CompleteUserAuth(w, req)
	if err != nil {
		gothic.BeginAuthHandler(w, req)
		return
	}

	user := UserResponse{User: GothUserToUser(gothuser)}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
