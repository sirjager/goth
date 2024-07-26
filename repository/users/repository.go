package users

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	repoerrors "github.com/sirjager/goth/repository/errors"
	usersrepomigrations "github.com/sirjager/goth/repository/users/migrations"
	"github.com/sirjager/goth/repository/users/sqlc"
)

type repo struct {
	logr  zerolog.Logger
	store sqlc.Store
}

func NewUsersRepo(conn *pgxpool.Pool, pgURL string, l zerolog.Logger) (UsersRepo, error) {
	if err := usersrepomigrations.MigrateUsersRepo(pgURL); err != nil {
		l.Error().Err(err).Msg("failed to migrate users repo")
		return nil, repoerrors.ErrFailedToMigrate
	}
	store := sqlc.NewStore(conn)
	repo := &repo{logr: l, store: store}
	l.Info().Msg("repository initialized")
	return repo, nil
}
