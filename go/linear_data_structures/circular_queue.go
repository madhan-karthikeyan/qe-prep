package linear_data_structures

// CircularQueue is a fixed-size FIFO buffer that wraps around.
type CircularQueue[T any] struct {
	buf  []T
	head int
	tail int
	size int
	cap  int
}

// NewCircularQueue creates a new CircularQueue with the given capacity.
func NewCircularQueue[T any](capacity int) *CircularQueue[T] {
	if capacity < 1 {
		capacity = 1
	}
	return &CircularQueue[T]{
		buf: make([]T, capacity),
		cap: capacity,
	}
}

// Enqueue adds an element to the back. Returns false if the queue is full.
func (q *CircularQueue[T]) Enqueue(item T) bool {
	if q.size == q.cap {
		return false
	}
	q.buf[q.tail] = item
	q.tail = (q.tail + 1) % q.cap
	q.size++
	return true
}

// Dequeue removes and returns the front element. Returns the zero value and
// false if the queue is empty.
func (q *CircularQueue[T]) Dequeue() (T, bool) {
	if q.size == 0 {
		var zero T
		return zero, false
	}
	item := q.buf[q.head]
	q.head = (q.head + 1) % q.cap
	q.size--
	return item, true
}

// Peek returns the front element without removing it. Returns the zero value
// and false if the queue is empty.
func (q *CircularQueue[T]) Peek() (T, bool) {
	if q.size == 0 {
		var zero T
		return zero, false
	}
	return q.buf[q.head], true
}

// IsFull returns true if the queue has reached its capacity.
func (q *CircularQueue[T]) IsFull() bool {
	return q.size == q.cap
}

// IsEmpty returns true if the queue has no elements.
func (q *CircularQueue[T]) IsEmpty() bool {
	return q.size == 0
}

// Size returns the number of elements in the queue.
func (q *CircularQueue[T]) Size() int {
	return q.size
}
