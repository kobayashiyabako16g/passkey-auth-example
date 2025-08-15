package auth

import "errors"

var (
	ErrUserExists         = errors.New("user already exists")
	ErrSessionNotFound    = errors.New("session not found")
	ErrFinishRegistration = errors.New("registration failed")
)
