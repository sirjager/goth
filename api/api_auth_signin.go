package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/sirjager/gopkg/utils"

	"github.com/sirjager/goth/entity"
	mw "github.com/sirjager/goth/middlewares"
	"github.com/sirjager/goth/repository/users"
	"github.com/sirjager/goth/vo"
)

type SignInResponse struct {
	User         *entity.Profile `json:"user,omitempty"`
	AccessToken  string          `json:"accessToken,omitempty"`
	RefreshToken string          `json:"refreshToken,omitempty"`
	SessionID    string          `json:"sessionID,omitempty"`
	Message      string          `json:"message,omitempty"`
} //	@name	SignInResponse

var errInvalidCredentials = errors.New("invalid credentials")

// Signin Request
//
//	@Summary		Signin
//	@Description	Signin using credentials
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Router			/auth/signin [get]
//	@Security		BasicAuth
//	@Param			user	query		bool			false	"If true, returns User in body"
//	@Param			cookies	query		bool			false	"If true, returns AccessToken, RefreshToken and SessionID in body"
//	@Success		200		{object}	SignInResponse	"SignInResponse"
func (s *Server) Signin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	identity, _password, ok := r.BasicAuth()
	if !ok {
		s.Failure(w, errors.New("invalid authorization header"), http.StatusBadRequest)
		return
	}

	password, err := vo.NewPassword(_password)
	if err != nil {
		s.Failure(w, errInvalidCredentials, http.StatusUnauthorized)
		return
	}

	var res users.UserReadResult

	// If identity is email, get user using email, else get using username
	if strings.Contains(identity, "@") {
		email, errEmail := vo.NewEmail(identity)
		if errEmail != nil {
			_, _ = password.HashPassword()
			s.Failure(w, errInvalidCredentials, http.StatusUnauthorized)
			return
		}
		res = s.repo.UserGetByEmail(ctx, email)
	} else {
		username, errUsername := vo.NewUsername(identity)
		if errUsername != nil {
			_, _ = password.HashPassword()
			s.Failure(w, errInvalidCredentials, http.StatusUnauthorized)
			return
		}
		res = s.repo.UserGetByUsername(ctx, username)
	}

	if res.Error != nil {
		if res.StatusCode == http.StatusNotFound {
			s.Failure(w, errInvalidCredentials, http.StatusUnauthorized)
			return
		}
		s.Failure(w, res.Error, res.StatusCode)
		return
	}

	// validate password hash
	if err = password.VerifyPassword(res.User.Password.Value()); err != nil {
		s.Failure(w, errInvalidCredentials, http.StatusUnauthorized)
		return
	}

	if !res.User.Verified {
		s.Failure(w, errEmailNotVerified, http.StatusUnauthorized)
		return
	}

	// NOTE: Create Session And Tokens
	sessionID := utils.XIDNew().String()
	accessData := mw.NewAccessPayload(res.User.ID, sessionID)
	accessTokenDur := s.config.AuthAccessTokenExpire
	accessToken, accessTokenPayload, err := s.tokens.CreateToken(accessData, accessTokenDur)
	if err != nil {
		s.Failure(w, err)
		return
	}

	refreshData := mw.NewRefreshPayload(res.User.ID, sessionID)
	refreshTokenDur := s.config.AuthRefreshTokenExpire
	refreshToken, refreshTokenPayload, err := s.tokens.CreateToken(refreshData, refreshTokenDur)
	if err != nil {
		s.Failure(w, err)
		return
	}

	// INFO: saving access token to cache
	accessKey := mw.TokenKey(res.User.ID.Value().String(), sessionID, mw.TokenTypeAccess)
	if err = s.cache.Set(ctx, accessKey, accessTokenPayload, accessTokenDur); err != nil {
		s.Failure(w, err)
		return
	}

	// INFO: saving refresh token to cache
	refreshKey := mw.TokenKey(res.User.ID.Value().String(), sessionID, mw.TokenTypeRefresh)
	if err = s.cache.Set(ctx, refreshKey, refreshTokenPayload, refreshTokenDur); err != nil {
		s.Failure(w, err)
		return
	}

	// INFO: Seting Cookies
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
		&http.Cookie{
			Name: mw.CookieRefreshToken, Value: refreshToken,
			Path: "/", Expires: refreshTokenPayload.ExpiresAt,
			HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
		},
	)

	userParam := r.URL.Query().Get("user")
	getUser := userParam == "true" || (r.URL.Query().Has("user") && userParam == "")

	cookiesParams := r.URL.Query().Get("cookies")
	getCookies := cookiesParams == "true" ||
		(r.URL.Query().Has("cookies") && cookiesParams == "")

	response := SignInResponse{Message: "signed in successfully"}

	if getUser {
		response.User = res.User.Profile()
	}

	if getCookies {
		response.SessionID = sessionID
		response.AccessToken = accessToken
		response.RefreshToken = refreshToken
	}

	s.Success(w, response, res.StatusCode)
}
