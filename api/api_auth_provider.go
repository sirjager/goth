package api

import (
	"context"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"

	mw "github.com/sirjager/goth/middlewares"
)

// @Summary		Login
// @Description	Authenticates a user with a specified provider
// @Tags			Auth
// @Produce		json
// @Router			/auth/{provider} [get]
// @Param			provider	path		string			true	"Provider Name"
// @Success		200			{object}	UserResponse	"User object"
func (a *API) AuthProvider(w http.ResponseWriter, r *http.Request) {
	refererURL := r.Header.Get("Referer")
	parsedURL, err := url.Parse(refererURL)
	if err != nil {
		http.Error(w, "invalid refer url", http.StatusBadRequest)
		return
	}
	// Reconstruct the base URL
	refererURL = parsedURL.Scheme + "://" + parsedURL.Host

	provider := chi.URLParam(r, "provider")
	if user, authenticated := mw.IsAuthenticated(r, a.repo, a.tokens, a.cache); authenticated {
		response := UserResponse{User: user.Profile()}
		a.Success(w, response)
		return
	}

	req := r.WithContext(context.WithValue(r.Context(), "provider", provider))

	gothic.SetState = func(req *http.Request) string {
		return refererURL
	}

	gothic.BeginAuthHandler(w, req)
}
