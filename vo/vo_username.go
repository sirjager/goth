package vo

import (
	"errors"
	"fmt"
	"regexp"
)

type Username struct {
	baseValueObject[string]
}

const (
	UsernameMinLength = 3
	UsernameMaxLength = 60
)

var isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString

func NewUsername(value string) (*Username, error) {
	username := &Username{baseValueObject[string]{value}}
	if err := username.Validate(); err != nil {
		return nil, err
	}
	return username, nil
}

func (v *Username) Validate() error {
	if len(v.value) < UsernameMinLength {
		return fmt.Errorf("username must be at least %d characters long", UsernameMinLength)
	}
	if len(v.value) > UsernameMaxLength {
		return fmt.Errorf("username must be at most %d characters long", UsernameMaxLength)
	}
	if !isValidUsername(v.value) {
		return errors.New("username must only contain letters, numbers, and underscores")
	}
	return nil
}

func MustParseUsername(value string) *Username {
	username, err := NewUsername(value)
	if err != nil {
		panic(err)
	}
	return username
}
