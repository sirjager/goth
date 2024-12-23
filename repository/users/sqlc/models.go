// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package sqlc

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Verified   bool      `json:"verified"`
	Blocked    bool      `json:"blocked"`
	Provider   string    `json:"provider"`
	GoogleID   string    `json:"google_id"`
	FullName   string    `json:"full_name"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	AvatarUrl  string    `json:"avatar_url"`
	PictureUrl string    `json:"picture_url"`
	Master     bool      `json:"master"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
