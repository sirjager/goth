package users

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/sirjager/goth/entity"
	repoerrors "github.com/sirjager/goth/repository/errors"
	usersrepomigrations "github.com/sirjager/goth/repository/users/migrations"
	"github.com/sirjager/goth/repository/users/sqlc"
	"github.com/sirjager/goth/vo"
)

type UserRepository interface {
	UserCreate(c context.Context, user *entity.User) UserReadResult
	UserDelete(c context.Context, userID *vo.ID) UserDeleteResult
	UserGetMaster(c context.Context) UserReadResult
	UserGetAll(c context.Context, limit, page int) UsersReadResult
	UserGetByID(c context.Context, userID *vo.ID) UserReadResult
	UserGetByEmail(c context.Context, email *vo.Email) UserReadResult
	UserGetByUsername(c context.Context, username *vo.Username) UserReadResult
	UserUpdate(c context.Context, user *entity.User) UserReadResult
	UserUpdateVerified(c context.Context, userID *vo.ID, status bool) UserReadResult
	UserUpdatePassword(c context.Context, userID *vo.ID, pass *vo.HashedPassword) UserReadResult
}

type userRepo struct {
	logr  zerolog.Logger
	store sqlc.Store
	pool  *pgxpool.Pool
}

func NewUsersRepo(conn *pgxpool.Pool, pgURL string, l zerolog.Logger) (UserRepository, error) {
	if err := usersrepomigrations.MigrateUsersRepo(pgURL); err != nil {
		l.Error().Err(err).Msg("failed to migrate users repo")
		return nil, repoerrors.ErrFailedToMigrate
	}

	store := sqlc.NewStore(conn)
	repo := &userRepo{logr: l, store: store, pool: conn}
	l.Info().Msg("users repository initialized")
	return repo, nil
}
