package jwt_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	jwtadapter "github.com/harmeetsingh/grpc-envoy/services/auth/internal/adapter/jwt"
)

func mustGenerateKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	return key
}

func TestSign_ProducesValidToken(t *testing.T) {
	key := mustGenerateKey(t)
	signer := jwtadapter.NewRS256Signer(key, "https://auth.example.com")

	token, err := signer.Sign("user-123", "acme")
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestSign_TokenContainsTenantClaim(t *testing.T) {
	key := mustGenerateKey(t)
	signer := jwtadapter.NewRS256Signer(key, "https://auth.example.com")

	token, err := signer.Sign("user-123", "acme")
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	claims, err := signer.Verify(token)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	tenant, ok := claims["tenant"].(string)
	if !ok || tenant != "acme" {
		t.Errorf("expected tenant 'acme', got %v", claims["tenant"])
	}
}

func TestSign_TokenContainsSubClaim(t *testing.T) {
	key := mustGenerateKey(t)
	signer := jwtadapter.NewRS256Signer(key, "https://auth.example.com")

	token, err := signer.Sign("user-456", "globex")
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	claims, err := signer.Verify(token)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub != "user-456" {
		t.Errorf("expected sub 'user-456', got %v", claims["sub"])
	}
}

func TestSign_TokenContainsIssuer(t *testing.T) {
	key := mustGenerateKey(t)
	signer := jwtadapter.NewRS256Signer(key, "https://auth.example.com")

	token, err := signer.Sign("user-123", "acme")
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	claims, err := signer.Verify(token)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	iss, ok := claims["iss"].(string)
	if !ok || iss != "https://auth.example.com" {
		t.Errorf("expected iss 'https://auth.example.com', got %v", claims["iss"])
	}
}

func TestVerify_InvalidToken_ReturnsError(t *testing.T) {
	key := mustGenerateKey(t)
	signer := jwtadapter.NewRS256Signer(key, "https://auth.example.com")

	_, err := signer.Verify("invalid.token.here")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestVerify_WrongKey_ReturnsError(t *testing.T) {
	key1 := mustGenerateKey(t)
	key2 := mustGenerateKey(t)
	signer1 := jwtadapter.NewRS256Signer(key1, "https://auth.example.com")
	signer2 := jwtadapter.NewRS256Signer(key2, "https://auth.example.com")

	token, err := signer1.Sign("user-123", "acme")
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	_, err = signer2.Verify(token)
	if err == nil {
		t.Fatal("expected error when verifying with wrong key")
	}
}
