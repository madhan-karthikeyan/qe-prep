package lru_cache

import (
	"sync"
	"testing"
)

func TestCacheGetPut(t *testing.T) {
	t.Parallel()
	c := NewCache(3)
	c.Put(1, 10)
	c.Put(2, 20)
	c.Put(3, 30)

	v, ok := c.Get(1)
	if !ok || v != 10 {
		t.Errorf("Get(1) = (%d, %v), want (10, true)", v, ok)
	}
	v, ok = c.Get(2)
	if !ok || v != 20 {
		t.Errorf("Get(2) = (%d, %v), want (20, true)", v, ok)
	}
	v, ok = c.Get(4)
	if ok {
		t.Errorf("Get(4) = (%d, %v), want (0, false)", v, ok)
	}
}

func TestCacheEviction(t *testing.T) {
	t.Parallel()
	c := NewCache(2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Put(3, 3) // evicts key 1

	if _, ok := c.Get(1); ok {
		t.Error("expected key 1 to be evicted")
	}
	if v, ok := c.Get(2); !ok || v != 2 {
		t.Errorf("Get(2) = (%d, %v)", v, ok)
	}
	if v, ok := c.Get(3); !ok || v != 3 {
		t.Errorf("Get(3) = (%d, %v)", v, ok)
	}
}

func TestCacheUpdate(t *testing.T) {
	t.Parallel()
	c := NewCache(2)
	c.Put(1, 10)
	c.Put(2, 20)
	c.Put(1, 100) // update, not eviction

	c.Put(3, 30) // evicts key 2

	if v, ok := c.Get(1); !ok || v != 100 {
		t.Errorf("Get(1) = (%d, %v), want (100, true)", v, ok)
	}
	if _, ok := c.Get(2); ok {
		t.Error("expected key 2 to be evicted")
	}
}

func TestCacheEvictOrder(t *testing.T) {
	t.Parallel()
	c := NewCache(3)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Put(3, 3)
	c.Get(1) // makes 1 most recent
	c.Put(4, 4) // evicts 2

	if _, ok := c.Get(2); ok {
		t.Error("expected key 2 to be evicted (least recently used)")
	}
}

func TestCacheZeroCapacity(t *testing.T) {
	c := NewCache(0)
	c.Put(1, 1)
	if _, ok := c.Get(1); ok {
		t.Error("expected no entry with zero capacity")
	}
}

func TestCacheConcurrent(t *testing.T) {
	c := NewCache(100)
	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			c.Put(k, k*10)
		}(i)
	}
	wg.Wait()
	for i := range 50 {
		v, ok := c.Get(i)
		if !ok {
			t.Errorf("key %d missing", i)
		} else if v != i*10 {
			t.Errorf("key %d: got %d, want %d", i, v, i*10)
		}
	}
}
