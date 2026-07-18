package linear_data_structures

import "testing"

func TestStack(t *testing.T) {
	t.Parallel()
	s := NewStack[int]()
	if !s.IsEmpty() {
		t.Error("new stack should be empty")
	}
	if size := s.Size(); size != 0 {
		t.Errorf("Size = %d, want 0", size)
	}
	s.Push(1)
	s.Push(2)
	s.Push(3)
	if s.Size() != 3 {
		t.Errorf("Size = %d, want 3", s.Size())
	}
	v, ok := s.Peek()
	if !ok || v != 3 {
		t.Errorf("Peek = (%d, %v), want (3, true)", v, ok)
	}
	v, ok = s.Pop()
	if !ok || v != 3 {
		t.Errorf("Pop = (%d, %v), want (3, true)", v, ok)
	}
	v, ok = s.Pop()
	if !ok || v != 2 {
		t.Errorf("Pop = (%d, %v), want (2, true)", v, ok)
	}
	v, ok = s.Pop()
	if !ok || v != 1 {
		t.Errorf("Pop = (%d, %v), want (1, true)", v, ok)
	}
	if !s.IsEmpty() {
		t.Error("stack should be empty after popping all")
	}
}

func TestStackPopEmpty(t *testing.T) {
	t.Parallel()
	s := NewStack[string]()
	v, ok := s.Pop()
	if ok {
		t.Errorf("expected empty, got (%q, %v)", v, ok)
	}
}

func TestStackPeekEmpty(t *testing.T) {
	t.Parallel()
	s := NewStack[float64]()
	v, ok := s.Peek()
	if ok {
		t.Errorf("expected empty, got (%f, %v)", v, ok)
	}
}

func TestStackString(t *testing.T) {
	t.Parallel()
	s := NewStack[string]()
	s.Push("a")
	s.Push("b")
	v, ok := s.Pop()
	if !ok || v != "b" {
		t.Errorf("Pop = (%q, %v), want (\"b\", true)", v, ok)
	}
}
