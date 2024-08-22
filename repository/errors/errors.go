package repoerrors

import (
	"errors"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUniqueKeyViolation = errors.New("uniqe key violation")

	ErrFailedToMigrate = errors.New("failed to migrate database")

	ErrRoleNotFound = errors.New("role not found")
)
