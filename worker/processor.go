package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/sirjager/gopkg/cache"
	"github.com/sirjager/gopkg/mail"
	"github.com/sirjager/gopkg/tokens"

	"github.com/sirjager/goth/config"
	"github.com/sirjager/goth/repository"
)

type TaskProcessor interface {
	Start() error
	Shutdown()
	SendEmailVerification(ctx context.Context, tasks *asynq.Task) error
}

const (
	PriorityCritical = "critical"
	PriorityUrgent   = "urgent"
	PriorityDefault  = "default"
	PriorityLow      = "low"
	PriorityLazy     = "lazy"
)

type proc struct {
	logr   zerolog.Logger
	repo   repository.Repo
	mail   mail.Sender
	cache  cache.Cache
	tokens tokens.TokenBuilder
	server *asynq.Server
	config *config.Config
}

func newTaskProcessor(
	logr zerolog.Logger,
	repo repository.Repo,
	mail mail.Sender,
	cache cache.Cache,
	tokens tokens.TokenBuilder,
	config *config.Config,
	opts asynq.RedisClientOpt,
) (TaskProcessor, error) {
	clientConfig := asynq.Config{
		Queues: map[string]int{
			PriorityCritical: 10,
			PriorityUrgent:   7,
			PriorityDefault:  5,
			PriorityLow:      3,
			PriorityLazy:     1,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(
			func(ctx context.Context, task *asynq.Task, err error) {
				logr.Error().Err(err).
					Str("type", task.Type()).
					Interface("payload", task.Payload()).
					Msg("failed to process task")
			},
		),
		Logger: NewLogger(logr),
	}

	server := asynq.NewServer(opts, clientConfig)

	return &proc{
		server: server,
		logr:   logr,
		repo:   repo,
		mail:   mail,
		cache:  cache,
		tokens: tokens,
		config: config,
	}, nil
}

// Start starts the RedisTaskProcessor
func (p *proc) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendEmailVerification, p.SendEmailVerification)

	return p.server.Start(mux)
}

func (p *proc) Shutdown() {
	p.server.Shutdown()
}
