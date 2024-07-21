package repoerrors

import (
	"errors"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrFailedToMigrate = errors.New("failed to migrate database")
)
