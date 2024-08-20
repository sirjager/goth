package api

import "errors"

var (
	errInvalidEmailVerificationCode = errors.New("invalid email verification code")
	errEmailNotVerified             = errors.New("email not verified")
)
