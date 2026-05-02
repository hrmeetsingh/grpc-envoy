package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/harmeetsingh/grpc-envoy/services/user/internal/domain"
	"github.com/harmeetsingh/grpc-envoy/services/user/internal/port"
)

type UserUseCase struct {
	repo port.UserRepository
}

func New(repo port.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, name string) (*domain.User, error) {
	if name == "" {
		return nil, domain.ErrEmptyName
	}
	u := &domain.User{
		ID:   newID(),
		Name: name,
	}
	if err := uc.repo.Save(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (uc *UserUseCase) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return uc.repo.FindByID(ctx, id)
}

func newID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
