package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmptyHost          = errors.New("host header required")
)

type Credentials struct {
	Username string
	Password string
}

type TokenResult struct {
	Token string
}
