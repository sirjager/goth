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
	UserCreate(ctx context.Context, user *entity.User) UserReadResult
	UserDelete(ctx context.Context, userID *vo.ID) UserDeleteResult
	UserDeleteTx(ctx context.Context, params UserDeleteTxParams) UserDeleteResult
	UserGetMaster(ctx context.Context) UserReadResult
	UserGetAll(ctx context.Context, limit, page int) UsersReadResult
	UserGetByID(ctx context.Context, userID *vo.ID) UserReadResult
	UserGetByEmail(ctx context.Context, email *vo.Email) UserReadResult
	UserGetByUsername(ctx context.Context, username *vo.Username) UserReadResult
	UserUpdate(ctx context.Context, user *entity.User) UserReadResult
	UserUpdateVerified(ctx context.Context, userID *vo.ID, status bool) UserReadResult
	UserUpdatePassword(ctx context.Context, userID *vo.ID, pass *vo.HashedPassword) UserReadResult
	UserUpdatePasswordTx(ctx context.Context, params UserUpdatePasswordTxParams) UserReadResult
}

type repo struct {
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
	repo := &repo{logr: l, store: store, pool: conn}
	l.Info().Msg("users repository initialized")
	return repo, nil
}
