package worker

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"

	"github.com/sirjager/goth/payload"
)

type TaskDistributor interface {
	Shutdown()
	ResetPassword(
		ctx context.Context,
		payload *payload.ResetPassword,
		opts ...asynq.Option,
	) error
	SendEmailVerification(
		ctx context.Context,
		payload *payload.VerifyEmail,
		opts ...asynq.Option,
	) error
}

type dist struct {
	client *asynq.Client
	logr   zerolog.Logger
}

func NewTaskDistributor(logr zerolog.Logger, redisOptions *redis.Options) TaskDistributor {
	client := asynq.NewClient(asynq.RedisClientOpt{
		DB:        redisOptions.DB,
		Addr:      redisOptions.Addr,
		Network:   redisOptions.Network,
		Password:  redisOptions.Password,
		Username:  redisOptions.Username,
		TLSConfig: redisOptions.TLSConfig,
		PoolSize:  redisOptions.PoolSize,
	})
	return &dist{client, logr}
}

// close redis
func (d *dist) Shutdown() {
	d.client.Close()
}
