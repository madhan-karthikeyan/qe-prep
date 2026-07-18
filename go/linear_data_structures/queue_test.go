package linear_data_structures

import "testing"

func TestQueue(t *testing.T) {
	t.Parallel()
	q := NewQueue[int]()
	if !q.IsEmpty() {
		t.Error("new queue should be empty")
	}
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	if q.Size() != 3 {
		t.Errorf("Size = %d, want 3", q.Size())
	}
	v, ok := q.Peek()
	if !ok || v != 1 {
		t.Errorf("Peek = (%d, %v), want (1, true)", v, ok)
	}
	v, ok = q.Dequeue()
	if !ok || v != 1 {
		t.Errorf("Dequeue = (%d, %v), want (1, true)", v, ok)
	}
	v, ok = q.Dequeue()
	if !ok || v != 2 {
		t.Errorf("Dequeue = (%d, %v), want (2, true)", v, ok)
	}
	q.Enqueue(4)
	v, ok = q.Dequeue()
	if !ok || v != 3 {
		t.Errorf("Dequeue = (%d, %v), want (3, true)", v, ok)
	}
	v, ok = q.Dequeue()
	if !ok || v != 4 {
		t.Errorf("Dequeue = (%d, %v), want (4, true)", v, ok)
	}
	if !q.IsEmpty() {
		t.Error("queue should be empty after dequeuing all")
	}
}

func TestQueueDequeueEmpty(t *testing.T) {
	t.Parallel()
	q := NewQueue[string]()
	v, ok := q.Dequeue()
	if ok {
		t.Errorf("expected empty, got (%q, %v)", v, ok)
	}
}

func TestQueuePeekEmpty(t *testing.T) {
	t.Parallel()
	q := NewQueue[float64]()
	v, ok := q.Peek()
	if ok {
		t.Errorf("expected empty, got (%f, %v)", v, ok)
	}
}
