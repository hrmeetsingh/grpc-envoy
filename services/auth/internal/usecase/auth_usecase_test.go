package usecase_test

import (
	"context"
	"testing"

	"github.com/harmeetsingh/grpc-envoy/services/auth/internal/domain"
	"github.com/harmeetsingh/grpc-envoy/services/auth/internal/usecase"
)

type stubStore struct {
	users map[string]string // username → password
}

func (s *stubStore) Authenticate(_ context.Context, username, password string) (string, error) {
	if pw, ok := s.users[username]; ok && pw == password {
		return "user-" + username, nil
	}
	return "", domain.ErrInvalidCredentials
}

type stubSigner struct {
	lastUserID string
	lastTenant string
}

func (s *stubSigner) Sign(userID string, tenant string) (string, error) {
	s.lastUserID = userID
	s.lastTenant = tenant
	return "token." + userID + "." + tenant, nil
}

func (s *stubSigner) Verify(tokenString string) (map[string]interface{}, error) {
	return nil, nil
}

func TestLogin_ValidCredentials_ReturnsToken(t *testing.T) {
	store := &stubStore{users: map[string]string{"alice": "secret"}}
	signer := &stubSigner{}
	uc := usecase.New(store, signer)

	result, err := uc.Login(context.Background(), domain.Credentials{
		Username: "alice",
		Password: "secret",
	}, "acme.example.com")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil || result.Token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestLogin_ValidCredentials_EmbedsTenantFromHost(t *testing.T) {
	store := &stubStore{users: map[string]string{"alice": "secret"}}
	signer := &stubSigner{}
	uc := usecase.New(store, signer)

	_, err := uc.Login(context.Background(), domain.Credentials{
		Username: "alice",
		Password: "secret",
	}, "acme.example.com")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if signer.lastTenant != "acme" {
		t.Errorf("expected tenant 'acme', got %q", signer.lastTenant)
	}
}

func TestLogin_ValidCredentials_PassesUserIDToSigner(t *testing.T) {
	store := &stubStore{users: map[string]string{"bob": "pass"}}
	signer := &stubSigner{}
	uc := usecase.New(store, signer)

	_, err := uc.Login(context.Background(), domain.Credentials{
		Username: "bob",
		Password: "pass",
	}, "globex.example.com")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if signer.lastUserID != "user-bob" {
		t.Errorf("expected userID 'user-bob', got %q", signer.lastUserID)
	}
}

func TestLogin_InvalidCredentials_ReturnsError(t *testing.T) {
	store := &stubStore{users: map[string]string{"alice": "secret"}}
	signer := &stubSigner{}
	uc := usecase.New(store, signer)

	_, err := uc.Login(context.Background(), domain.Credentials{
		Username: "alice",
		Password: "wrong",
	}, "acme.example.com")

	if err != domain.ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_EmptyHost_ReturnsError(t *testing.T) {
	store := &stubStore{users: map[string]string{"alice": "secret"}}
	signer := &stubSigner{}
	uc := usecase.New(store, signer)

	_, err := uc.Login(context.Background(), domain.Credentials{
		Username: "alice",
		Password: "secret",
	}, "")

	if err != domain.ErrEmptyHost {
		t.Errorf("expected ErrEmptyHost, got %v", err)
	}
}

func TestLogin_HostWithPort_ExtractsSubdomain(t *testing.T) {
	store := &stubStore{users: map[string]string{"alice": "secret"}}
	signer := &stubSigner{}
	uc := usecase.New(store, signer)

	_, err := uc.Login(context.Background(), domain.Credentials{
		Username: "alice",
		Password: "secret",
	}, "acme.example.com:8080")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if signer.lastTenant != "acme" {
		t.Errorf("expected tenant 'acme', got %q", signer.lastTenant)
	}
}
