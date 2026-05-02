package memory_test

import (
	"context"
	"testing"

	"github.com/harmeetsingh/grpc-envoy/services/user/internal/adapter/memory"
	"github.com/harmeetsingh/grpc-envoy/services/user/internal/domain"
)

func TestSaveAndFindByID(t *testing.T) {
	repo := memory.New()
	user := &domain.User{ID: "1", Name: "Alice"}

	if err := repo.Save(context.Background(), user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := repo.FindByID(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.Name != "Alice" {
		t.Fatalf("expected Alice, got %s", found.Name)
	}
}

func TestFindByID_NotFound(t *testing.T) {
	repo := memory.New()
	_, err := repo.FindByID(context.Background(), "missing")
	if err != domain.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
