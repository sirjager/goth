package mw

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/markbates/goth/gothic"
	"github.com/sirjager/gopkg/cache"
	"github.com/sirjager/gopkg/tokens"

	"github.com/sirjager/goth/entity"
	"github.com/sirjager/goth/repository"
	"github.com/sirjager/goth/vo"
)

func authenticateUsingTokens(
	r *http.Request,
	repo *repository.Repo,
	sessions cache.Cache,
	builder tokens.TokenBuilder,
) (*entity.User, error) {
	token := ""
	// get access token from cookies
	if cookie, cookieErr := r.Cookie("accessToken"); cookieErr == nil {
		token = cookie.Value
	}

	// get token from authorization header
	if token == "" {
		token, _ = extractAuthorizationToken(r)
	}

	// return err if no token
	if token == "" {
		return nil, ErrUnAuthorized
	}

	// check if incoming token is still valid or not
	payload, err := builder.VerifyToken(token)
	if err != nil {
		return nil, ErrUnAuthorized
	}

	data := payload.Payload.(map[string]interface{})

	userID := data["user_id"].(string)
	sessionID := data["session_id"].(string)
	tokenType := data["token_type"].(string)

	key := TokenKey(userID, sessionID, tokenType)

	var stored tokens.Payload
	if err = sessions.Get(r.Context(), key, &stored); err != nil {
		return nil, ErrUnAuthorized
	}

	// check if stored payload is still valid or not
	if err = stored.Valid(); err != nil {
		return nil, err
	}

	result := repo.UserReadByID(r.Context(), vo.MustParseID(userID))
	if result.Error != nil {
		if result.StatusCode == http.StatusNotFound {
			return nil, ErrUnAuthorized
		}
		return nil, result.Error
	}

	return result.User, err
}

func SessionKey(userID, sessionID string) string {
	return fmt.Sprintf("sess:%s:%s", userID, sessionID)
}

func UserSessionsKey(userID string) string {
	return fmt.Sprintf("sess:%s", userID)
}

func TokenKey(userID, sessionID string, tokenType string) string {
	return fmt.Sprintf(
		"sess:%s:%s:%s", // sess:userID:sessionID:(refresh|access)
		userID, sessionID, tokenType,
	)
}

func authenticateUsingOAuth(r *http.Request, repo *repository.Repo) (*entity.User, error) {
	session, err := gothic.Store.Get(r, SessionCookieName)
	if err != nil {
		return nil, err
	}

	value, ok := session.Values["email"].(string)
	if !ok || value == "" {
		return nil, ErrUnAuthorized
	}

	email, err := vo.NewEmail(value)
	if err != nil {
		return nil, err
	}
	result := repo.UserReadByEmail(r.Context(), email)
	if result.Error != nil {
		return nil, result.Error
	}
	return result.User, nil
}

func extractAuthorizationToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header not found")
	}
	values := strings.Split(authHeader, " ")
	if len(values) != 2 {
		return "", fmt.Errorf("invalid authorization header")
	}
	if values[0] != "Bearer" {
		return "", fmt.Errorf("unsupported authorization header")
	}
	return values[1], nil
}
