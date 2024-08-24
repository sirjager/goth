package mw

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/markbates/goth/gothic"
	"github.com/rs/zerolog"
	"github.com/sirjager/gopkg/cache"
	"github.com/sirjager/gopkg/mail"
	"github.com/sirjager/gopkg/tokens"

	"github.com/sirjager/goth/config"
	"github.com/sirjager/goth/entity"
	"github.com/sirjager/goth/repository"
	"github.com/sirjager/goth/vo"
	"github.com/sirjager/goth/worker"
)

type Adapters struct {
	Repo     repository.Repo
	Toknb    tokens.TokenBuilder
	Cache    cache.Cache
	Logr     zerolog.Logger
	Validate *validator.Validate
	Config   *config.Config
	Mail     mail.Sender
	Tasks    worker.TaskDistributor
}

func LoadAdapters(
	config *config.Config,
	repository repository.Repo,
	tokenBuilder tokens.TokenBuilder,
	cache cache.Cache,
	logger zerolog.Logger,
	mailSender mail.Sender,
	taskDistributor worker.TaskDistributor,
) *Adapters {
	return &Adapters{
		Repo:     repository,
		Toknb:    tokenBuilder,
		Cache:    cache,
		Logr:     logger,
		Config:   config,
		Mail:     mailSender,
		Tasks:    taskDistributor,
		Validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

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
