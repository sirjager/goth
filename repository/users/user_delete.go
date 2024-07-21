package users

import (
	"context"

	"github.com/google/uuid"

	repoerrors "github.com/sirjager/goth/repository/errors"
)

func (r *repo) UserDelete(ctx context.Context, userID string) error {
	_, err := r.store.UserDelete(ctx, uuid.MustParse(userID))
	if err != nil {
		if isRecordNotFound(err) {
			return repoerrors.ErrUserNotFound
		}
		return err
	}

	return nil
}
