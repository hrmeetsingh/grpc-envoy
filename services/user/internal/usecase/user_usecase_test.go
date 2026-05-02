package usecase_test

import (
	"context"
	"testing"

	"github.com/harmeetsingh/grpc-envoy/services/user/internal/domain"
	"github.com/harmeetsingh/grpc-envoy/services/user/internal/usecase"
)

// stubRepo implements port.UserRepository for unit tests.
type stubRepo struct {
	users map[string]*domain.User
}

func newStubRepo() *stubRepo {
	return &stubRepo{users: make(map[string]*domain.User)}
}

func (r *stubRepo) Save(_ context.Context, u *domain.User) error {
	r.users[u.ID] = u
	return nil
}

func (r *stubRepo) FindByID(_ context.Context, id string) (*domain.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return u, nil
}

func TestCreateUser_Success(t *testing.T) {
	uc := usecase.New(newStubRepo())
	user, err := uc.CreateUser(context.Background(), "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Name != "Alice" {
		t.Fatalf("expected name Alice, got %s", user.Name)
	}
	if user.ID == "" {
		t.Fatal("expected non-empty ID")
	}
}

func TestCreateUser_EmptyName(t *testing.T) {
	uc := usecase.New(newStubRepo())
	_, err := uc.CreateUser(context.Background(), "")
	if err != domain.ErrEmptyName {
		t.Fatalf("expected ErrEmptyName, got %v", err)
	}
}

func TestGetUser_Found(t *testing.T) {
	repo := newStubRepo()
	uc := usecase.New(repo)
	created, _ := uc.CreateUser(context.Background(), "Bob")
	found, err := uc.GetUser(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.Name != "Bob" {
		t.Fatalf("expected Bob, got %s", found.Name)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	uc := usecase.New(newStubRepo())
	_, err := uc.GetUser(context.Background(), "nonexistent")
	if err != domain.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
