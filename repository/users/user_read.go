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

func (r *repo) UsersRead(
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
		users = append(users, r.toUserEntity(u))
	}

	res.StatusCode = http.StatusOK
	res.Users = users
	return
}

func (r *repo) UserReadByID(ctx context.Context, userID string) (res UserReadResult) {
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
	res.User = r.toUserEntity(dbuser)
	return
}

func (r *repo) UserReadByEmail(ctx context.Context, email string) (res UserReadResult) {
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
	res.User = r.toUserEntity(dbuser)
	return
}

func (r *repo) UserReadMaster(ctx context.Context) (res UserReadResult) {
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
	res.User = r.toUserEntity(dbuser)
	return
}
