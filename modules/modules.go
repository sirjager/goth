package modules

import (
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/sirjager/gopkg/cache"
	"github.com/sirjager/gopkg/mail"
	"github.com/sirjager/gopkg/tokens"

	"github.com/sirjager/goth/config"
	"github.com/sirjager/goth/repository"
	"github.com/sirjager/goth/worker"
)

// Modules holds the core components for the server.
type Modules struct {
	logger          zerolog.Logger         // Logger for logging information.
	cache           cache.Cache            // Cache for storing and retrieving data.
	tokenBuilder    tokens.TokenBuilder    // TokenBuilder for generating and managing tokens.
	repository      repository.Repo        // Repository for data access.
	validator       *validator.Validate    // Validator for validating structs.
	configs         *config.Config         // Configuration settings for the server.
	mailSender      mail.Sender            // MailSender for sending emails.
	taskDistrubutor worker.TaskDistributor // TaskDistributor for distributing tasks.
}

// NewModules creates a new instance of Modules with the provided components.
func NewModules(
	configs *config.Config, // Configuration settings.
	logger zerolog.Logger, // Logger instance.
	cache cache.Cache, // Cache instance.
	repository repository.Repo, // Repository instance.
	tokenBuilder tokens.TokenBuilder, // TokenBuilder instance.
	mailSender mail.Sender, // MailSender instance.
	taskDistributor worker.TaskDistributor, // TaskDistributor instance.
) *Modules {
	validator := validator.New(validator.WithRequiredStructEnabled()) // Create a new validator.
	return &Modules{
		logger:          logger,
		cache:           cache,
		tokenBuilder:    tokenBuilder,
		repository:      repository,
		validator:       validator,
		configs:         configs,
		mailSender:      mailSender,
		taskDistrubutor: taskDistributor,
	}
}

// Logger returns the logger instance.
func (m *Modules) Logger() *zerolog.Logger {
	return &m.logger
}

// Cache returns the cache instance.
func (m *Modules) Cache() cache.Cache {
	return m.cache
}

// Mailer returns the mail sender instance.
func (m *Modules) Mailer() mail.Sender {
	return m.mailSender
}

// Tokens returns the token builder instance.
func (m *Modules) Tokens() tokens.TokenBuilder {
	return m.tokenBuilder
}

// Repo returns the repository instance.
func (m *Modules) Repo() repository.Repo {
	return m.repository
}

// Validator returns the validator instance.
func (m *Modules) Validator() *validator.Validate {
	return m.validator
}

// Config returns the configuration settings.
func (m *Modules) Config() *config.Config {
	return m.configs
}

// Tasks returns the task distributor instance.
func (m *Modules) Tasks() worker.TaskDistributor {
	return m.taskDistrubutor
}
