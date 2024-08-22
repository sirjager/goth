package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/sirjager/gopkg/utils"
	"golang.org/x/net/context"

	"github.com/sirjager/goth/entity"
	mw "github.com/sirjager/goth/middlewares"
	"github.com/sirjager/goth/payload"
	"github.com/sirjager/goth/repository"
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

	// get user from repository result

	res := _getUser(r.Context(), identity, s.repo)
	if res.Error != nil {
		if res.StatusCode == http.StatusBadRequest {
			s.Failure(w, errInvalidCredentials, http.StatusUnauthorized)
			return
		}
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

	sessionID := utils.XIDNew().String()

	// create and save access token payload
	accessData := payload.NewAccessPayload(res.User, sessionID)
	accessExpiry := s.config.AuthAccessTokenExpire
	accessToken, accessPayload, err := s.toknb.CreateToken(accessData, accessExpiry)
	if err != nil {
		s.Failure(w, err)
		return
	}
	accessKey := payload.SessionAccessKey(res.User.ID.Value().String(), sessionID)
	if err = s.cache.Set(ctx, accessKey, accessPayload, accessExpiry); err != nil {
		s.Failure(w, err)
		return
	}

	// create and save refresh token payload
	refreshData := payload.NewRefreshPayload(res.User, sessionID)
	refreshExpiry := s.config.AuthRefreshTokenExpire
	refreshToken, refreshPayload, err := s.toknb.CreateToken(refreshData, refreshExpiry)
	if err != nil {
		s.Failure(w, err)
		return
	}
	refreshKey := payload.SessionRefreshKey(res.User.ID.Value().String(), sessionID)
	if err = s.cache.Set(ctx, refreshKey, refreshPayload, refreshExpiry); err != nil {
		s.Failure(w, err)
		return
	}

	// INFO: Seting Cookies
	s.SetCookies(w,
		&http.Cookie{
			Name: mw.CookieSessionID, Value: sessionID,
			Path: "/", Expires: accessPayload.ExpiresAt,
			HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
		},
		&http.Cookie{
			Name: mw.CookieAccessToken, Value: accessToken,
			Path: "/", Expires: accessPayload.ExpiresAt,
			HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
		},
		&http.Cookie{
			Name: mw.CookieRefreshToken, Value: refreshToken,
			Path: "/", Expires: refreshPayload.ExpiresAt,
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

func _getUser(ctx context.Context, id string, repo repository.Repo) users.UserReadResult {
	if strings.Contains(id, "@") {
		email, errEmail := vo.NewEmail(id)
		if errEmail != nil {
			return users.UserReadResult{StatusCode: http.StatusBadRequest, Error: errEmail}
		}
		return repo.UserGetByEmail(ctx, email)
	}
	username, errUsername := vo.NewUsername(id)
	if errUsername != nil {
		return users.UserReadResult{StatusCode: http.StatusBadRequest, Error: errUsername}
	}
	return repo.UserGetByUsername(ctx, username)
}
