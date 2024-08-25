package config

import (
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"github.com/sirjager/gopkg/utils"
	"github.com/spf13/viper"
)

// LoadConfigs loads the configuration from the specified YAML file.
func LoadConfigs(path string, name string, env ...string) *Config {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	var config Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("failed to read configs")
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal().Err(err).Msg("failed to unmarshal configs")
	}
	config.StartTime = time.Now()

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get hostname")
	}

	// Construct the DBUrl using the DBConfig values.
	if config.ServerName == "" {
		config.ServerName = hostname
	}

	redisOptions, err := redis.ParseURL(config.RedisURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse redis url")
	}
	config.RedisOptions = redisOptions
	config.RedisURLShort = strings.Split(config.RedisURL, "@")[1]

	config.GoEnv = utils.GetFirstOrFallback("dev", env...)

	validator := validator.New(validator.WithRequiredStructEnabled())

	if err = validator.Struct(config); err != nil {
		log.Fatal().Err(err).Msg("invalid configs")
	}

	return &config
}
