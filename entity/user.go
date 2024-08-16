// Package entity provides entity models such as User
package entity

import (
	"time"

	"github.com/sirjager/goth/vo"
)

type User struct {
	CreatedAt  time.Time          `json:"created_at,omitempty"`
	UpdatedAt  time.Time          `json:"updated_at,omitempty"`
	Email      *vo.Email          `json:"email,omitempty"`
	Password   *vo.HashedPassword `json:"password,omitempty"`
	Username   *vo.Username       `json:"username,omitempty"`
	ID         *vo.ID             `json:"id,omitempty"`
	LastName   string             `json:"last_name,omitempty"`
	NickName   string             `json:"nick_name,omitempty"`
	FirstName  string             `json:"first_name,omitempty"`
	GoogleID   string             `json:"google_id,omitempty"`
	AvatarURL  string             `json:"avatar_url,omitempty"`
	PictureURL string             `json:"picture_url,omitempty"`
	Location   string             `json:"location,omitempty"`
	Name       string             `json:"name,omitempty"`
	Provider   string             `json:"provider,omitempty"`
	Master     bool               `json:"master,omitempty"`
	Verified   bool               `json:"verified,omitempty"`
	Blocked    bool               `json:"blocked,omitempty"`
}

// Profile is for public display, and can be used to send out diretcly.
// It strips out confidentials and private information
type Profile struct {
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	PictureURL string    `json:"picture_url"`
	Name       string    `json:"name"`
	Username   string    `json:"username"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Verified   bool      `json:"verified"`
	Blocked    bool      `json:"blocked"`
} // @name User

func (user *User) Profile() *Profile {
	profile := &Profile{
		ID:         user.ID.Value().String(),
		Email:      user.Email.Value(),
		Username:   user.Username.Value(),
		Verified:   user.Verified,
		Blocked:    user.Blocked,
		Name:       user.Name,
		PictureURL: user.PictureURL,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
	if profile.PictureURL == "" {
		profile.PictureURL = user.AvatarURL
	}
	return profile
}
