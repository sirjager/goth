package api

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sirjager/gopkg/cache"
	"github.com/sirjager/gopkg/tokens"

	"github.com/sirjager/goth/config"
	"github.com/sirjager/goth/oauth"
	"github.com/sirjager/goth/worker"
)

var (
	testConfig *config.Config
	testCache  cache.Cache
	testTokens tokens.TokenBuilder
	testLogr   zerolog.Logger
	testTasks  worker.TaskDistributor
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
	tasks := worker.NewTaskDistributor(logr, redisOptions)

	tokens, err := tokens.NewPasetoBuilder(config.AuthTokenSecret)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create token builder")
	}

	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	redirect := fmt.Sprintf("http://%s", address)
	// initializing sessions manager using redis backed
	oauth := oauth.NewOAuth(redirect, config, logr)
	if err = oauth.InitializeRedisStore(config.RedisURLShort, config.AuthTokenSecret); err != nil {
		log.Fatal().Err(err).Msg("failed to initialize redis store")
	}

	testConfig = config
	testTokens = tokens
	testCache = cache
	testTasks = tasks
	os.Exit(M.Run())
}
