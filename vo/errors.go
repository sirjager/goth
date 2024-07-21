package vo

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrShortPassword = fmt.Errorf(
		"password too short, must be at least %d characters long",
		PasswordMinLength,
	)
	ErrPasswordTooLong = fmt.Errorf(
		"password too long, must be at most %d characters long",
		PasswordMaxLength,
	)
	ErrInvalidPasword = errors.New(
		"password must contain at least one uppercase letter, one lowercase letter, one number, and one special character",
	)

	ErrInvalidEmail         = errors.New("invalid email")
	ErrInvalidEmailDomain   = errors.New("invalid email domain")
	ErrInvalidEmailUsername = errors.New("invalid email username")

	ErrPasswordHashesDoNotMatch = errors.New("password hashes do not match")

	ErrInvalidUsername = errors.New(
		"username must only contain letters, numbers, and underscores",
	)

	ErrMismatchedHashAndPassword = bcrypt.ErrMismatchedHashAndPassword
)
