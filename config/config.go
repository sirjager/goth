package config

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/sirjager/goth/logger"
	"github.com/sirjager/goth/oauth"
)

// Config represents the application configuration.
type Config struct {
	StartTime     time.Time
	Logger        logger.Config
	ServiceName   string `mapstructure:"SERVICE_NAME"`
	ServerName    string `mapstructure:"SERVER_NAME"`
	Host          string `mapstructure:"HOST"`
	RedisURL      string `mapstructure:"REDIS_URL"`
	RedisURLShort string
	PostgresURL   string `mapstructure:"POSTGRES_URL"`
	SecretKey     string `mapstructure:"SECRET_KEY"`
	OAuth         oauth.Config
	RESTPort      int `mapstructure:"REST_PORT"`
}

// LoadConfigs loads the configuration from the specified YAML file.
func LoadConfigs(path string, name string, startTime time.Time) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	if err = viper.Unmarshal(&config); err != nil {
		return
	}

	if err = viper.Unmarshal(&config.Logger); err != nil {
		return
	}

	if err = viper.Unmarshal(&config.OAuth); err != nil {
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		return Config{}, err
	}

	if config.Host == "" {
		config.Host = "localhost"
	}

	// Construct the DBUrl using the DBConfig values.
	if config.ServerName == "" {
		config.ServerName = hostname
		config.Logger.ServerName = config.ServerName
	}

	config.StartTime = startTime
	config.RedisURLShort = strings.Split(config.RedisURL, "@")[1]

	return
}
