package mw

import (
	"errors"
	"net/http"

	"github.com/sirjager/gopkg/cache"
	"github.com/sirjager/gopkg/tokens"

	"github.com/sirjager/goth/entity"
	"github.com/sirjager/goth/modules"
	"github.com/sirjager/goth/payload"
	"github.com/sirjager/goth/vo"
)

func IsAuthenticated(r *http.Request, modules *modules.Modules) (*entity.User, bool) {
	user, _, err := getAuthenticatedUser(r, modules, CookieAccessToken)
	if err != nil || user == nil {
		return nil, false
	}
	if err = user.Email.Validate(); err != nil {
		return nil, false
	}
	return user, true
}

func getAuthenticatedUser(
	r *http.Request,
	modules *modules.Modules,
	cookieName string,
) (*entity.User, int, error) {
	ctx := r.Context()
	var err error
	var user *entity.User

	// not allowing to use refresh token
	user, err = authenticateUsingOAuth(r, modules.Repo())
	if err == nil && user != nil && user.Email.Value() != "" {
		// if we have user and no error we return user
		return user, http.StatusOK, err
	}

	accessToken := extractAuthToken(r, cookieName)
	var accessData payload.BaseAuthPayload

	// if invalid access tokens
	if _, err = modules.Tokens().VerifyToken(accessToken, &accessData); err != nil {
		modules.Logger().Error().Msg("access token verification failed")
		return nil, http.StatusUnauthorized, errors.New(unauthorized)
	}

	// retriving stored session, to check if sessions is valid and not expired
	var accessPayload tokens.Payload
	accessKey := payload.SessionAccessKey(accessData.UserID, accessData.SessionID)
	if cookieName == CookieRefreshToken {
		accessKey = payload.SessionRefreshKey(accessData.UserID, accessData.SessionID)
	}

	if err = modules.Cache().Get(ctx, accessKey, &accessPayload); err != nil {
		// any other internal error, than not found
		if !errors.Is(err, cache.ErrNoRecord) {
			return nil, http.StatusInternalServerError, err
		}
		// if session is not found, means its expired
		modules.Logger().Error().Msg("session not found")
		return nil, http.StatusUnauthorized, errors.New(unauthorized)
	}

	// extract payload from stored session to match with incoming access token's payload
	var storedPayload payload.BaseAuthPayload
	if err = modules.Tokens().ReadPayload(&accessPayload, &storedPayload); err != nil {
		modules.Logger().Error().Msg("failed to read stored payload")
		return nil, http.StatusUnauthorized, errors.New(unauthorized)
	}

	// parse user id from incoming access token payload
	userID, err := vo.NewIDFrom(accessData.UserID)
	if err != nil {
		modules.Logger().Error().Msg("invalid user id")
		return nil, http.StatusUnauthorized, errors.New(unauthorized)
	}

	// match each field of incoming access token payload with stored access token payload
	if storedPayload.CreatedAt != accessData.CreatedAt {
		return nil, http.StatusUnauthorized, errors.New(unauthorized)
	}
	if storedPayload.TokenType != accessData.TokenType {
		return nil, http.StatusUnauthorized, errors.New(unauthorized)
	}
	if storedPayload.SessionID != accessData.SessionID {
		return nil, http.StatusUnauthorized, errors.New(unauthorized)
	}
	if storedPayload.UserID != accessData.UserID {
		return nil, http.StatusUnauthorized, errors.New(unauthorized)
	}

	// get user from repository, to ensure that user exists, and is valid
	res := modules.Repo().UserGetByID(ctx, userID)
	if res.Error != nil {
		// if its internal error, return as it is
		if res.StatusCode != http.StatusNotFound {
			return nil, res.StatusCode, res.Error
		}
		// if user not found, return unauthorized
		modules.Logger().Error().Msg("user not found")
		return nil, http.StatusUnauthorized, errors.New(unauthorized)
	}
	err = res.Error
	user = res.User

	// INFO: ONLY SUCCESSFUL RETURN
	if user != nil && user.ID.IsEqual(userID.Value()) && err == nil {
		return user, http.StatusOK, err
	}

	return nil, http.StatusUnauthorized, errors.New(unauthorized)
}
