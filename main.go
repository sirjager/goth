package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sirjager/gopkg/db"
	"golang.org/x/sync/errgroup"

	"github.com/sirjager/goth/api"
	"github.com/sirjager/goth/config"
	"github.com/sirjager/goth/logger"
	"github.com/sirjager/goth/oauth"
	"github.com/sirjager/goth/repository"
)

var startTime time.Time

// NOTE: Listenting to thse signals for gracefull shutdown
var interuptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func init() {
	startTime = time.Now()
}

//	@title			OAuthAPI
//	@version		0.1.0
//	@description	OAuth API for 3rd party authentication

//	@contact.name	Ankur Kumar
//	@contact.url	https://github.com/sirjager

// @BasePath	/
func main() {
	// NOTE: change name of .env file here. For defaults, use "defaults"
	config, err := config.LoadConfigs(".", "defaults", startTime)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configurations")
	}

	logger, err := logger.NewLogger(config.Logger)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize logger")
	}
	defer logger.Close()
	logr := logger.Logr

	ctx, cancel := signal.NotifyContext(context.Background(), interuptSignals...)
	defer cancel()
	wg, ctx := errgroup.WithContext(ctx)

	// NOTE: http server address
	address := fmt.Sprintf("%s:%d", config.Host, config.RESTPort)

	redirect := fmt.Sprintf("http://%s", address)

	// initializing sessions manager using redis backed
	oauth := oauth.NewOAuth(redirect, config.OAuth, logr)
	if err = oauth.InitializeRedisStore(config.RedisURL, config.SecretKey); err != nil {
		log.Fatal().Err(err).Msg("failed to initialize redis store")
	}
	defer oauth.Close(ctx, wg)

	database, conn, err := db.NewDatabae(ctx, db.Config{PostgresURL: config.PostgresURL}, logr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize database")
	}
	defer database.Close()

	repo, err := repository.NewRepository(conn, config.PostgresURL, logr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize repository")
	}

	server := api.NewServer(repo, logr, config)

	server.StartServer(address, ctx, wg)

	err = wg.Wait()
	if err != nil {
		logr.Fatal().Err(err).Msg("error from wait group")
	}
}
