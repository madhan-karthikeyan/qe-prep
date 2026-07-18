package thread_pool

import (
	"sync/atomic"
	"testing"
)

func TestPoolSubmitTasks(t *testing.T) {
	p := New(2, 10)
	var counter atomic.Int32

	for range 10 {
		p.Submit(func() {
			counter.Add(1)
		})
	}
	p.Stop()
	p.Wait()

	if n := counter.Load(); n != 10 {
		t.Errorf("expected 10 tasks completed, got %d", n)
	}
}

func TestPoolResultsViaCallback(t *testing.T) {
	p := New(2, 10)
	results := make([]int, 5)

	for i := range 5 {
		i := i // capture
		p.Submit(func() {
			results[i] = i * 2
		})
	}
	p.Stop()
	p.Wait()

	for i, v := range results {
		if v != i*2 {
			t.Errorf("results[%d] = %d, want %d", i, v, i*2)
		}
	}
}

func TestPoolGracefulShutdown(t *testing.T) {
	p := New(2, 5)
	var counter atomic.Int32

	for range 20 {
		p.Submit(func() {
			counter.Add(1)
		})
	}
	// Stop while tasks may still be in the queue
	p.Stop()
	p.Wait()

	// After stop, Submit should return false
	if ok := p.Submit(func() {}); ok {
		t.Error("Submit after Stop should return false")
	}
}

func TestPoolNoTasks(t *testing.T) {
	p := New(2, 10)
	p.Stop()
	p.Wait()
	// should not hang or panic
}

func TestPoolSubmitAndWait(t *testing.T) {
	p := New(2, 10)
	var result int
	ok := p.SubmitAndWait(func() {
		result = 42
	})
	if !ok {
		t.Fatal("SubmitAndWait returned false")
	}
	if result != 42 {
		t.Errorf("result = %d, want 42", result)
	}
	p.Stop()
	p.Wait()
}

func TestPoolStress(t *testing.T) {
	const numTasks = 1000
	const numWorkers = 8

	p := New(numWorkers, numTasks)
	var counter atomic.Int64

	for range numTasks {
		p.Submit(func() {
			counter.Add(1)
		})
	}
	p.Stop()
	p.Wait()

	if n := counter.Load(); n != numTasks {
		t.Errorf("expected %d tasks, got %d", numTasks, n)
	}
}
