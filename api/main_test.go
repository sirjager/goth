package api

import (
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sirjager/gopkg/cache"
	"github.com/sirjager/gopkg/mail"
	"github.com/sirjager/gopkg/tokens"

	"github.com/sirjager/goth/config"
	"github.com/sirjager/goth/oauth"
)

var (
	testConfig *config.Config
	testCache  cache.Cache
	testTokens tokens.TokenBuilder
	testLogr   zerolog.Logger
	testMail   mail.Sender
)

func TestMain(M *testing.M) {
	logr := zerolog.New(os.Stdout)
	config := config.LoadConfigs("..", "defaults", "test")
	redisOptions, redisOptionsParseErr := redis.ParseURL(config.RedisURL)
	if redisOptionsParseErr != nil {
		logr.Fatal().Err(redisOptionsParseErr).Msg("failed to parse redis url")
	}

	redisClient := redis.NewClient(redisOptions)
	cache := cache.NewCacheRedis(redisClient, logr)

	tokens, err := tokens.NewPasetoBuilder(config.AuthTokenSecret)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create token builder")
	}

	// initializing sessions manager using redis backed
	oauth := oauth.NewOAuth(config, logr)
	if err = oauth.InitializeRedisStore(config.RedisURLShort, config.AuthTokenSecret); err != nil {
		log.Fatal().Err(err).Msg("failed to initialize redis store")
	}

	testConfig = config
	testTokens = tokens
	testCache = cache
	os.Exit(M.Run())
}
