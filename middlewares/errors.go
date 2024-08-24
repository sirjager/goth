package mw

import "errors"

const (
	insufficientPermissions = "insufficient permissions"
	unauthorized            = "Unauthorized"
)

var (
	ErrUnAuthorized     = errors.New("Unauthorized")
	ErrEmailNotVerified = errors.New("email not verified")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
)
