package users

import (
	"context"
	"net/http"

	"github.com/sirjager/goth/entity"
	repoerrors "github.com/sirjager/goth/repository/errors"
	"github.com/sirjager/goth/repository/users/sqlc"
)

func (r *UserRepo) UserUpdate(ctx context.Context, u *entity.User) (res UserReadResult) {
	updated, err := r.store.UserUpdate(ctx, sqlc.UserUpdateParams{
		ID:         u.ID.Value(),
		Name:       u.Name,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		NickName:   u.NickName,
		PictureUrl: u.PictureURL,
		AvatarUrl:  u.AvatarURL,
	})
	if err != nil {
		res.Error = err
		res.StatusCode = http.StatusInternalServerError
		if isRecordNotFound(err) {
			res.StatusCode = http.StatusNotFound
			res.Error = repoerrors.ErrUserNotFound
			return
		}
		return
	}

	res.StatusCode = http.StatusOK
	res.User = r.ToUserEntity(updated)
	return
}
