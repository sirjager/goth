package oauth

import (
	"context"
	"fmt"
	"time"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"gopkg.in/boj/redistore.v1"
)

type Config struct {
	GoogleClientID     string        `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string        `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GithubClientID     string        `mapstructure:"GITHUB_CLIENT_ID"`
	GithubClientSecret string        `mapstructure:"GITHUB_CLIENT_SECRET"`
	SessionsMaxAge     time.Duration `mapstructure:"SESSIONS_MAX_AGE"`
	SecureCookies      bool          `mapstructure:"SECURE_COOKIES"`
}

type OAuth struct {
	logr     zerolog.Logger
	store    *redistore.RediStore
	redirect string
	config   Config
}

func NewOAuth(redirect string, config Config, logr zerolog.Logger) *OAuth {
	return &OAuth{
		config:   config,
		logr:     logr,
		redirect: redirect,
	}
}

func (o *OAuth) InitializeRedisStore(address, secretKey string) (err error) {
	store, err := redistore.NewRediStore(20, "tcp", address, "", []byte(secretKey))
	if err != nil {
		return err
	}
	store.Options.HttpOnly = true
	store.Options.Secure = o.config.SecureCookies
	store.SetMaxAge(int(o.config.SessionsMaxAge.Seconds()))

	o.store = store

	c := o.config
	gothic.Store = o.store

	goth.UseProviders(
		google.New(c.GoogleClientID, c.GoogleClientSecret, callbackURL(o, "google")),
		github.New(c.GithubClientID, c.GithubClientSecret, callbackURL(o, "github"), "user:email"),
	)
	return
}

func (o *OAuth) Close(ctx context.Context, wg *errgroup.Group) {
	wg.Go(func() error {
		<-ctx.Done()
		o.logr.Info().Msg("gracefully shutting redis store")

		if err := o.store.Close(); err != nil {
			o.logr.Error().Err(err).Msg("failed to shutdown redis store")
			return err
		}
		o.logr.Info().Msg("redis store gracefully shutdown complete")
		return nil
	})
}

func callbackURL(o *OAuth, provider string) string {
	return fmt.Sprintf("%s/auth/%s/callback", o.redirect, provider)
}
