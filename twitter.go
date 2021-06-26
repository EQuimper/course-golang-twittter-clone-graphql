package twitter

import "errors"

var (
	ErrBadCredentials     = errors.New("email/password wrong combination")
	ErrNotFound           = errors.New("not found")
	ErrValidation         = errors.New("validation error")
	ErrInvalidAccessToken = errors.New("invalid access token")
)
