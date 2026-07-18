package rate_limiter

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestTokenBucketAllow(t *testing.T) {
	tb := NewTokenBucket(5, 10)
	for i := 0; i < 5; i++ {
		if !tb.Allow() {
			t.Errorf("iteration %d: expected Allow() to be true", i)
		}
	}
	if tb.Allow() {
		t.Error("expected Allow() to be false when bucket is empty")
	}
}

func TestTokenBucketRefill(t *testing.T) {
	tb := NewTokenBucket(5, 100)
	for i := 0; i < 5; i++ {
		tb.Allow()
	}
	time.Sleep(60 * time.Millisecond)
	if !tb.Allow() {
		t.Error("expected Allow() to be true after refill")
	}
}

func TestTokenBucketAllowN(t *testing.T) {
	t.Parallel()
	tb := NewTokenBucket(10, 10)
	if !tb.AllowN(5) {
		t.Error("expected AllowN(5) to be true")
	}
	if tb.AllowN(7) {
		t.Error("expected AllowN(7) to be false (only 5 remain)")
	}
}

func TestTokenBucketExactCapacity(t *testing.T) {
	t.Parallel()
	tb := NewTokenBucket(3, 5)
	if !tb.AllowN(3) {
		t.Error("expected AllowN(3) to be true")
	}
	if tb.Allow() {
		t.Error("expected Allow() to be false after consuming all tokens")
	}
}

func TestTokenBucketStress(t *testing.T) {
	tb := NewTokenBucket(10, 1000)
	var allowed atomic.Int64
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if tb.Allow() {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()
	if n := allowed.Load(); n > 10 {
		t.Errorf("stress: allowed %d tokens, expected at most 10", n)
	}
}

func TestTokenBucketConcurrentSafe(t *testing.T) {
	tb := NewTokenBucket(100, 1000)
	var wg sync.WaitGroup
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tb.Allow()
		}()
	}
	wg.Wait()
}
