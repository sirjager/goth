package users

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/sirjager/goth/entity"
	repoerrors "github.com/sirjager/goth/repository/errors"
	"github.com/sirjager/goth/repository/users/sqlc"
)

func (r *userRepo) UserCreate(c context.Context, user *entity.User) (res UserReadResult) {
	dbuser, err := r.store.UserCreate(c, sqlc.UserCreateParams{
		ID:         uuid.New(),
		Email:      user.Email.Value(),
		Password:   user.Password.Value(),
		Username:   user.Username.Value(),
		Verified:   user.Verified,
		Blocked:    user.Blocked,
		Provider:   user.Provider,
		GoogleID:   user.GoogleID,
		FullName:   user.FullName,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		AvatarUrl:  user.AvatarURL,
		PictureUrl: user.PictureURL,
		Master:     user.Master,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	if err != nil {
		res.Error = err
		res.StatusCode = http.StatusInternalServerError
		if isUniqueViolation(err) {
			res.StatusCode = http.StatusConflict
			res.Error = repoerrors.ErrUserAlreadyExists
			return
		}
		return
	}

	res.StatusCode = http.StatusCreated
	res.User = r.ToUserEntity(dbuser)
	return
}
