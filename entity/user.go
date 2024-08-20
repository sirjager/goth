// Package entity provides entity models such as User
package entity

import (
	"time"

	"github.com/sirjager/goth/vo"
)

type User struct {
	CreatedAt  time.Time          `json:"createdAt,omitempty"  bson:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt,omitempty"  bson:"updatedAt"`
	Email      *vo.Email          `json:"email,omitempty"      bson:"email"`
	Password   *vo.HashedPassword `json:"password,omitempty"   bson:"password"`
	Username   *vo.Username       `json:"username,omitempty"   bson:"username"`
	ID         *vo.ID             `json:"id,omitempty"         bson:"id"`
	Provider   string             `json:"provider,omitempty"   bson:"provider"`
	AvatarURL  string             `json:"avatarURL,omitempty"  bson:"avatarURL"`
	FullName   string             `json:"fullName,omitempty"   bson:"fullName"`
	FirstName  string             `json:"firstName,omitempty"  bson:"firstName"`
	GoogleID   string             `json:"googleID,omitempty"   bson:"googleID"`
	LastName   string             `json:"lastName,omitempty"   bson:"lastName"`
	PictureURL string             `json:"pictureURL,omitempty" bson:"pictureURL"`
	Blocked    bool               `json:"blocked,omitempty"    bson:"blocked"`
	Verified   bool               `json:"verified,omitempty"   bson:"verified"`
	Master     bool               `json:"master,omitempty"     bson:"master"`
}

// Profile is for public display, and can be used to send out diretcly.
// It strips out confidentials and private information
type Profile struct {
	CreatedAt  time.Time `json:"createdAt,omitempty"  bson:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt,omitempty"  bson:"updatedAt"`
	ID         string    `json:"id,omitempty"         bson:"id"`
	Email      string    `json:"email,omitempty"      bson:"email"`
	PictureURL string    `json:"pictureURL,omitempty" bson:"pictureURL"`
	Username   string    `json:"username,omitempty"   bson:"username"`
	FullName   string    `json:"fullName,omitempty"   bson:"fullName"`
	FirstName  string    `json:"firstName,omitempty"  bson:"firstName"`
	LastName   string    `json:"lastName,omitempty"   bson:"lastName"`
	Verified   bool      `json:"verified,omitempty"   bson:"verified"`
	Blocked    bool      `json:"blocked,omitempty"    bson:"blocked"`
} // @name User

func (user *User) Profile() *Profile {
	profile := &Profile{
		ID:         user.ID.Value().String(),
		Email:      user.Email.Value(),
		Username:   user.Username.Value(),
		Verified:   user.Verified,
		Blocked:    user.Blocked,
		FullName:   user.FullName,
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
