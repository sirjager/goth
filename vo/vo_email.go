package vo

import (
	"regexp"

	"github.com/sirjager/goth/utils"
)

type Email struct {
	value string
}

var (
	isValidEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString
)

func NewEmail(value string) (*Email, error) {
	email := &Email{value}
	if err := email.Validate(); err != nil {
		return nil, err
	}
	return email, nil
}

func GenerateTestEmail() *Email {
	return &Email{utils.RandomEmail()}
}

func (v *Email) IsEqual(other *Email) bool {
	return v.value == other.value
}

// MustParseEmail returns email if valid or panics
func MustParseEmail(value string) *Email {
	email := &Email{value}
	if err := email.Validate(); err != nil {
		panic(err)
	}
	return email
}

func (e *Email) Value() string {
	return string(e.value)
}

func (v *Email) Validate() error {
	if !isValidEmail(v.value) {
		return ErrInvalidEmail
	}
	return nil
}
