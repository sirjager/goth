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
	repo repository.Repository,
	sessions cache.Cache,
	builder tokens.TokenBuilder,
	allowRefreshToken ...bool,
) (*entity.User, error) {
	// tokens can be at multiple places in the request
	// In cookies, in authorization header, or as query parameters
	// due to security reasons we will not be accepting tokens as query params
	// as query params are directly expored in url.

	token := "" // initialize token as empty

	// extract token from cookie
	if cookie, cookieErr := r.Cookie(CookieAccessToken); cookieErr == nil {
		token = cookie.Value
	}

	// if token is still empty, extract it from authorization header
	if token == "" {
		token, _ = extractAuthorizationToken(r)
	}

	if token == "" && len(allowRefreshToken) == 1 && allowRefreshToken[0] {
		if cookie, cookieErr := r.Cookie(CookieRefreshToken); cookieErr == nil {
			token = cookie.Value
		}
	}

	// return error if token is still empty
	if token == "" {
		return nil, ErrUnAuthorized
	}

	// incoming token payload
	var incoming TokenCustomPayload
	// check if incoming token is still valid or not
	_, err := builder.VerifyToken(token, &incoming)
	if err != nil {
		return nil, ErrUnAuthorized
	}

	// stored token payload, to make sure token payload is ours  and still active.
	var storedPayload tokens.Payload
	storedPayloadKey := TokenKey(incoming.UserID, incoming.SessionID, incoming.TokenType)
	if err = sessions.Get(r.Context(), storedPayloadKey, &storedPayload); err != nil {
		return nil, ErrUnAuthorized
	}

	// stored payload custom data
	var storedPayloadData TokenCustomPayload
	if err = builder.ReadPayload(&storedPayload, &storedPayloadData); err != nil {
		return nil, ErrUnAuthorized
	}

	// though it is not necessarily to match or validate anything else
	// if tokens is validated, and we have stored payload everything is fine
	// but if we want we can validate ourselves both custom payload data
	if incoming.UserID != storedPayloadData.UserID {
		return nil, ErrUnAuthorized
	}
	if incoming.SessionID != storedPayloadData.SessionID {
		return nil, ErrUnAuthorized
	}

	result := repo.UserGetByID(r.Context(), vo.MustParseID(incoming.UserID))
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

func authenticateUsingOAuth(r *http.Request, repo repository.Repository) (*entity.User, error) {
	session, err := gothic.Store.Get(r, CookieGothicSession)
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
	result := repo.UserGetByEmail(r.Context(), email)
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
