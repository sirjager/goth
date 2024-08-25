package oauth

import (
	"context"
	"fmt"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"gopkg.in/boj/redistore.v1"

	"github.com/sirjager/goth/config"
)

type OAuth struct {
	logr     zerolog.Logger
	store    *redistore.RediStore
	config   *config.Config
	redirect string
}

func NewOAuth(config *config.Config, logr zerolog.Logger) *OAuth {
	redirect := fmt.Sprintf("http://%s:%d", config.Host, config.Port)
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
	store.Options.Secure = o.config.AuthSecureCookies
	store.Options.MaxAge = int(o.config.AuthOAuthTokensExpire.Seconds())
	store.SetMaxAge(int(o.config.AuthOAuthTokensExpire.Seconds()))

	o.store = store

	c := o.config
	gothic.Store = o.store

	goth.UseProviders(
		google.New(c.AuthGoogleClientID, c.AuthGoogleClientSecret, callback(o, "google")),
		github.New(
			c.AuthGithubClientID,
			c.AuthGithubClientSecret,
			callback(o, "github"),
			"user:email",
		),
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

func callback(o *OAuth, provider string) string {
	return fmt.Sprintf("%s/auth/%s/callback", o.redirect, provider)
}
