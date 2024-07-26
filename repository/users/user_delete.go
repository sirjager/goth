package users

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	repoerrors "github.com/sirjager/goth/repository/errors"
)

type UserDeleteResult struct {
	Error      error
	StatusCode int
}

func (r *UserRepo) UserDelete(ctx context.Context, userID string) (res UserDeleteResult) {
	_, err := r.store.UserDelete(ctx, uuid.MustParse(userID))
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
