package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/sirjager/gopkg/utils"

	"github.com/sirjager/goth/entity"
	mw "github.com/sirjager/goth/middlewares"
	"github.com/sirjager/goth/repository/users"
	"github.com/sirjager/goth/vo"
)

type SignInResponse struct {
	User         *entity.Profile `json:"user,omitempty"`
	AccessToken  string          `json:"access_token,omitempty"`
	RefreshToken string          `json:"refresh_token,omitempty"`
	SessionID    string          `json:"session_id,omitempty"`
	Message      string          `json:"message,omitempty"`
} //	@name	SignInResponse

const (
	TokenTypeAccess  = "0"
	TokenTypeRefresh = "1"
)

var errInvalidCredentials = errors.New("invalid credentials")

// Signin Request
//
//	@Summary		Signin
//	@Description	Signin using email and password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Router			/auth/signin [get]
//	@Security		BasicAuth
//	@Param			user	query		bool			false	"If true, returns User in body"
//	@Param			cookies	query		bool			false	"If true, returns AccessToken, RefreshToken and SessionID in body"
//	@Success		200		{object}	SignInResponse	"SignInResponse"
func (a *API) Signin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	identity, _password, ok := r.BasicAuth()
	if !ok {
		a.Failure(w, errors.New("invalid authorization header"), http.StatusBadRequest)
		return
	}

	password, err := vo.NewPassword(_password)
	if err != nil {
		a.Failure(w, errInvalidCredentials, http.StatusUnauthorized)
		return
	}

	var result users.UserReadResult

	if strings.Contains(identity, "@") {
		email, errEmail := vo.NewEmail(identity)
		if errEmail != nil {
			_, _ = password.HashPassword()
			a.Failure(w, errInvalidCredentials, http.StatusUnauthorized)
			return
		}
		result = a.repo.UserReadByEmail(ctx, email)
	} else {
		username, errUsername := vo.NewUsername(identity)
		if errUsername != nil {
			_, _ = password.HashPassword()
			a.Failure(w, errInvalidCredentials, http.StatusUnauthorized)
			return
		}
		result = a.repo.UserReadByUsername(ctx, username)
	}

	if result.Error != nil {
		_, _ = password.HashPassword()
		a.Failure(w, errInvalidCredentials, http.StatusUnauthorized)
		return
	}

	user := result.User

	// validate password hash
	if err = password.VerifyPassword(user.Password.Value()); err != nil {
		a.Failure(w, errInvalidCredentials, http.StatusUnauthorized)
		return
	}

	// if !user.Verified {
	// 	a.Failure(w, errors.New("email not verified"), http.StatusUnauthorized)
	// 	return
	// }

	// NOTE: Create Session And Tokens
	sessionID := utils.XIDNew().String()

	accessData := &mw.TokenPayloadData{
		UserID:    user.ID.Value().String(),
		SessionID: sessionID,
		TokenType: TokenTypeAccess,
	}
	// INFO: creating access token with OAuth.SessionsMaxAge duration
	accessTokenDur := a.config.OAuth.SessionsMaxAge
	accessToken, accessTokenPayload, err := a.tokens.CreateToken(accessData, accessTokenDur)
	if err != nil {
		a.Failure(w, err)
		return
	}

	refreshData := &mw.TokenPayloadData{
		UserID:    user.ID.Value().String(),
		SessionID: sessionID,
		TokenType: TokenTypeRefresh,
	}

	// INFO: creating refresh token 40% more than on accessTokenDur
	refreshTokenDur := time.Duration(float64(accessTokenDur) * 1.4)
	refreshToken, refreshTokenPayload, err := a.tokens.CreateToken(refreshData, refreshTokenDur)
	if err != nil {
		a.Failure(w, err)
		return
	}

	// INFO: saving access token to cache
	accessKey := mw.TokenKey(user.ID.Value().String(), sessionID, TokenTypeAccess)
	if err = a.cache.Set(ctx, accessKey, accessTokenPayload, accessTokenDur); err != nil {
		a.Failure(w, err)
		return
	}

	// INFO: saving refresh token to cache
	refreshKey := mw.TokenKey(user.ID.Value().String(), sessionID, TokenTypeRefresh)
	if err = a.cache.Set(ctx, refreshKey, refreshTokenPayload, refreshTokenDur); err != nil {
		a.Failure(w, err)
		return
	}

	// INFO: Seting Cookies
	a.SetCookies(w,
		// &http.Cookie{
		// 	Name: mw.SessionCookieName, Value: "",
		// 	Path: "/", Expires: time.Now(),
		// 	HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
		// },
		&http.Cookie{
			Name: "sessionId", Value: sessionID,
			Path: "/", Expires: accessTokenPayload.ExpiresAt,
			HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
		},
		&http.Cookie{
			Name: "accessToken", Value: accessToken,
			Path: "/", Expires: accessTokenPayload.ExpiresAt,
			HttpOnly: true, SameSite: http.SameSiteDefaultMode, Secure: false,
		},
		&http.Cookie{
			Name: "refreshToken", Value: refreshToken,
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
		response.User = result.User.Profile()
	}

	if getCookies {
		response.SessionID = sessionID
		response.AccessToken = accessToken
		response.RefreshToken = refreshToken
	}

	a.Success(w, response, result.StatusCode)
}
