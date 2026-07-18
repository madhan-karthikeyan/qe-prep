package producer_consumer

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestBlockingQueuePutGet(t *testing.T) {
	t.Parallel()
	q := NewBlockingQueue[int](3)
	ctx := context.Background()

	if err := q.Put(ctx, 1); err != nil {
		t.Fatal(err)
	}
	if err := q.Put(ctx, 2); err != nil {
		t.Fatal(err)
	}
	v, err := q.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if v != 1 {
		t.Errorf("got %d, want 1", v)
	}
	v, err = q.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if v != 2 {
		t.Errorf("got %d, want 2", v)
	}
}

func TestBlockingQueueBlockingBehavior(t *testing.T) {
	q := NewBlockingQueue[int](1)
	ctx := context.Background()
	_ = q.Put(ctx, 1) // fills the queue

	done := make(chan struct{})
	go func() {
		// this Put should block until we dequeue
		if err := q.Put(ctx, 2); err != nil {
			t.Errorf("Put: %v", err)
		}
		close(done)
	}()

	select {
	case <-done:
		t.Error("Put should have blocked")
	case <-time.After(50 * time.Millisecond):
		// expected: Put blocked
	}

	// dequeue to unblock the putter
	v, _ := q.Get(ctx)
	if v != 1 {
		t.Errorf("got %d, want 1", v)
	}

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Error("Put did not unblock after dequeue")
	}
}

func TestBlockingQueueTimeout(t *testing.T) {
	q := NewBlockingQueue[int](1)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_ = q.Put(context.Background(), 1)
	err := q.Put(ctx, 2)
	if err != context.DeadlineExceeded && err != context.Canceled {
		t.Errorf("expected deadline exceeded, got %v", err)
	}
}

func TestBlockingQueueClose(t *testing.T) {
	q := NewBlockingQueue[int](3)
	ctx := context.Background()
	q.Put(ctx, 1)
	q.Put(ctx, 2)
	q.Close()

	// Get should still work for remaining items
	v, err := q.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if v != 1 {
		t.Errorf("got %d, want 1", v)
	}

	v, err = q.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if v != 2 {
		t.Errorf("got %d, want 2", v)
	}

	// after draining, Get should return ErrQueueClosed
	_, err = q.Get(ctx)
	if err != ErrQueueClosed {
		t.Errorf("expected ErrQueueClosed, got %v", err)
	}

	// Put should also return ErrQueueClosed
	err = q.Put(ctx, 3)
	if err != ErrQueueClosed {
		t.Errorf("expected ErrQueueClosed, got %v", err)
	}
}

func TestBlockingQueueMultipleGoroutines(t *testing.T) {
	q := NewBlockingQueue[int](10)
	ctx := context.Background()
	n := 100

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range n {
			if err := q.Put(ctx, i); err != nil {
				return
			}
		}
	}()

	received := make([]bool, n)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range n {
			v, err := q.Get(ctx)
			if err != nil {
				return
			}
			received[v] = true
		}
	}()

	wg.Wait()
	for i, ok := range received {
		if !ok {
			t.Errorf("item %d not received", i)
		}
	}
}

func TestBlockingQueueStress(t *testing.T) {
	q := NewBlockingQueue[int](100)
	ctx := context.Background()
	const numProducers = 10
	const numConsumers = 10
	const itemsPerProducer = 1000

	var wg sync.WaitGroup
	var produced sync.WaitGroup
	produced.Add(numProducers)
	for i := range numProducers {
		go func(base int) {
			defer produced.Done()
			for j := range itemsPerProducer {
				q.Put(ctx, base*itemsPerProducer+j)
			}
		}(i)
	}

	var mu sync.Mutex
	received := make(map[int]int)
	wg.Add(numConsumers)
	for range numConsumers {
		go func() {
			defer wg.Done()
			for {
				v, err := q.Get(ctx)
				if err != nil {
					return
				}
				mu.Lock()
				received[v]++
				mu.Unlock()
			}
		}()
	}

	produced.Wait()
	q.Close()
	wg.Wait()

	total := numProducers * itemsPerProducer
	if len(received) != total {
		t.Errorf("received %d unique items, want %d", len(received), total)
	}
	for k, v := range received {
		if v != 1 {
			t.Errorf("item %d received %d times, want 1", k, v)
		}
	}
}
