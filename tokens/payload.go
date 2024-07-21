package tokens

import (
	"errors"
	"time"

	"github.com/sirjager/goth/vo"
)

var (
	// ErrExpiredToken is returned when a token has expired
	ErrExpiredToken = errors.New("expired token")

	// ErrInvalidToken is returned when a token is invalid
	ErrInvalidToken = errors.New("invalid token")
)

// Data contains the payload data of the token
type Data struct {
	SessionID *vo.ID `json:"session_id,omitempty"`
	UserID    *vo.ID `json:"user_id,omitempty"`
	ClientIP  string `json:"client_ip,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	Type      string `json:"type,omitempty"`

	// Code stores otps for emails,
	// Never in tokens, or anywhere else
	Code string `json:"code,omitempty"`
	// UserEmail stores email for emails
	// Never sent in tokens, or anywhere else
	UserEmail *vo.Email `json:"user_email,omitempty"`
	NewEmail  *vo.Email `json:"new_email,omitempty"`
}

// Payload contains the payload data of the token
type Payload struct {
	IssuedAt  time.Time `json:"iat,omitempty"`
	ExpiresAt time.Time `json:"expires,omitempty"`
	Data      *Data     `json:"payload,omitempty"`
	ID        *vo.ID    `json:"id,omitempty"`
}

// NewPayload creates a new payload for a specific username and duration
func NewPayload(data *Data, duration time.Duration) (*Payload, error) {
	id, err := vo.NewID()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        id,
		Data:      data,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
	return payload, nil
}

// Valid checks if the token payload is not expired
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}
