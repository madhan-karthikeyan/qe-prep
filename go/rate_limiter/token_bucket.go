package rate_limiter

import (
	"sync"
	"time"
)

// TokenBucket implements a token bucket rate limiter with burst capacity and
// continuous refill. Thread-safe via sync.Mutex.
type TokenBucket struct {
	mu         sync.Mutex
	capacity   float64
	tokens     float64
	rate       float64
	lastRefill time.Time
}

// NewTokenBucket creates a new TokenBucket with the given burst capacity and
// refill rate (tokens per second).
func NewTokenBucket(capacity int, rate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   float64(capacity),
		tokens:     float64(capacity),
		rate:       rate,
		lastRefill: time.Now(),
	}
}

// refill adds tokens based on elapsed time since last refill.
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	if elapsed > 0 {
		tb.tokens += elapsed * tb.rate
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastRefill = now
	}
}

// Allow checks if a single token can be consumed. Returns true if within rate
// limit.
func (tb *TokenBucket) Allow() bool {
	return tb.AllowN(1)
}

// AllowN checks if n tokens can be consumed atomically. Returns true if within
// rate limit.
func (tb *TokenBucket) AllowN(n int) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.refill()
	if tb.tokens >= float64(n) {
		tb.tokens -= float64(n)
		return true
	}
	return false
}
