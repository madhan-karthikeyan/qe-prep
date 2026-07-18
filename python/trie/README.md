# Trie

Difficulty: Medium
Estimated Interview Time: 25 min
Prerequisites: Trees, recursion, string algorithms

## Problem Statement

Implement a Trie (prefix tree) supporting insert, search, starts_with, and delete operations with an optional case-insensitive mode.

## Requirements

- insert(word), search(word), starts_with(prefix), delete(word)
- Case-insensitive option (normalize to lowercase)
- Count words with given prefix
- Delete reduces size counter correctly

## Implementation Notes

- _TrieNode stores dict of children and a boolean is_end flag
- Case-insensitive mode lower()s all input
- Delete recursively removes nodes that are no longer needed
- _count_words does a DFS from a given node to count all words in that subtree

## Test Strategy
- Unit: insert/search, starts_with, delete, overlapping prefixes, empty string, case sensitivity toggle, count_prefix, size tracking

## Edge Cases

- Empty string insert/search/delete
- Deleting a non-existent word returns False
- Overlapping prefixes: deleting "app" does not affect "apple"
- Case-insensitive: "Hello" matches "hello" and "HELLO"

## Failure Cases

- None or non-string input (TypeError from lower())
- Deleting a word that is a prefix of another word

## Complexity
- Time: O(L) for all operations (L = word length)
- Space: O(N * L) for N words of average length L

## Progression Path
Basic → Case-insensitive → Delete → Count prefix → Compressed trie (radix tree)

## Common Interview Follow-ups

- How would you implement autocomplete (return all words with a prefix)?
- How would you compress the trie into a radix tree?
- How would you find the longest common prefix of all words?
- How would you handle 1 million words efficiently?

## Possible Production Improvements

- Compressed/radix tree for memory efficiency
- Serialization/deserialization
- Aho-Corasick automaton for pattern matching
- Unicode normalization
- Persistent trie for concurrent access
