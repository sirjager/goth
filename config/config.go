package config

import (
	"time"

	"github.com/go-redis/redis/v8"
)

// Conf holds the configuration settings for the application.
type Config struct {
	// GoEnv specifies the environment in which the application is running: dev/development, prod/production, or test/testing.
	GoEnv string `mapstructure:"GO_ENV" validate:"required,oneof=dev development prod production test testing"`

	// StartTime records the time when the configurations are loaded, useful for tracking service uptime.
	StartTime time.Time `validate:"required"`

	// Host denotes the host address where the API service is running.
	Host string `mapstructure:"HOST" validate:"required,hostname"`

	// Port specifies the port number on which the API service will listen.
	Port int `mapstructure:"PORT" validate:"required,number"`

	// ServiceName represents the alpha numberic name of the API service, used in logs and other areas.
	ServiceName string `mapstructure:"SERVICE_NAME" validate:"required,alphanum"`

	// ServerName is alphanum name that identifies the specific instance of the service, distinguishing it from other instances.
	ServerName string `mapstructure:"SERVER_NAME" validate:"required,alphanum"`

	// DocsSpecURL is the file path to the Swagger-generated JSON file for API documentation.
	DocsSpecURL string `validate:"required"`

	// PostgresURL is the connection string for the PostgreSQL database used by the service.
	PostgresURL string `mapstructure:"POSTGRES_URL" validate:"required"`

	// RedisURL is the connection string for the Redis instance used for async workers, sessions, and caching.
	RedisURL string `mapstructure:"REDIS_URL" validate:"required"`
	// RedisOptions holds additional configuration options for Redis.
	RedisOptions *redis.Options `                         validate:"required"`
	// RedisURLShort is a shorter version of the Redis connection string.
	RedisURLShort string `                         validate:"required"`

	// LoggerLogfile specifies the file path where logs will be written.
	LoggerLogfile string `mapstructure:"LOGGER_LOG_FILE"`

	// AuthTokenSecret is the secret key used for signing authentication tokens.
	AuthTokenSecret string `mapstructure:"AUTH_TOKEN_SECRET"            validate:"required,len=32"`
	// AuthSecureCookies indicates whether secure cookies should be used for authentication.
	AuthSecureCookies bool `mapstructure:"AUTH_SECURE_COOKIES"          validate:"boolean"`
	// AuthOAuthTokensExpire defines the expiration duration for OAuth tokens.
	AuthOAuthTokensExpire time.Duration `mapstructure:"AUTH_OAUTH_TOKEN_EXPIRE"      validate:"required,gte=30s"`
	// AuthAccessTokenExpire defines the expiration duration for access tokens.
	AuthAccessTokenExpire time.Duration `mapstructure:"AUTH_ACCESS_TOKEN_EXPIRE"     validate:"required,gte=30s"`
	// AuthRefreshTokenExpire defines the expiration duration for refresh tokens.
	AuthRefreshTokenExpire time.Duration `mapstructure:"AUTH_REFRESH_TOKEN_EXPIRE"    validate:"required,gte=30s"`
	// AuthEmailVerifyCooldown specifies the cooldown period between email verification requests.
	AuthEmailVerifyCooldown time.Duration `mapstructure:"AUTH_EMAIL_VERIFY_COOLDOWN"   validate:"required,gte=30s"`
	// AuthEmailVerifyExpire defines the expiration duration for email verification links.
	AuthEmailVerifyExpire time.Duration `mapstructure:"AUTH_EMAIL_VERIFY_EXPIRE"     validate:"required,gte=30s"`
	// AuthUserDeleteCooldown specifies the cooldown period between user deletion requests.
	AuthUserDeleteCooldown time.Duration `mapstructure:"AUTH_USER_DELETE_COOLDOWN"    validate:"required,gte=30s"`
	// AuthUserDeleteExpire defines the expiration duration for user deletion requests.
	AuthUserDeleteExpire time.Duration `mapstructure:"AUTH_USER_DELETE_EXPIRE"      validate:"required,gte=30s"`
	// AuthEmailChangeCooldown specifies the cooldown period between email change requests.
	AuthEmailChangeCooldown time.Duration `mapstructure:"AUTH_EMAIL_CHANGE_COOLDOWN"   validate:"required,gte=30s"`
	// AuthEmailChangeExpire defines the expiration duration for email change links.
	AuthEmailChangeExpire time.Duration `mapstructure:"AUTH_EMAIL_CHANGE_EXPIRE"     validate:"required,gte=30s"`
	// AuthPasswordResetCooldown specifies the cooldown period between password reset requests.
	AuthPasswordResetCooldown time.Duration `mapstructure:"AUTH_PASSWORD_RESET_COOLDOWN" validate:"required,gte=30s"`
	// AuthPasswordResetExpire defines the expiration duration for password reset links.
	AuthPasswordResetExpire time.Duration `mapstructure:"AUTH_PASSWORD_RESET_EXPIRE"   validate:"required,gte=30s"`
	// AuthGoogleClientID is the client ID for Google OAuth integration.
	AuthGoogleClientID string `mapstructure:"AUTH_GOOGLE_CLIENT_ID"        validate:"required,gte=50"`
	// AuthGoogleClientSecret is the client secret for Google OAuth integration.
	AuthGoogleClientSecret string `mapstructure:"AUTH_GOOGLE_CLIENT_SECRET"    validate:"required,gte=28"`
	// AuthGithubClientID is the client ID for GitHub OAuth integration.
	AuthGithubClientID string `mapstructure:"AUTH_GITHUB_CLIENT_ID"        validate:"required,gte=14"`
	// AuthGithubClientSecret is the client secret for GitHub OAuth integration.
	AuthGithubClientSecret string `mapstructure:"AUTH_GITHUB_CLIENT_SECRET"    validate:"required,gte=32"`

	// MailSMTPName is the email address that will appear as the sender in SMTP emails.
	MailSMTPName string `mapstructure:"MAIL_SMTP_Name" validate:"required,gte=4"`
	// MailSMTPUser is the username for authenticating with the SMTP server.
	MailSMTPUser string `mapstructure:"MAIL_SMTP_USER" validate:"required,email"`
	// MailSMTPPass is the password for authenticating with the SMTP server.
	MailSMTPPass string `mapstructure:"MAIL_SMTP_PASS" validate:"required,gte=14"`
}
