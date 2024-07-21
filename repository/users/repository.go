package users

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/sirjager/goth/entity"
	repoerrors "github.com/sirjager/goth/repository/errors"
	usersrepomigrations "github.com/sirjager/goth/repository/users/migrations"
	"github.com/sirjager/goth/repository/users/sqlc"
)

type UsersRepo interface {
	UserCreate(c context.Context, u *entity.User) (*entity.User, error)
	UserReadByID(c context.Context, uid string) (*entity.User, error)
	UserReadByEmail(c context.Context, email string) (*entity.User, error)
	UserDelete(c context.Context, uid string) error
}

type repo struct {
	logr  zerolog.Logger
	store sqlc.Store
}

func NewUsersRepo(conn *pgxpool.Pool, pgURL string, l zerolog.Logger) (UsersRepo, error) {
	if err := usersrepomigrations.MigrateUsersRepo(pgURL); err != nil {
		return nil, repoerrors.ErrFailedToMigrate
	}
	store := sqlc.NewStore(conn)
	repo := &repo{logr: l, store: store}
	l.Info().Msg("repository initialized")
	return repo, nil
}
