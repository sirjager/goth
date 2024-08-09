package api

import (
	"time"

	"github.com/markbates/goth"

	"github.com/sirjager/goth/entity"
	"github.com/sirjager/goth/vo"
)

type User struct {
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	ID         string    `json:"id,omitempty"`
	Email      string    `json:"email,omitempty"`
	PictureURL string    `json:"picture_url,omitempty"`
	Name       string    `json:"name,omitempty"`
	FirstName  string    `json:"first_name,omitempty"`
	LastName   string    `json:"last_name,omitempty"`
	Verified   bool      `json:"verified,omitempty"`
	Blocked    bool      `json:"blocked,omitempty"`
} // @name User

func EntityToUser(user *entity.User) User {
	u := User{
		ID:         user.ID.Value().String(),
		Email:      user.Email.Value(),
		Verified:   user.Verified,
		Blocked:    user.Blocked,
		Name:       user.Name,
		PictureURL: user.AvatarURL,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
	if user.PictureURL != "" {
		u.PictureURL = user.PictureURL
	}
	return u
}

func EntitiesToUsers(users []*entity.User) []User {
	result := []User{}
	for _, u := range users {
		result = append(result, EntityToUser(u))
	}
	return result
}

func GothUserToEntityUser(user goth.User) *entity.User {
	userEntity := &entity.User{
		// ID:         vo.MustParseID(user.UserID), -- this is google id, it is slightly different, so we have to skip this
		Email:      vo.MustParseEmail(user.Email),
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

func GothUserToUser(user goth.User) User {
	entityUser := GothUserToEntityUser(user)
	return EntityToUser(entityUser)
}
