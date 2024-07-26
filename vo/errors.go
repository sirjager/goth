package vo

import (
	"errors"
)

var (
	ErrInvalidEmail         = errors.New("invalid email")
	ErrInvalidEmailDomain   = errors.New("invalid email domain")
	ErrInvalidEmailUsername = errors.New("invalid email username")
)
