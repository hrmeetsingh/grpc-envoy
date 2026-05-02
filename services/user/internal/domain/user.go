package domain

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmptyName    = errors.New("name must not be empty")
)

type User struct {
	ID   string
	Name string
}
