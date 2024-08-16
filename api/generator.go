package api

import (
	"github.com/markbates/goth"

	"github.com/sirjager/goth/entity"
	"github.com/sirjager/goth/vo"
)

func EntitiesToProfiles(users []*entity.User) []*entity.Profile {
	result := []*entity.Profile{}
	for _, u := range users {
		result = append(result, u.Profile())
	}
	return result
}

func GothUserToEntityUser(user goth.User) *entity.User {
	userEntity := &entity.User{
		// ID:         vo.MustParseID(user.UserID), -- this is google id, it is slightly different, so we have to skip this
		Email:      vo.MustParseEmail(user.Email),
		Username:   vo.GenerateUsername(),
		Password:   &vo.HashedPassword{}, // this user does not have password
		Provider:   user.Provider,
		Name:       user.Name,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Location:   user.Location,
		NickName:   user.NickName,
		GoogleID:   user.UserID,
		Verified:   false,
		PictureURL: user.AvatarURL, //  this will be set by user,
		AvatarURL:  user.AvatarURL, // this comes from auth provider
	}

	if user.RawData["verified_email"] != nil {
		userEntity.Verified = user.RawData["verified_email"].(bool)
	}
	return userEntity
}

func GothUserToUser(gothuser goth.User) *entity.Profile {
	user := GothUserToEntityUser(gothuser)
	return user.Profile()
}
