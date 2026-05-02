package domain

import "errors"

var ErrEmptyName = errors.New("name must not be empty")

type Greeting struct {
	Message string
}
