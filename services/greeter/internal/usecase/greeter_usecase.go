package usecase

import (
	"context"
	"fmt"

	"github.com/harmeetsingh/grpc-envoy/services/greeter/internal/domain"
	"github.com/harmeetsingh/grpc-envoy/services/greeter/internal/port"
)

type GreeterUseCase struct {
	clock port.Clock
}

func New(clock port.Clock) *GreeterUseCase {
	return &GreeterUseCase{clock: clock}
}

func (uc *GreeterUseCase) SayHello(ctx context.Context, name string) (*domain.Greeting, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if name == "" {
		return nil, domain.ErrEmptyName
	}
	t := uc.clock.Now().UTC().Format("2006-01-02T15:04:05Z")
	return &domain.Greeting{
		Message: fmt.Sprintf("Hello, %s! The time is %s.", name, t),
	}, nil
}

func (uc *GreeterUseCase) SayGoodbye(ctx context.Context, name string) (*domain.Greeting, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if name == "" {
		return nil, domain.ErrEmptyName
	}
	t := uc.clock.Now().UTC().Format("2006-01-02T15:04:05Z")
	return &domain.Greeting{
		Message: fmt.Sprintf("Goodbye, %s! The time is %s.", name, t),
	}, nil
}
