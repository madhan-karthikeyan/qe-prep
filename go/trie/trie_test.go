package trie

import (
	"testing"
)

func TestTrieInsertSearch(t *testing.T) {
	t.Parallel()
	tr := New(false)
	tr.Insert("hello")
	tr.Insert("world")

	if !tr.Search("hello") {
		t.Error("expected to find 'hello'")
	}
	if !tr.Search("world") {
		t.Error("expected to find 'world'")
	}
	if tr.Search("hell") {
		t.Error("should not find 'hell' (only prefix)")
	}
	if tr.Search("notfound") {
		t.Error("should not find 'notfound'")
	}
}

func TestTrieStartsWith(t *testing.T) {
	t.Parallel()
	tr := New(false)
	tr.Insert("hello")
	tr.Insert("help")

	if !tr.StartsWith("hel") {
		t.Error("StartsWith('hel') should be true")
	}
	if !tr.StartsWith("hello") {
		t.Error("StartsWith('hello') should be true")
	}
	if tr.StartsWith("world") {
		t.Error("StartsWith('world') should be false")
	}
}

func TestTrieDelete(t *testing.T) {
	t.Parallel()
	tr := New(false)
	tr.Insert("hello")
	tr.Insert("help")

	if err := tr.Delete("help"); err != nil {
		t.Fatal(err)
	}
	if tr.Search("help") {
		t.Error("should not find 'help' after delete")
	}
	if !tr.Search("hello") {
		t.Error("should still find 'hello'")
	}
	if !tr.StartsWith("hel") {
		t.Error("should still have words with prefix 'hel'")
	}
}

func TestTrieDeleteNonExistent(t *testing.T) {
	t.Parallel()
	tr := New(false)
	tr.Insert("hello")
	if err := tr.Delete("world"); err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestTrieCaseInsensitive(t *testing.T) {
	t.Parallel()
	tr := New(true)
	tr.Insert("Hello")
	if !tr.Search("hello") {
		t.Error("should find 'hello' case-insensitively")
	}
	if !tr.Search("HELLO") {
		t.Error("should find 'HELLO' case-insensitively")
	}
	if !tr.Search("Hello") {
		t.Error("should find 'Hello' case-insensitively")
	}
}

func TestTrieEmptyString(t *testing.T) {
	t.Parallel()
	tr := New(false)
	tr.Insert("")
	if !tr.Search("") {
		t.Error("should find empty string")
	}
}

func TestTrieCountWordsWithPrefix(t *testing.T) {
	t.Parallel()
	tr := New(false)
	words := []string{"a", "ab", "abc", "abcd", "abcde"}
	for _, w := range words {
		tr.Insert(w)
	}
	tests := []struct {
		prefix string
		want   int
	}{
		{"a", 5},
		{"ab", 4},
		{"abc", 3},
		{"abcd", 2},
		{"abcde", 1},
		{"x", 0},
	}
	for _, tc := range tests {
		t.Run(tc.prefix, func(t *testing.T) {
			if got := tr.CountWordsWithPrefix(tc.prefix); got != tc.want {
				t.Errorf("CountWordsWithPrefix(%q) = %d, want %d", tc.prefix, got, tc.want)
			}
		})
	}
}

func TestTrieUnicode(t *testing.T) {
	t.Parallel()
	tr := New(false)
	tr.Insert("café")
	if !tr.Search("café") {
		t.Error("should find 'café'")
	}
	if tr.Search("cafe") {
		t.Error("should not find 'cafe' (different unicode)")
	}
	if !tr.StartsWith("caf") {
		t.Error("StartsWith('caf') should be true")
	}
	if got := tr.CountWordsWithPrefix("caf"); got != 1 {
		t.Errorf("CountWordsWithPrefix('caf') = %d, want 1", got)
	}
}

func TestTrieInsertDuplicate(t *testing.T) {
	t.Parallel()
	tr := New(false)
	tr.Insert("hello")
	tr.Insert("hello")
	if !tr.Search("hello") {
		t.Error("should find 'hello'")
	}
	if err := tr.Delete("hello"); err != nil {
		t.Fatal(err)
	}
	if tr.Search("hello") {
		t.Error("should not find after one delete if we inserted twice")
	}
}
