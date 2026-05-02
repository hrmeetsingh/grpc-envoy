package userstore

import (
	"context"

	"github.com/harmeetsingh/grpc-envoy/services/auth/internal/domain"
)

var users = map[string]struct {
	password string
	id       string
}{
	"alice": {password: "password123", id: "usr-alice-001"},
	"bob":   {password: "password456", id: "usr-bob-002"},
}

type Hardcoded struct{}

func NewHardcoded() *Hardcoded {
	return &Hardcoded{}
}

func (h *Hardcoded) Authenticate(_ context.Context, username, password string) (string, error) {
	u, ok := users[username]
	if !ok || u.password != password {
		return "", domain.ErrInvalidCredentials
	}
	return u.id, nil
}
