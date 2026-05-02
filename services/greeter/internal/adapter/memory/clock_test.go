package memory_test

import (
	"testing"
	"time"

	"github.com/harmeetsingh/grpc-envoy/services/greeter/internal/adapter/memory"
)

func TestRealClock_ReturnsCurrentTime(t *testing.T) {
	clock := memory.NewRealClock()
	before := time.Now()
	got := clock.Now()
	after := time.Now()

	if got.Before(before) || got.After(after) {
		t.Fatalf("expected time between %v and %v, got %v", before, after, got)
	}
}
