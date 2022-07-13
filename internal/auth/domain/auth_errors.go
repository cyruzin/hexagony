package domain

import "errors"

var (
	ErrAuth         = errors.New("authentication failed")
	ErrEmptyClaim   = errors.New("claim is empty")
	ErrSign         = errors.New("failed to sign the key")
	ErrUserNotFound = errors.New("user not found")
	ErrPassword     = errors.New("wrong password")
)
