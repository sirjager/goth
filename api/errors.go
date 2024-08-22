package api

import "errors"

var (
	errInvalidCode      = errors.New("error invalid code")
	errEmailNotVerified = errors.New("email not verified")
)
