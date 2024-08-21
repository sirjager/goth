package users

import (
	"context"

	"github.com/sirjager/goth/entity"
	"github.com/sirjager/goth/repository/users/sqlc"
	"github.com/sirjager/goth/vo"
)

func (r *userRepo) UserUpdate(ctx context.Context, u *entity.User) UserReadResult {
	return r._userReadCommon(func() (sqlc.User, error) {
		return r.store.UserUpdate(ctx, sqlc.UserUpdateParams{
			ID:         u.ID.Value(),
			FullName:   u.FullName,
			FirstName:  u.FirstName,
			LastName:   u.LastName,
			PictureUrl: u.PictureURL,
			AvatarUrl:  u.AvatarURL,
		})
	})
}

func (r *userRepo) UserUpdateVerified(ctx context.Context, uid *vo.ID, status bool) UserReadResult {
	return r._userReadCommon(func() (sqlc.User, error) {
		return r.store.UserUpdateVerified(ctx, sqlc.UserUpdateVerifiedParams{
			ID:       uid.Value(),
			Verified: status,
		})
	})
}

func (r *userRepo) UserUpdatePassword(ctx context.Context, uid *vo.ID, pass *vo.HashedPassword) UserReadResult {
	return r._userReadCommon(func() (sqlc.User, error) {
		return r.store.UserUpdatePassword(ctx, sqlc.UserUpdatePasswordParams{
			ID:       uid.Value(),
			Password: pass.Value(),
		})
	})
}
