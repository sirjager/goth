package worker

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

type TaskDistributor interface {
	Shutdown()
	SendEmailVerification(
		ctx context.Context,
		payload SendEmailVerificationPayload,
		opts ...asynq.Option,
	) error
}

type distributor struct {
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
	return &distributor{client, logr}
}

// close redis
func (d *distributor) Shutdown() {
	d.client.Close()
}
