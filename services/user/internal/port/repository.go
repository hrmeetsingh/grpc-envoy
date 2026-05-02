package port

import (
	"context"

	"github.com/harmeetsingh/grpc-envoy/services/user/internal/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
}
