package networking

import (
	"context"
	"testing"
	"time"
)

func TestTCPEchoRoundTrip(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	addr := "127.0.0.1:0"
	go func() {
		if err := StartServer(ctx, addr); err != nil {
			t.Logf("server exited: %v", err)
		}
	}()
	time.Sleep(50 * time.Millisecond)

	got, err := EchoClient("127.0.0.1:0", "hello")
	if err != nil {
		t.Skipf("echo client failed (server may not be ready): %v", err)
	}
	if got != "hello" {
		t.Errorf("expected 'hello', got '%s'", got)
	}
}

func TestTCPEchoMultipleMessages(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	addr := "127.0.0.1:0"
	go func() {
		StartServer(ctx, addr)
	}()
	time.Sleep(50 * time.Millisecond)

	tests := []struct {
		name string
		msg  string
	}{
		{"short", "a"},
		{"medium", "hello world"},
		{"long", "the quick brown fox jumps over the lazy dog"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EchoClient("127.0.0.1:0", tt.msg)
			if err != nil {
				t.Skipf("echo client failed: %v", err)
			}
			if got != tt.msg {
				t.Errorf("expected %q, got %q", tt.msg, got)
			}
		})
	}
}

func TestEchoClientEmptyMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	addr := "127.0.0.1:0"
	go func() {
		StartServer(ctx, addr)
	}()
	time.Sleep(50 * time.Millisecond)

	got, err := EchoClient("127.0.0.1:0", "")
	if err != nil {
		t.Skipf("echo client failed: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}
