package core

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

// App holds the core components for the server.
type App struct {
	logger          zerolog.Logger         // Logger for logging information.
	cache           cache.Cache            // Cache for storing and retrieving data.
	tokenBuilder    tokens.TokenBuilder    // TokenBuilder for generating and managing tokens.
	repository      repository.Repo        // Repository for data access.
	validator       *validator.Validate    // Validator for validating structs.
	configs         *config.Config         // Configuration settings for the server.
	mailSender      mail.Sender            // MailSender for sending emails.
	taskDistrubutor worker.TaskDistributor // TaskDistributor for distributing tasks.
}

func NewCoreApp(
	configs *config.Config, // Configuration settings.
	logger zerolog.Logger, // Logger instance.
	cache cache.Cache, // Cache instance.
	repository repository.Repo, // Repository instance.
	tokenBuilder tokens.TokenBuilder, // TokenBuilder instance.
	mailSender mail.Sender, // MailSender instance.
	taskDistributor worker.TaskDistributor, // TaskDistributor instance.

) *App {
	_validator := validator.New(validator.WithRequiredStructEnabled()) // Create a new validator.
	return &App{
		logger:          logger,
		cache:           cache,
		tokenBuilder:    tokenBuilder,
		repository:      repository,
		validator:       _validator,
		configs:         configs,
		mailSender:      mailSender,
		taskDistrubutor: taskDistributor,
	}

}

// Logger returns the logger instance.
func (m *App) Logger() *zerolog.Logger {
	return &m.logger
}

// Cache returns the cache instance.
func (m *App) Cache() cache.Cache {
	return m.cache
}

// Mailer returns the mail sender instance.
func (m *App) Mailer() mail.Sender {
	return m.mailSender
}

// Tokens returns the token builder instance.
func (m *App) Tokens() tokens.TokenBuilder {
	return m.tokenBuilder
}

// Repo returns the repository instance.
func (m *App) Repo() repository.Repo {
	return m.repository
}

// Validator returns the validator instance.
func (m *App) Validator() *validator.Validate {
	return m.validator
}

// Config returns the configuration settings.
func (m *App) Config() *config.Config {
	return m.configs
}

// Tasks returns the task distributor instance.
func (m *App) Tasks() worker.TaskDistributor {
	return m.taskDistrubutor
}
