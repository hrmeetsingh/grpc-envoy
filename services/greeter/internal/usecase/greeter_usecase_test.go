package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/harmeetsingh/grpc-envoy/services/greeter/internal/domain"
	"github.com/harmeetsingh/grpc-envoy/services/greeter/internal/usecase"
)

type stubClock struct {
	fixed time.Time
}

func (c *stubClock) Now() time.Time { return c.fixed }

func TestSayHello_Success(t *testing.T) {
	clock := &stubClock{fixed: time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC)}
	uc := usecase.New(clock)
	g, err := uc.SayHello(context.Background(), "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "Hello, Alice! The time is 2026-05-01T10:00:00Z."
	if g.Message != expected {
		t.Fatalf("expected %q, got %q", expected, g.Message)
	}
}

func TestSayHello_EmptyName(t *testing.T) {
	uc := usecase.New(&stubClock{})
	_, err := uc.SayHello(context.Background(), "")
	if err != domain.ErrEmptyName {
		t.Fatalf("expected ErrEmptyName, got %v", err)
	}
}

func TestSayGoodbye_Success(t *testing.T) {
	clock := &stubClock{fixed: time.Date(2026, 5, 1, 18, 0, 0, 0, time.UTC)}
	uc := usecase.New(clock)
	g, err := uc.SayGoodbye(context.Background(), "Bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "Goodbye, Bob! The time is 2026-05-01T18:00:00Z."
	if g.Message != expected {
		t.Fatalf("expected %q, got %q", expected, g.Message)
	}
}

func TestSayGoodbye_EmptyName(t *testing.T) {
	uc := usecase.New(&stubClock{})
	_, err := uc.SayGoodbye(context.Background(), "")
	if err != domain.ErrEmptyName {
		t.Fatalf("expected ErrEmptyName, got %v", err)
	}
}
