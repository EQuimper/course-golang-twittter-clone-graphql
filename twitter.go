package twitter

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrValidation = errors.New("validation error")
)
