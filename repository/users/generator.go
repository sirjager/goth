package users

import (
	"github.com/google/uuid"

	"github.com/sirjager/goth/entity"
	"github.com/sirjager/goth/repository/users/sqlc"
)

func (r *repo) toUserEntity(dbUser sqlc.User) *entity.User {
	return &entity.User{
		ID:         dbUser.ID.String(),
		Email:      dbUser.Email,
		Verified:   dbUser.Verified,
		Blocked:    dbUser.Blocked,
		Provider:   dbUser.Provider,
		GoogleID:   dbUser.GoogleID,
		Name:       dbUser.Name,
		FirstName:  dbUser.FirstName,
		LastName:   dbUser.LastName,
		NickName:   dbUser.NickName,
		Location:   dbUser.Location,
		AvatarURL:  dbUser.AvatarUrl,
		CreatedAt:  dbUser.CreatedAt,
		UpdatedAt:  dbUser.UpdatedAt,
	}
}

func (r *repo) ToDatabaseUser(user *entity.User) sqlc.User {
	return sqlc.User{
		ID:         uuid.MustParse(user.ID),
		Email:      user.Email,
		Verified:   user.Verified,
		Blocked:    user.Blocked,
		Provider:   user.Provider,
		GoogleID:   user.GoogleID,
		Name:       user.Name,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		NickName:   user.NickName,
		Location:   user.Location,
		AvatarUrl:  user.AvatarURL,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}
