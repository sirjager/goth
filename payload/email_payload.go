package payload

import (
	"time"

	"github.com/sirjager/gopkg/utils"

	"github.com/sirjager/goth/entity"
)

type BaseEmailPayload struct {
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UserID    string    `json:"userID,omitempty"`
	Email     string    `json:"email,omitempty"`
	Code      string    `json:"code,omitempty"`
	TokenType string    `json:"tokenType,omitempty"`
}

func (b *BaseEmailPayload) IsExpired(maxDuration time.Duration) bool {
	return time.Since(b.CreatedAt) >= maxDuration
}

func (b *BaseEmailPayload) IsValidType(validType string) bool {
	return b.TokenType == validType
}

type VerifyEmail struct {
	BaseEmailPayload
}

type ResetPassword struct {
	BaseEmailPayload
}

type UserDelete struct {
	BaseEmailPayload
}

type ChangeEmail struct {
	BaseEmailPayload
}

func NewChangeEmailPayload(user *entity.User) *ChangeEmail {
	return &ChangeEmail{
		BaseEmailPayload: BaseEmailPayload{
			UserID:    user.ID.Value().String(),
			CreatedAt: time.Now(),
			Email:     user.Email.Value(),
			Code:      utils.RandomNumberAsString(6),
			TokenType: TypeChangeEmail,
		},
	}
}

func NewUserDeletePayload(user *entity.User) *UserDelete {
	return &UserDelete{
		BaseEmailPayload: BaseEmailPayload{
			UserID:    user.ID.Value().String(),
			CreatedAt: time.Now(),
			Email:     user.Email.Value(),
			Code:      utils.RandomNumberAsString(6),
			TokenType: TypeUserDelete,
		},
	}
}

func NewVerifyEmailPayload(user *entity.User) *VerifyEmail {
	return &VerifyEmail{
		BaseEmailPayload: BaseEmailPayload{
			UserID:    user.ID.Value().String(),
			CreatedAt: time.Now(),
			Email:     user.Email.Value(),
			Code:      utils.RandomNumberAsString(6),
			TokenType: TypeVerifyEmail,
		},
	}
}

func NewResetPasswordPayload(user *entity.User) *ResetPassword {
	return &ResetPassword{
		BaseEmailPayload: BaseEmailPayload{
			UserID:    user.ID.Value().String(),
			CreatedAt: time.Now(),
			Email:     user.Email.Value(),
			Code:      utils.RandomNumberAsString(6),
			TokenType: TypeResetPassword,
		},
	}
}
