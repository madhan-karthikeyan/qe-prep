package linear_data_structures

import "testing"

func TestCircularQueue(t *testing.T) {
	t.Parallel()
	q := NewCircularQueue[int](3)
	if !q.IsEmpty() {
		t.Error("new queue should be empty")
	}
	if q.IsFull() {
		t.Error("new queue should not be full")
	}
	if ok := q.Enqueue(1); !ok {
		t.Error("Enqueue(1) should succeed")
	}
	if ok := q.Enqueue(2); !ok {
		t.Error("Enqueue(2) should succeed")
	}
	if ok := q.Enqueue(3); !ok {
		t.Error("Enqueue(3) should succeed")
	}
	if !q.IsFull() {
		t.Error("queue should be full")
	}
	if ok := q.Enqueue(4); ok {
		t.Error("Enqueue(4) should fail (full)")
	}
	v, ok := q.Dequeue()
	if !ok || v != 1 {
		t.Errorf("Dequeue = (%d, %v), want (1, true)", v, ok)
	}
	if ok := q.Enqueue(4); !ok {
		t.Error("Enqueue(4) after dequeue should succeed")
	}
	v, ok = q.Dequeue()
	if !ok || v != 2 {
		t.Errorf("Dequeue = (%d, %v), want (2, true)", v, ok)
	}
	v, ok = q.Dequeue()
	if !ok || v != 3 {
		t.Errorf("Dequeue = (%d, %v), want (3, true)", v, ok)
	}
	v, ok = q.Dequeue()
	if !ok || v != 4 {
		t.Errorf("Dequeue = (%d, %v), want (4, true)", v, ok)
	}
	if !q.IsEmpty() {
		t.Error("queue should be empty")
	}
}

func TestCircularQueueWrapAround(t *testing.T) {
	q := NewCircularQueue[int](3)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	q.Dequeue() // 1 out
	q.Enqueue(4)
	q.Dequeue() // 2 out
	q.Enqueue(5)
	q.Dequeue() // 3 out
	q.Enqueue(6)
	// buffer should now be [4,5,6] in order
	v, ok := q.Dequeue()
	if !ok || v != 4 {
		t.Errorf("expected 4, got (%d, %v)", v, ok)
	}
	v, ok = q.Dequeue()
	if !ok || v != 5 {
		t.Errorf("expected 5, got (%d, %v)", v, ok)
	}
	v, ok = q.Dequeue()
	if !ok || v != 6 {
		t.Errorf("expected 6, got (%d, %v)", v, ok)
	}
}

func TestCircularQueuePeekEmpty(t *testing.T) {
	q := NewCircularQueue[int](1)
	v, ok := q.Peek()
	if ok {
		t.Errorf("expected empty, got (%d, %v)", v, ok)
	}
}

func TestCircularQueueDequeueEmpty(t *testing.T) {
	q := NewCircularQueue[int](1)
	v, ok := q.Dequeue()
	if ok {
		t.Errorf("expected empty, got (%d, %v)", v, ok)
	}
}

func TestCircularQueueSize(t *testing.T) {
	q := NewCircularQueue[int](5)
	if s := q.Size(); s != 0 {
		t.Errorf("Size = %d, want 0", s)
	}
	q.Enqueue(1)
	q.Enqueue(2)
	if s := q.Size(); s != 2 {
		t.Errorf("Size = %d, want 2", s)
	}
	q.Dequeue()
	if s := q.Size(); s != 1 {
		t.Errorf("Size = %d, want 1", s)
	}
}
