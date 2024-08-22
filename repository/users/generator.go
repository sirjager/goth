package users

import (
	"github.com/sirjager/goth/entity"
	"github.com/sirjager/goth/repository/users/sqlc"
	"github.com/sirjager/goth/vo"
)

func (r *repo) ToUserEntity(dbUser sqlc.User) *entity.User {
	return &entity.User{
		ID:                          vo.MustParseID(dbUser.ID.String()),
		Email:                       vo.MustParseEmail(dbUser.Email),
		Username:                    vo.MustParseUsername(dbUser.Username),
		Password:                    vo.MustParseHashedPassword(dbUser.Password),
		Verified:                    dbUser.Verified,
		Blocked:                     dbUser.Blocked,
		Provider:                    dbUser.Provider,
		GoogleID:                    dbUser.GoogleID,
		FullName:                    dbUser.FullName,
		FirstName:                   dbUser.FirstName,
		LastName:                    dbUser.LastName,
		AvatarURL:                   dbUser.AvatarUrl,
		PictureURL:                  dbUser.PictureUrl,
		CreatedAt:                   dbUser.CreatedAt,
		UpdatedAt:                   dbUser.UpdatedAt,
		Master:                      dbUser.Master,
	}
}
