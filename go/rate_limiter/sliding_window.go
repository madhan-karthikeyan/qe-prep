package rate_limiter

import (
	"sync"
	"time"
)

// SlidingWindow implements a sliding window log rate limiter. It tracks
// request timestamps and evicts those outside the window on each Allow call.
type SlidingWindow struct {
	mu       sync.Mutex
	max      int
	window   time.Duration
	requests []time.Time
}

// NewSlidingWindow creates a new sliding window rate limiter.
func NewSlidingWindow(max int, window time.Duration) *SlidingWindow {
	return &SlidingWindow{
		max:      max,
		window:   window,
		requests: make([]time.Time, 0, max),
	}
}

// Allow checks whether a new request falls within the rate limit. Returns true
// if the request is allowed.
func (sw *SlidingWindow) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	now := time.Now()
	cutoff := now.Add(-sw.window)
	var active []time.Time
	for _, t := range sw.requests {
		if t.After(cutoff) {
			active = append(active, t)
		}
	}
	sw.requests = active
	if len(sw.requests) < sw.max {
		sw.requests = append(sw.requests, now)
		return true
	}
	return false
}
