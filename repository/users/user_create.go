package users

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/sirjager/goth/entity"
	repoerrors "github.com/sirjager/goth/repository/errors"
	"github.com/sirjager/goth/repository/users/sqlc"
)

func (r *repo) UserCreate(c context.Context, u *entity.User) (*entity.User, error) {
	dbuser, err := r.store.UserCreate(c, sqlc.UserCreateParams{
		ID:        uuid.New(),
		Email:     u.Email,
		Verified:  u.Verified,
		Blocked:   u.Blocked,
		Provider:  u.Provider,
		GoogleID:  u.GoogleID,
		Name:      u.Name,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		NickName:  u.NickName,
		AvatarUrl: u.AvatarURL,
		Location:  u.Location,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, repoerrors.ErrUserAlreadyExists
		}
		return nil, nil
	}
	return r.toUserEntity(dbuser), nil
}
