package vo

import (
	"regexp"
)

type Email struct {
	baseValueObject[string]
}

var isValidEmail = regexp.MustCompile(
	`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
).MatchString

func NewEmail(value string) (*Email, error) {
	email := &Email{baseValueObject[string]{value}}
	if err := email.Validate(); err != nil {
		return nil, err
	}
	return email, nil
}

func (v *Email) Validate() error {
	if !isValidEmail(v.value) {
		return ErrInvalidEmail
	}
	return nil
}

// MustParseEmail returns email if valid or panics
func MustParseEmail(value string) *Email {
	email := &Email{baseValueObject[string]{value}}
	if err := email.Validate(); err != nil {
		panic(err)
	}
	return email
}
