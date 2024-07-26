package api

import (
	"net/http"

	"github.com/markbates/goth/gothic"
)

// @Summary		Logout
// @Description	Logout from a provider
// @Tags			Auth
// @Produce		json
// @Param			provider	path	string	true	"Provider Name"
// @Router			/auth/logout/{provider} [get]
func (a *API) AuthLogout(w http.ResponseWriter, r *http.Request) {
	if err := gothic.Logout(w, r); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
