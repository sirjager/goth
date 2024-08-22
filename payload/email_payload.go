package payload

import (
	"fmt"
	"time"
)

type EmailType int

const (
	EmailVerification EmailType = iota
	PasswordReset
	EmailChange
	UserDeletion
)

// String method to convert the enum to a string for easy printing
func (e EmailType) String() string {
	emailTypes := []string{"EmailVerification", "PasswordReset", "EmailChange", "UserDeletion"}
	if e < EmailVerification || e > UserDeletion {
		return "Unknown"
	}
	return emailTypes[e]
}

type EmailPayload struct {
	CreatedAt time.Time     `json:"createdAt,omitempty"`
	Body      string        `json:"body,omitempty"`
	Subject   string        `json:"subject,omitempty"`
	Email     string        `json:"email,omitempty"`
	Code      string        `json:"code,omitempty"`
	CacheKey  string        `json:"cacheKey,omitempty"`
	Type      EmailType     `json:"type,omitempty"`
	CacheExp  time.Duration `json:"cacheExp,omitempty"`
}

func EmailKey(email string, emailType EmailType) string {
	return fmt.Sprintf("send:%d:%s", emailType, email)
}
