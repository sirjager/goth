package payload

import (
	"encoding/json"
	"time"

	"github.com/sirjager/goth/entity"
)

type BaseAuthPayload struct {
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UserID    string    `json:"userID,omitempty"`
	SessionID string    `json:"sessionID,omitempty"`
	TokenType string    `json:"tokenType,omitempty"`
}

func (b *RefreshToken) IsEqual(other RefreshToken) bool {
	if b.CreatedAt != other.CreatedAt {
		return false
	}
	if b.TokenType != other.TokenType {
		return false
	}
	if b.SessionID != other.SessionID {
		return false
	}
	if b.UserID != other.UserID {
		return false
	}
	return true
}


func (b *AccessToken) IsEqual(other AccessToken) bool {
	if b.CreatedAt != other.CreatedAt {
		return false
	}
	if b.TokenType != other.TokenType {
		return false
	}
	if b.SessionID != other.SessionID {
		return false
	}
	if b.UserID != other.UserID {
		return false
	}
	return true
}

// Marshal for AccessToken
func (a *BaseAuthPayload) Marshal() ([]byte, error) {
	return json.Marshal(a)
}

// Unmarshal for AccessToken
func (a *BaseAuthPayload) Unmarshal(data []byte) error {
	return json.Unmarshal(data, a)
}

type AccessToken struct {
	BaseAuthPayload
}

type RefreshToken struct {
	BaseAuthPayload
}

func NewAccessPayload(user *entity.User, sessionID string) *AccessToken {
	return &AccessToken{
		BaseAuthPayload: BaseAuthPayload{
			CreatedAt: time.Now(),
			SessionID: sessionID,
			UserID:    user.ID.Value().String(),
			TokenType: TypeAccess,
		},
	}
}

func NewRefreshPayload(user *entity.User, sessionID string) *RefreshToken {
	return &RefreshToken{
		BaseAuthPayload: BaseAuthPayload{
			CreatedAt: time.Now(),
			SessionID: sessionID,
			UserID:    user.ID.Value().String(),
			TokenType: TypeRefresh,
		},
	}
}
