// Package entity provides entity models such as User
package entity

import (
	"time"

	"github.com/sirjager/goth/vo"
)

type User struct {
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	Provider   string    `json:"provider,omitempty"`
	FirstName  string    `json:"first_name,omitempty"`
	LastName   string    `json:"last_name,omitempty"`
	NickName   string    `json:"nick_name,omitempty"`
	ID         *vo.ID    `json:"id,omitempty"`
	GoogleID   string    `json:"google_id,omitempty"`
	AvatarURL  string    `json:"avatar_url,omitempty"`
	PictureURL string    `json:"picture_url,omitempty"`
	Location   string    `json:"location,omitempty"`
	Name       string    `json:"name,omitempty"`
	Email      *vo.Email `json:"email,omitempty"`
	Master     bool      `json:"master,omitempty"`
	Verified   bool      `json:"verified,omitempty"`
	Blocked    bool      `json:"blocked,omitempty"`
}
