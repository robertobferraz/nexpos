package errors

import (
	"errors"
)

var ErrDuplicatedEmail = errors.New("email already been taken")
var ErrDuplicatedUsername = errors.New("username already been taken")
var ErrDuplicatedExternalUID = errors.New("external uid already been taken")
var ErrInvalidToken = errors.New("unauthorized Token")
var ErrInvalidTokenFormat = errors.New("invalid Token format")
