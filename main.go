package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"github.com/sirjager/gopkg/cache"
	"github.com/sirjager/gopkg/db"
	"github.com/sirjager/gopkg/mail"
	"github.com/sirjager/gopkg/tokens"
	"golang.org/x/sync/errgroup"

	"github.com/sirjager/goth/api"
	"github.com/sirjager/goth/config"
	"github.com/sirjager/goth/logger"
	"github.com/sirjager/goth/modules"
	"github.com/sirjager/goth/oauth"
	"github.com/sirjager/goth/repository"
	"github.com/sirjager/goth/worker"
)

var startTime time.Time

// NOTE: Listenting to thse signals for gracefull shutdown
var interuptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

//	@title			OAuthAPI
//	@version		0.1.0
//	@description	OAuth API for 3rd party authentication

//	@contact.name				Ankur Kumar
//	@contact.url				https://github.com/sirjager
//	@securityDefinitions.basic	BasicAuth

// @BasePath	/
func main() {
	// NOTE: change name of .env file here. For defaults, use "defaults"
	config := config.LoadConfigs(".", "defaults")
	logger, err := logger.NewLogger(config.ServerName, config.LoggerLogfile, config.GoEnv)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize logger")
	}
	defer logger.Close()
	logr := logger.Logr

	ctx, cancel := signal.NotifyContext(context.Background(), interuptSignals...)
	defer cancel()
	wg, ctx := errgroup.WithContext(ctx)

	// initializing sessions manager using redis backed
	oauth := oauth.NewOAuth(config, logr)
	if err = oauth.InitializeRedisStore(config.RedisURLShort, config.AuthTokenSecret); err != nil {
		fmt.Println("i am here")
		log.Fatal().Err(err).Msg("failed to initialize redis store")
	}
	defer oauth.Close(ctx, wg)

	database, conn, err := db.NewDatabae(ctx, db.Config{PostgresURL: config.PostgresURL}, logr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize database")
	}
	defer database.Close()

	// database repository
	repo, err := repository.NewRepository(conn, config.PostgresURL, logr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize repository")
	}

	redisClient := redis.NewClient(config.RedisOptions)
	if pingErr := redisClient.Ping(ctx).Err(); pingErr != nil {
		logger.Logr.Fatal().Err(pingErr).Msg("failed to ping redis client")
	}
	defer redisClient.Close()

	cache := cache.NewCacheRedis(redisClient, logger.Logr)

	mailConfig := mail.Config{
		SMTPSender: config.MailSMTPName,
		SMTPUser:   config.MailSMTPUser,
		SMTPPass:   config.MailSMTPPass,
	}
	mail, err := mail.NewGmailSender(mailConfig)
	if err != nil {
		logr.Fatal().Err(err).Msg("failed to initialize gmail smtp")
	}

	tokens, err := tokens.NewPasetoBuilder(config.AuthTokenSecret)
	if err != nil {
		logr.Fatal().Err(err).Msg("failed to create token builder")
	}

	tasks := worker.NewTaskDistributor(logr, config.RedisOptions)
	defer tasks.Shutdown()

	modules := modules.NewModules(config, logr, cache, repo, tokens, mail, tasks)
	worker.RunTaskProcessor(ctx, wg, logr, repo, mail, cache, tokens, config)

	server := api.NewServer(modules)

	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	server.StartServer(address, ctx, wg)

	err = wg.Wait()
	if err != nil {
		logr.Fatal().Err(err).Msg("error from wait group")
	}
}
