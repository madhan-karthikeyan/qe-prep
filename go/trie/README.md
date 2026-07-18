# Trie
Difficulty: Medium
Estimated Interview Time: 25 min
Prerequisites: trees, recursion, prefix search

## Problem Statement
Implement a trie (prefix tree) for efficient word storage and prefix-based queries.

## Requirements
- Insert, Search, StartsWith, Delete, CountWordsWithPrefix
- Case-insensitive option
- Unicode rune-based nodes

## Implementation Notes
- Each node has a map[rune]*node for children
- Insert/Search/StartsWith all O(k) where k is key length
- Delete recursively cleans up unused nodes
- CountWordsWithPrefix traverses subtree and counts isEnd markers
- Case-insensitive mode normalizes to lowercase

## Test Strategy
- Insert/Search/Delete round-trip
- Pre-existing prefix after partial delete
- Case-insensitive matching
- Unicode (multi-byte runes)
- CountWordsWithPrefix accuracy
- Empty string as valid word

## Edge Cases
- Empty string insertion and search
- Delete non-existent word
- Insert duplicate word
- Words that share prefixes

## Failure Cases
- Delete of non-existent word returns ErrNotFound

## Complexity (Time + Space)
- Insert: O(k) time, O(k) space per new node
- Search/StartsWith: O(k) time, O(1) space
- Delete: O(k) time, O(1) space
- CountWordsWithPrefix: O(n) where n is subtree size
- Space: O(total characters in all words)

## Progression Path (Basic → Intermediate → Advanced → Production)
- Basic: ASCII-only trie with insert/search
- Intermediate: Delete, StartsWith, case-insensitive
- Advanced: Unicode support, prefix count
- Production: Compressed trie (radix tree), on-disk storage

## Common Interview Follow-ups
- How would you implement autocomplete?
- How would you compress the trie?
- What are the space trade-offs vs a hash set?

## Possible Production Improvements
- Implement a radix tree (compressed trie) for memory efficiency
- Add serialization/deserialization
- Add fuzzy/approximate search
- Implement autocomplete with top-k suggestions
