package memory

import (
	"context"
	"sync"

	"github.com/harmeetsingh/grpc-envoy/services/user/internal/domain"
)

type UserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func New() *UserRepository {
	return &UserRepository{users: make(map[string]*domain.User)}
}

func (r *UserRepository) Save(_ context.Context, u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.ID] = u
	return nil
}

func (r *UserRepository) FindByID(_ context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return u, nil
}
