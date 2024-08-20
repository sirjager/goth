package mw

import (
	"time"

	"github.com/sirjager/gopkg/utils"

	"github.com/sirjager/goth/entity"
	"github.com/sirjager/goth/vo"
)

const (
	TokenTypeAccess            = "0"
	TokenTypeRefresh           = "1"
	TokenTypeEmailVerification = "2"
)

type TokenCustomPayload struct {
	CreatedAt             time.Time `json:"createdAt,omitempty"`
	UserID                string    `json:"userID,omitempty"`
	UserEmail             string    `json:"userEmail,omitempty"`
	SessionID             string    `json:"sessionID,omitempty"`
	TokenType             string    `json:"tokenType,omitempty"`
	EmailVerificationCode string    `json:"emailVerificationCode,omitempty"`
}

func NewAccessPayload(userID *vo.ID, sessionID string) *TokenCustomPayload {
	return &TokenCustomPayload{
		UserID:    userID.Value().String(),
		SessionID: sessionID,
		CreatedAt: time.Now(),
		TokenType: TokenTypeAccess,
	}
}

func NewRefreshPayload(userID *vo.ID, sessionID string) *TokenCustomPayload {
	return &TokenCustomPayload{
		UserID:    userID.Value().String(),
		SessionID: sessionID,
		CreatedAt: time.Now(),
		TokenType: TokenTypeRefresh,
	}
}

func NewEmailVerificationPayload(user *entity.User) *TokenCustomPayload {
	return &TokenCustomPayload{
		UserID:                user.ID.Value().String(),
		UserEmail:             user.Email.Value(),
		CreatedAt:             time.Now(),
		TokenType:             TokenTypeEmailVerification,
		EmailVerificationCode: utils.RandomNumberAsString(6),
	}
}
