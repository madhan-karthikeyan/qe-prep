package trie

import (
	"errors"
	"strings"
	"unicode"
)

// ErrNotFound is returned by Delete when the word is not in the trie.
var ErrNotFound = errors.New("word not found in trie")

// node represents a single node in the trie.
type node struct {
	children map[rune]*node
	isEnd    bool
	count    int // number of words passing through this node
}

// Trie is a prefix tree supporting case-insensitive and unicode rune-based
// operations.
type Trie struct {
	root            *node
	caseInsensitive bool
}

// New creates a new Trie. Set caseInsensitive to true to ignore case.
func New(caseInsensitive bool) *Trie {
	return &Trie{
		root:            &node{children: make(map[rune]*node)},
		caseInsensitive: caseInsensitive,
	}
}

// normalizeRune converts a rune to lowercase if the trie is case-insensitive.
func (t *Trie) normalizeRune(r rune) rune {
	if t.caseInsensitive {
		return unicode.ToLower(r)
	}
	return r
}

// normalizeWord normalizes the word for case-insensitive mode.
func (t *Trie) normalizeWord(word string) string {
	if t.caseInsensitive {
		return strings.ToLower(word)
	}
	return word
}

// Insert adds a word to the trie.
func (t *Trie) Insert(word string) {
	word = t.normalizeWord(word)
	n := t.root
	for _, r := range word {
		r = t.normalizeRune(r)
		child, ok := n.children[r]
		if !ok {
			child = &node{children: make(map[rune]*node)}
			n.children[r] = child
		}
		n.count++
		n = child
	}
	n.count++
	n.isEnd = true
}

// Search returns true if the word was inserted.
func (t *Trie) Search(word string) bool {
	word = t.normalizeWord(word)
	n := t.traverse(word)
	return n != nil && n.isEnd
}

// StartsWith returns true if any inserted word has the given prefix.
func (t *Trie) StartsWith(prefix string) bool {
	prefix = t.normalizeWord(prefix)
	return t.traverse(prefix) != nil
}

// Delete removes a word from the trie. Returns ErrNotFound if the word does not
// exist.
func (t *Trie) Delete(word string) error {
	word = t.normalizeWord(word)
	if !t.Search(word) {
		return ErrNotFound
	}
	t.delete(t.root, []rune(word), 0)
	return nil
}

// delete recursively removes a word from the trie.
func (t *Trie) delete(n *node, runes []rune, depth int) {
	if depth == len(runes) {
		n.isEnd = false
		n.count--
		return
	}
	r := t.normalizeRune(runes[depth])
	child := n.children[r]
	t.delete(child, runes, depth+1)
	n.count--
	if child.count == 0 && !child.isEnd {
		delete(n.children, r)
	}
}

// traverse follows the given key and returns the last node reached, or nil if
// the key is not in the trie.
func (t *Trie) traverse(key string) *node {
	n := t.root
	for _, r := range key {
		r = t.normalizeRune(r)
		child, ok := n.children[r]
		if !ok {
			return nil
		}
		n = child
	}
	return n
}

// countWordsFromNode returns the number of complete words in the subtree rooted
// at n.
func countWordsFromNode(n *node) int {
	if n == nil {
		return 0
	}
	total := 0
	if n.isEnd {
		total++
	}
	for _, child := range n.children {
		total += countWordsFromNode(child)
	}
	return total
}

// CountWordsWithPrefix returns the number of inserted words that have the given
// prefix.
func (t *Trie) CountWordsWithPrefix(prefix string) int {
	prefix = t.normalizeWord(prefix)
	n := t.traverse(prefix)
	return countWordsFromNode(n)
}
