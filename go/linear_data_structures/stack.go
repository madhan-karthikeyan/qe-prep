package linear_data_structures

// Stack is a generic LIFO data structure backed by a slice.
type Stack[T any] struct {
	items []T
}

// NewStack creates a new empty Stack.
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Push adds an element to the top of the stack.
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top element. Returns the zero value and false if
// the stack is empty.
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	idx := len(s.items) - 1
	item := s.items[idx]
	s.items = s.items[:idx]
	return item, true
}

// Peek returns the top element without removing it. Returns the zero value and
// false if the stack is empty.
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// IsEmpty returns true if the stack has no elements.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Size returns the number of elements in the stack.
func (s *Stack[T]) Size() int {
	return len(s.items)
}
