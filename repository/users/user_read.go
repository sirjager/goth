package users

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/sirjager/goth/entity"
	repoerrors "github.com/sirjager/goth/repository/errors"
	"github.com/sirjager/goth/repository/users/sqlc"
)

type UsersReadResult struct {
	Error      error
	Users      []*entity.User
	StatusCode int
}

func (r *UserRepo) UsersRead(
	ctx context.Context,
	optionalLimit, optionalPage int,
) (res UsersReadResult) {
	arg := sqlc.UsersReadParams{}
	if optionalLimit > 0 {
		arg.Limit = pgtype.Int4{Int32: int32(optionalLimit), Valid: true}
	}
	if optionalPage > 0 && optionalLimit > 0 {
		arg.Offset = pgtype.Int4{
			Int32: int32((optionalPage - 1) * int(arg.Limit.Int32)),
			Valid: true,
		}
	}

	dbUsers, err := r.store.UsersRead(ctx, arg)
	if err != nil {
		res.Error = err
		res.StatusCode = http.StatusInternalServerError
		return
	}

	users := []*entity.User{}
	for _, u := range dbUsers {
		users = append(users, r.ToUserEntity(u))
	}

	res.StatusCode = http.StatusOK
	res.Users = users
	return
}

type UserReadResult struct {
	Error      error
	User       *entity.User
	StatusCode int
}

func (r *UserRepo) UserReadByID(ctx context.Context, userID string) (res UserReadResult) {
	dbuser, err := r.store.UserRead(ctx, uuid.MustParse(userID))
	if err != nil {
		res.StatusCode = http.StatusInternalServerError
		res.Error = err
		if isRecordNotFound(err) {
			res.StatusCode = http.StatusNotFound
			res.Error = repoerrors.ErrUserNotFound
			return
		}
		return
	}
	res.StatusCode = http.StatusOK
	res.User = r.ToUserEntity(dbuser)
	return
}

func (r *UserRepo) UserReadByEmail(ctx context.Context, email string) (res UserReadResult) {
	dbuser, err := r.store.UserReadByEmail(ctx, email)
	if err != nil {
		res.StatusCode = http.StatusInternalServerError
		res.Error = err
		if isRecordNotFound(err) {
			res.StatusCode = http.StatusNotFound
			res.Error = repoerrors.ErrUserNotFound
			return
		}
		return
	}
	res.StatusCode = http.StatusOK
	res.User = r.ToUserEntity(dbuser)
	return
}

func (r *UserRepo) UserReadMaster(ctx context.Context) (res UserReadResult) {
	dbuser, err := r.store.UserReadMaster(ctx)
	if err != nil {
		res.StatusCode = http.StatusInternalServerError
		res.Error = err
		if isRecordNotFound(err) {
			res.StatusCode = http.StatusNotFound
			res.Error = repoerrors.ErrUserNotFound
			return
		}
		return
	}
	res.StatusCode = http.StatusOK
	res.User = r.ToUserEntity(dbuser)
	return
}
