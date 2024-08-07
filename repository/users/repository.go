package users

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	repoerrors "github.com/sirjager/goth/repository/errors"
	usersrepomigrations "github.com/sirjager/goth/repository/users/migrations"
	"github.com/sirjager/goth/repository/users/sqlc"
)

type UserRepo struct {
	logr  zerolog.Logger
	store sqlc.Store
	pool  *pgxpool.Pool
}

func NewUsersRepo(conn *pgxpool.Pool, pgURL string, l zerolog.Logger) (*UserRepo, error) {
	if err := usersrepomigrations.MigrateUsersRepo(pgURL); err != nil {
		l.Error().Err(err).Msg("failed to migrate users repo")
		return nil, repoerrors.ErrFailedToMigrate
	}
	store := sqlc.NewStore(conn)
	repo := &UserRepo{logr: l, store: store, pool: conn}
	l.Info().Msg("users repository initialized")
	return repo, nil
}
