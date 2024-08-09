package users

import (
	"github.com/sirjager/goth/entity"
	"github.com/sirjager/goth/repository/users/sqlc"
	"github.com/sirjager/goth/vo"
)

func (r *UserRepo) ToUserEntity(dbUser sqlc.User) *entity.User {
	return &entity.User{
		ID:         vo.MustParseID(dbUser.ID.String()),
		Email:      vo.MustParseEmail(dbUser.Email),
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
		PictureURL: dbUser.PictureUrl,
		CreatedAt:  dbUser.CreatedAt,
		UpdatedAt:  dbUser.UpdatedAt,
		Master:     dbUser.Master,
	}
}

func (r *UserRepo) ToDatabaseUser(user *entity.User) sqlc.User {
	return sqlc.User{
		ID:        user.ID.Value(),
		Email:     user.Email.Value(),
		Verified:  user.Verified,
		Blocked:   user.Blocked,
		Provider:  user.Provider,
		GoogleID:  user.GoogleID,
		Name:      user.Name,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Master:    user.Master,
		NickName:  user.NickName,
		Location:  user.Location,
		AvatarUrl: user.AvatarURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
