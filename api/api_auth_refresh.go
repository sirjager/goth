package api

import (
	"net/http"

	"github.com/sirjager/gopkg/utils"

	mw "github.com/sirjager/goth/middlewares"
)

func (s *Server) RefreshToken(w http.ResponseWriter, r *http.Request) {
	user := mw.UserOrPanic(r)

	sessionID := utils.XIDNew().String()
	accessData := mw.NewAccessPayload(user.ID, sessionID)
	accessTokenDur := s.config.AuthAccessTokenExpire
	accessToken, accessTokenPayload, err := s.tokens.CreateToken(accessData, accessTokenDur)
	if err != nil {
		s.Failure(w, err)
		return
	}

	accessKey := mw.TokenKey(user.ID.Value().String(), sessionID, mw.TokenTypeAccess)
	if err = s.cache.Set(r.Context(), accessKey, accessTokenPayload, accessTokenDur); err != nil {
		s.Failure(w, err)
		return
	}

	s.SetCookies(w,
		&http.Cookie{
			Name: mw.CookieSessionID, Value: sessionID,
			Path: "/", Expires: accessTokenPayload.ExpiresAt,
			HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
		},
		&http.Cookie{
			Name: mw.CookieAccessToken, Value: accessToken,
			Path: "/", Expires: accessTokenPayload.ExpiresAt,
			HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
		},
	)

	s.SuccessOK(w, "access tokens refreshed")
}
