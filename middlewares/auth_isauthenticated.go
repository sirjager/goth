package mw

import (
	"net/http"

	"github.com/sirjager/gopkg/cache"
	"github.com/sirjager/gopkg/tokens"

	"github.com/sirjager/goth/entity"
	"github.com/sirjager/goth/payload"
	"github.com/sirjager/goth/repository"
	"github.com/sirjager/goth/vo"
)

func IsAuthenticated(
	r *http.Request,
	repo repository.Repo,
	tokens tokens.TokenBuilder,
	cache cache.Cache,
) (*entity.User, bool) {
	ctx := r.Context()
	var err error
	var user *entity.User
	user, err = authenticateUsingOAuth(r, repo)
	if err != nil || user == nil {
		// this will only be executed if has oauth session
		accessToken := extractAuthToken(r, CookieAccessToken)
		var incoming payload.AccessToken
		_, err = tokens.VerifyToken(accessToken, &incoming)
		if err != nil {
			return nil, false
		}
		tokenKey := payload.SessionAccessKey(incoming.UserID, incoming.SessionID)
		// session should be valid, check if sessions is not logged out
		var stored payload.AccessToken
		if err := cache.Get(ctx, tokenKey, &stored); err != nil {
			return nil, false
		}

		// match and verify if incoming and stored are valid
		userID, err := vo.NewIDFrom(incoming.UserID)
		if err != nil {
			return nil, false
		}
		if userID.Value().String() != stored.UserID {
			return nil, false
		}
		if incoming.SessionID != stored.SessionID {
			return nil, false
		}
		if incoming.TokenType != payload.TypeAccess {
			return nil, false
		}

		// now fetch user and add it to context
		res := repo.UserGetByID(ctx, userID)
		if res.Error != nil {
			return nil, false
		}
		user = res.User
	}
	if user != nil {
		return user, true
	}
	return nil, false
}
