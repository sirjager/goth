package users

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/sirjager/goth/entity"
	repoerrors "github.com/sirjager/goth/repository/errors"
	"github.com/sirjager/goth/repository/users/sqlc"
	"github.com/sirjager/goth/vo"
)

type UserReadResult struct {
	Error      error
	User       *entity.User
	StatusCode int
}

type UsersReadResult struct {
	Error      error
	Users      []*entity.User
	StatusCode int
}

func (r *repo) UserGetByID(ctx context.Context, userID *vo.ID) (res UserReadResult) {
	return r._userReadCommon(func() (sqlc.User, error) {
		return r.store.UserRead(ctx, userID.Value())
	})
}

func (r *repo) UserGetByEmail(ctx context.Context, email *vo.Email) (res UserReadResult) {
	return r._userReadCommon(func() (sqlc.User, error) {
		return r.store.UserReadByEmail(ctx, email.Value())
	})
}

func (r *repo) UserGetByUsername(ctx context.Context, u *vo.Username) (res UserReadResult) {
	return r._userReadCommon(func() (sqlc.User, error) {
		return r.store.UserReadByUsername(ctx, u.Value())
	})
}

func (r *repo) UserGetMaster(ctx context.Context) (res UserReadResult) {
	return r._userReadCommon(func() (sqlc.User, error) {
		return r.store.UserReadMaster(ctx)
	})
}

func (r *repo) UserGetAll(ctx context.Context, limit, page int) (res UsersReadResult) {
	arg := sqlc.UsersReadParams{}
	if limit > 0 {
		arg.Limit = pgtype.Int4{Int32: int32(limit), Valid: true}
	}
	if page > 0 && limit > 0 {
		arg.Offset = pgtype.Int4{
			Int32: int32((page - 1) * int(arg.Limit.Int32)),
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

func (r *repo) _userReadCommon(call func() (sqlc.User, error)) (res UserReadResult) {
	dbuser, err := call()
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
