package rate_limiter

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSlidingWindowUnderLimit(t *testing.T) {
	sw := NewSlidingWindow(5, time.Second)
	for i := 0; i < 5; i++ {
		if !sw.Allow() {
			t.Errorf("expected Allow() to be true, iteration %d", i)
		}
	}
}

func TestSlidingWindowOverLimit(t *testing.T) {
	sw := NewSlidingWindow(3, time.Second)
	for i := 0; i < 3; i++ {
		sw.Allow()
	}
	if sw.Allow() {
		t.Error("expected Allow() to be false when over limit")
	}
}

func TestSlidingWindowSlides(t *testing.T) {
	sw := NewSlidingWindow(2, 50*time.Millisecond)
	if !sw.Allow() {
		t.Fatal("expected Allow() to be true")
	}
	if !sw.Allow() {
		t.Fatal("expected Allow() to be true")
	}
	if sw.Allow() {
		t.Fatal("expected Allow() to be false (at limit)")
	}
	time.Sleep(60 * time.Millisecond)
	if !sw.Allow() {
		t.Error("expected Allow() to be true after window slides")
	}
}

func TestSlidingWindowEmpty(t *testing.T) {
	sw := NewSlidingWindow(1, time.Second)
	if !sw.Allow() {
		t.Error("expected Allow() to be true on empty window")
	}
}

func TestSlidingWindowStress(t *testing.T) {
	sw := NewSlidingWindow(10, time.Second)
	var allowed atomic.Int64
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if sw.Allow() {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()
	if n := allowed.Load(); n > 10 {
		t.Errorf("stress: allowed %d requests, expected at most 10", n)
	}
}
