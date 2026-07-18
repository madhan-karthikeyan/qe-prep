package linear_data_structures

import "container/list"

// Queue is a generic FIFO data structure backed by a linked list.
type Queue[T any] struct {
	items *list.List
}

// NewQueue creates a new empty Queue.
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{items: list.New()}
}

// Enqueue adds an element to the back of the queue.
func (q *Queue[T]) Enqueue(item T) {
	q.items.PushBack(item)
}

// Dequeue removes and returns the front element. Returns the zero value and
// false if the queue is empty.
func (q *Queue[T]) Dequeue() (T, bool) {
	front := q.items.Front()
	if front == nil {
		var zero T
		return zero, false
	}
	q.items.Remove(front)
	return front.Value.(T), true
}

// Peek returns the front element without removing it. Returns the zero value
// and false if the queue is empty.
func (q *Queue[T]) Peek() (T, bool) {
	front := q.items.Front()
	if front == nil {
		var zero T
		return zero, false
	}
	return front.Value.(T), true
}

// IsEmpty returns true if the queue has no elements.
func (q *Queue[T]) IsEmpty() bool {
	return q.items.Len() == 0
}

// Size returns the number of elements in the queue.
func (q *Queue[T]) Size() int {
	return q.items.Len()
}
