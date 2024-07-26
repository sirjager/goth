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

func (r *repo) UserCreate(c context.Context, u *entity.User) (res UserReadResult) {
	dbuser, err := r.store.UserCreate(c, sqlc.UserCreateParams{
		ID:         uuid.New(),
		Email:      u.Email,
		Verified:   u.Verified,
		Blocked:    u.Blocked,
		Provider:   u.Provider,
		GoogleID:   u.GoogleID,
		Name:       u.Name,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		NickName:   u.NickName,
		AvatarUrl:  u.AvatarURL,
		PictureUrl: u.PictureURL,
		Master:     u.Master,
		Location:   u.Location,
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

	res.StatusCode = 201
	res.User = r.toUserEntity(dbuser)
	return
}
