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
func (a *Server) Signout(w http.ResponseWriter, r *http.Request) {
	// INFO: Clear Cookies
	a.SetCookies(w,
		&http.Cookie{
			Name: mw.CookieGothicSession, Value: "",
			Path: "/", Expires: time.Now().Add(-24*time.Hour),
			HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
		},
		&http.Cookie{
			Name: mw.CookieSessionID, Value: "",
			Path: "/", Expires: time.Now().Add(-24*time.Hour),
			HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
		},
		&http.Cookie{
			Name: mw.CookieAccessToken, Value: "",
			Path: "/", Expires: time.Now().Add(-24*time.Hour),
			HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
		},
		&http.Cookie{
			Name: mw.CookieRefreshToken, Value: "",
			Path: "/", Expires: time.Now().Add(-24*time.Hour),
			HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
		},
	)
	_ = gothic.Logout(w, r)

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
