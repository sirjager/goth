package api

import (
	"time"

	"github.com/markbates/goth"

	"github.com/sirjager/goth/entity"
)

type UserResponse struct {
	User User `json:"user,omitempty"`
} // @name UserResponse

type User struct {
	ID        string    `json:"id,omitempty"`
	Email     string    `json:"email,omitempty"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Verified  bool      `json:"verified,omitempty"`
	Blocked   bool      `json:"blocked,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
} // @name User

func EntityToUser(user *entity.User) User {
	return User{
		ID:        user.ID,
		Email:     user.Email,
		Verified:  user.Verified,
		Blocked:   user.Blocked,
		AvatarURL: user.AvatarURL,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func GothUserToEntityUser(user goth.User) *entity.User {
	userEntity := &entity.User{
		ID:        user.UserID,
		Email:     user.Email,
		Provider:  user.Provider,
		Name:      user.Name,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Location:  user.Location,
		NickName:  user.NickName,
		GoogleID:  user.UserID,
		Verified:  false,
		AvatarURL: user.AvatarURL,
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
