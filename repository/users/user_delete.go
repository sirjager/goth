package users

import (
	"context"
	"net/http"

	repoerrors "github.com/sirjager/goth/repository/errors"
	"github.com/sirjager/goth/vo"
)

type UserDeleteResult struct {
	Error      error
	StatusCode int
}

func (r *repo) UserDelete(ctx context.Context, userID *vo.ID) (res UserDeleteResult) {
	_, err := r.store.UserDelete(ctx, userID.Value())
	if err != nil {
		res.Error = err
		res.StatusCode = http.StatusInternalServerError
		if isRecordNotFound(err) {
			res.StatusCode = http.StatusNotFound
			res.Error = repoerrors.ErrUserNotFound
		}
		return
	}

	res.StatusCode = http.StatusOK
	return
}

type UserDeleteTxParams struct {
	UserID       *vo.ID
	BeforeUpdate func() error
	AfterUpdate  func() error
}

func (r *repo) UserDeleteTx(ctx context.Context, params UserDeleteTxParams) UserDeleteResult {
	var res UserDeleteResult
	err := r.ExecTx(ctx, func() error {
		if params.BeforeUpdate != nil {
			if err := params.BeforeUpdate(); err != nil {
				return err
			}
		}
		_, err := r.store.UserDelete(ctx, params.UserID.Value())
		if err != nil {
			return err
		}

		if params.AfterUpdate != nil {
			if err = params.AfterUpdate(); err != nil {
				return err
			}
		}
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
