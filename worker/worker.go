package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/sirjager/gopkg/cache"
	"github.com/sirjager/gopkg/mail"
	"github.com/sirjager/gopkg/tokens"
	"golang.org/x/sync/errgroup"

	"github.com/sirjager/goth/config"
	"github.com/sirjager/goth/repository"
)

func RunTaskProcessor(
	ctx context.Context,
	wg *errgroup.Group,
	logr zerolog.Logger,
	repo repository.Repo,
	mail mail.Sender,
	cache cache.Cache,
	tokens tokens.TokenBuilder,
	config *config.Config,
) {
	redisOptions := config.RedisOptions
	opts := asynq.RedisClientOpt{
		DB:        redisOptions.DB,
		Addr:      redisOptions.Addr,
		Network:   redisOptions.Network,
		Password:  redisOptions.Password,
		Username:  redisOptions.Username,
		TLSConfig: redisOptions.TLSConfig,
		PoolSize:  redisOptions.PoolSize,
	}
	processor, err := newTaskProcessor(logr, repo, mail, cache, tokens, config, opts)
	if err != nil {
		logr.Fatal().Err(err).Msg("failed to create task processor")
	}
	logr.Info().Msgf("started task processor")
	if err := processor.Start(); err != nil {
		logr.Fatal().Err(err).Msg("failed to start task processor")
	}

	wg.Go(func() error {
		<-ctx.Done()
		logr.Info().Msg("gracefully shutting down task processor...")
		processor.Shutdown()
		logr.Info().Msg("task process has been stopped")
		return nil
	})
}
