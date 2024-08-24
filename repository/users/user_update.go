package users

import (
	"context"
	"net/http"

	"github.com/sirjager/goth/entity"
	repoerrors "github.com/sirjager/goth/repository/errors"
	"github.com/sirjager/goth/repository/users/sqlc"
	"github.com/sirjager/goth/vo"
)

func (r *repo) UserUpdate(ctx context.Context, u *entity.User) UserReadResult {
	return r._userReadCommon(func() (sqlc.User, error) {
		return r.store.UserUpdate(ctx, sqlc.UserUpdateParams{
			ID:         u.ID.Value(),
			FullName:   u.FullName,
			FirstName:  u.FirstName,
			Username:   u.Username.Value(),
			LastName:   u.LastName,
			PictureUrl: u.PictureURL,
			AvatarUrl:  u.AvatarURL,
		})
	})
}

func (r *repo) UserUpdateVerified(ctx context.Context, uid *vo.ID, status bool) UserReadResult {
	return r._userReadCommon(func() (sqlc.User, error) {
		return r.store.UserUpdateVerified(ctx, sqlc.UserUpdateVerifiedParams{
			ID:       uid.Value(),
			Verified: status,
		})
	})
}

func (r *repo) UserUpdatePassword(
	ctx context.Context,
	userID *vo.ID,
	password *vo.HashedPassword,
) UserReadResult {
	return r._userReadCommon(func() (sqlc.User, error) {
		return r.store.UserUpdatePassword(ctx, sqlc.UserUpdatePasswordParams{
			ID:       userID.Value(),
			Password: password.Value(),
		})
	})
}

type UserUpdatePasswordTxParams struct {
	UserID       *vo.ID
	Password     *vo.HashedPassword
	BeforeUpdate func() error
	AfterUpdate  func() error
}

func (r *repo) UserUpdatePasswordTx(
	ctx context.Context,
	params UserUpdatePasswordTxParams,
) UserReadResult {
	var res UserReadResult
	err := r.ExecTx(ctx, func() error {
		if params.BeforeUpdate != nil {
			if err := params.BeforeUpdate(); err != nil {
				return err
			}
		}

		dbUser, err := r.store.UserUpdatePassword(ctx, sqlc.UserUpdatePasswordParams{
			ID:       params.UserID.Value(),
			Password: params.Password.Value(),
		})
		if err != nil {
			return err
		}

		if params.AfterUpdate != nil {
			if err = params.AfterUpdate(); err != nil {
				return err
			}
		}
		res.User = r.ToUserEntity(dbUser)
		return nil
	})
	if err != nil {
		if isRecordNotFound(err) {
			res.Error = repoerrors.ErrUserNotFound
			res.StatusCode = http.StatusNotFound
			return res
		}
		if isUniqueViolation(err) {
			res.Error = repoerrors.ErrUniqueKeyViolation
			res.StatusCode = http.StatusConflict
			return res
		}
		res.Error = err
		res.StatusCode = http.StatusInternalServerError
		return res
	}

	res.StatusCode = http.StatusOK
	return res
}
