package usecase

import (
	"context"
	"strings"

	"github.com/harmeetsingh/grpc-envoy/services/auth/internal/domain"
	"github.com/harmeetsingh/grpc-envoy/services/auth/internal/port"
)

type AuthUseCase struct {
	store  port.UserStore
	signer port.TokenSigner
}

func New(store port.UserStore, signer port.TokenSigner) *AuthUseCase {
	return &AuthUseCase{store: store, signer: signer}
}

func (uc *AuthUseCase) Login(ctx context.Context, creds domain.Credentials, host string) (*domain.TokenResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if host == "" {
		return nil, domain.ErrEmptyHost
	}

	userID, err := uc.store.Authenticate(ctx, creds.Username, creds.Password)
	if err != nil {
		return nil, err
	}

	tenant := extractSubdomain(host)
	token, err := uc.signer.Sign(userID, tenant)
	if err != nil {
		return nil, err
	}

	return &domain.TokenResult{Token: token}, nil
}

func extractSubdomain(host string) string {
	// Strip port if present
	h := host
	if idx := strings.LastIndex(h, ":"); idx != -1 {
		h = h[:idx]
	}
	// First label is the subdomain
	parts := strings.SplitN(h, ".", 2)
	return parts[0]
}
