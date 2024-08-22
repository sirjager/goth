package api

import (
	"net/http"

	"github.com/sirjager/gopkg/utils"

	"github.com/sirjager/goth/entity"
	mw "github.com/sirjager/goth/middlewares"
	"github.com/sirjager/goth/payload"
)

type RefreshTokenResponse struct {
	Message     string          `json:"message,omitempty"`
	User        *entity.Profile `json:"user,omitempty"`
	AccessToken string          `json:"accessToken,omitempty"`
	SessionID   string          `json:"sessionID,omitempty"`
} //	@name	RefreshTokenResponse

// Refresh Tokens
//
//	@Summary		Refresh
//	@Description	Refreshes Access Token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Router			/auth/refresh [get]
//	@Param			user	query		bool					false	"If true, returns User in body"
//	@Param			cookies	query		bool					false	"If true, returns AccessToken and SessionID in body"
//	@Success		200		{object}	RefreshTokenResponse	"RefreshTokenResponse"
func (s *Server) RefreshToken(w http.ResponseWriter, r *http.Request) {
	user := mw.UserOrPanic(r)
	sessionID := utils.XIDNew().String()
	accessData := payload.NewAccessPayload(user, sessionID)
	accessTokenDur := s.config.AuthAccessTokenExpire
	accessToken, accessTokenPayload, err := s.toknb.CreateToken(accessData, accessTokenDur)
	if err != nil {
		s.Failure(w, err)
		return
	}

	accessKey := payload.SessionAccessKey(user.ID.Value().String(), sessionID)
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

	userParam := r.URL.Query().Get("user")
	getUser := userParam == "true" || (r.URL.Query().Has("user") && userParam == "")

	cookiesParams := r.URL.Query().Get("cookies")
	getCookies := cookiesParams == "true" ||
		(r.URL.Query().Has("cookies") && cookiesParams == "")

	response := RefreshTokenResponse{Message: "access tokens refreshed"}

	if getUser {
		response.User = user.Profile()
	}

	if getCookies {
		response.SessionID = sessionID
		response.AccessToken = accessToken
	}

	s.Success(w, response)
}
