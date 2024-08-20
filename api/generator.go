package api

import (
	"fmt"
	"strings"

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

func GothUserToEntityUser(gothUser goth.User) *entity.User {
	user := &entity.User{
		// ID:         vo.MustParseID(user.UserID), -- this is google id, it is slightly different, so we have to skip this
		Email:      vo.MustParseEmail(gothUser.Email),
		Username:   vo.MustParseUsername(strings.Split(gothUser.Email, "@")[0]),
		Password:   &vo.HashedPassword{}, // this user does not have password
		Provider:   gothUser.Provider,
		FullName:   fmt.Sprintf("%s %s", gothUser.FirstName, gothUser.LastName),
		FirstName:  gothUser.FirstName,
		LastName:   gothUser.LastName,
		GoogleID:   gothUser.UserID,
		Verified:   false,
		PictureURL: gothUser.AvatarURL, //  this will be set by user,
		AvatarURL:  gothUser.AvatarURL, // this comes from auth provider
	}

	if gothUser.RawData["verified_email"] != nil {
		if _emailVerified, ok := gothUser.RawData["verified_email"].(bool); ok {
			user.Verified = _emailVerified
		}
	}
	return user
}

func GothUserToUser(gothuser goth.User) *entity.Profile {
	user := GothUserToEntityUser(gothuser)
	return user.Profile()
}
