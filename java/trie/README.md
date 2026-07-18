# Trie
Difficulty: Medium
Estimated Interview Time: 40 min
Prerequisites: Tree traversal, recursion

## Problem Statement
Implement a trie (prefix tree) with insert, search, startsWith, delete, and countPrefix operations.

## Requirements
- Character-level nodes
- Case-insensitive option
- countPrefix returns number of complete words under prefix

## Implementation Notes
- HashMap-based children for O(1) per character
- Recursive delete that prunes unused branches
- normalize() method for case handling

## Test Strategy
- Insert/search verification
- Prefix matching
- Delete with branch cleanup
- countPrefix accuracy
- Case sensitivity toggle

## Edge Cases
- Empty string as a key
- Non-existent delete (returns false)
- Null arguments

## Failure Cases
- IllegalArgumentException on null insert/delete

## Complexity
- Time: O(L) per operation, L = word length
- Space: O(N * L) for N words of average length L

## Progression Path
1. Insert/Search → 2. StartsWith → 3. Delete → 4. CountPrefix → 5. Wildcard search

## Common Interview Follow-ups
- Implement wildcard search (. or *)
- How to serialize/deserialize a trie?
- Autocomplete with top-k suggestions?

## Possible Production Improvements
- Array-based children for ASCII-only
- Weighted nodes for autocomplete ranking
- Compressed/radix tree variant
