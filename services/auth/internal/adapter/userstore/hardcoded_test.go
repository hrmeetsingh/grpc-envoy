package userstore_test

import (
	"context"
	"testing"

	"github.com/harmeetsingh/grpc-envoy/services/auth/internal/adapter/userstore"
	"github.com/harmeetsingh/grpc-envoy/services/auth/internal/domain"
)

func TestHardcoded_ValidUser_ReturnsID(t *testing.T) {
	store := userstore.NewHardcoded()

	id, err := store.Authenticate(context.Background(), "alice", "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty user ID")
	}
}

func TestHardcoded_InvalidPassword_ReturnsError(t *testing.T) {
	store := userstore.NewHardcoded()

	_, err := store.Authenticate(context.Background(), "alice", "wrongpass")
	if err != domain.ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestHardcoded_UnknownUser_ReturnsError(t *testing.T) {
	store := userstore.NewHardcoded()

	_, err := store.Authenticate(context.Background(), "unknown", "password123")
	if err != domain.ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}
