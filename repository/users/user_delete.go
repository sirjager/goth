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

func (r *UserRepo) UserDelete(ctx context.Context, userID *vo.ID) (res UserDeleteResult) {
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
