package auth

import "errors"

var (
	ErrUsernameExists    = errors.New("username already exists")
	ErrEmailExists       = errors.New("email already exists")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrInvalidRole       = errors.New("invalid role")
	ErrUserNotFound      = errors.New("user not found")
	ErrHashingPassword   = errors.New("password hashing error")
)
