package users

import (
	"context"

	"github.com/google/uuid"

	"github.com/sirjager/goth/entity"
	repoerrors "github.com/sirjager/goth/repository/errors"
)

func (r *repo) UserReadByID(ctx context.Context, userID string) (*entity.User, error) {
	dbuser, err := r.store.UserRead(ctx, uuid.MustParse(userID))
	if err != nil {
		if isRecordNotFound(err) {
			return nil, repoerrors.ErrUserNotFound
		}
		return nil, err
	}
	return r.toUserEntity(dbuser), nil
}

func (r *repo) UserReadByEmail(ctx context.Context, email string) (*entity.User, error) {
	dbuser, err := r.store.UserReadByEmail(ctx, email)
	if err != nil {
		if isRecordNotFound(err) {
			return nil, repoerrors.ErrUserNotFound
		}
		return nil, err
	}
	return r.toUserEntity(dbuser), nil
}
