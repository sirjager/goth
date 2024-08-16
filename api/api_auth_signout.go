package api

import (
	"net/http"
	"time"

	"github.com/markbates/goth/gothic"

	mw "github.com/sirjager/goth/middlewares"
)

// @Summary		Signout
// @Description	Signout from a provider
// @Tags			Auth
// @Produce		json
// @Param			provider	path	string	true	"Provider Name"
// @Router			/auth/signout/{provider} [get]
func (a *API) Signout(w http.ResponseWriter, r *http.Request) {
	// INFO: Clear Cookies
	a.SetCookies(w,
		&http.Cookie{
			Name: mw.SessionCookieName, Value: "",
			Path: "/", Expires: time.Now(),
			HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
		},
		&http.Cookie{
			Name: "sessionId", Value: "",
			Path: "/", Expires: time.Now(),
			HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
		},
		&http.Cookie{
			Name: "accessToken", Value: "",
			Path: "/", Expires: time.Now(),
			HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
		},
		&http.Cookie{
			Name: "refreshToken", Value: "",
			Path: "/", Expires: time.Now(),
			HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
		},
	)
	_ = gothic.Logout(w, r)

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
