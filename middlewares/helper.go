package mw

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/markbates/goth/gothic"

	"github.com/sirjager/goth/entity"
	"github.com/sirjager/goth/repository"
	"github.com/sirjager/goth/vo"
)

func extractAuthToken(r *http.Request, cookieName string) (token string) {
	if cookie, cookieErr := r.Cookie(cookieName); cookieErr == nil {
		token = cookie.Value
	}
	// If access token is already in cookies return it
	if token != "" {
		return
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return
	}
	values := strings.Split(authHeader, " ")
	if len(values) != 2 {
		return
	}
	if values[0] != "Bearer" {
		return
	}
	token = values[1]
	return
}

func authenticateUsingOAuth(r *http.Request, repo repository.Repo) (*entity.User, error) {
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

func extractAccessToken(r *http.Request) string {
	token, _ := extractAuthorizationToken(r)
	if cookie, cookieErr := r.Cookie(CookieAccessToken); cookieErr == nil {
		if cookie.Value != "" {
			token = cookie.Value
		}
	}
	return token
}
