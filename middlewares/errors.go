package mw

import "errors"

const (
	emailNotVerified        = "Email not verified"
	insufficientPermissions = "Insufficient permissions"
	unauthorized            = "Unauthorized"
)


var (
	ErrUnAuthorized = errors.New("Unauthorized")
)
