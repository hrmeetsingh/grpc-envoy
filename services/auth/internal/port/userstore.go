package port

import "context"

type UserStore interface {
	Authenticate(ctx context.Context, username, password string) (userID string, err error)
}
