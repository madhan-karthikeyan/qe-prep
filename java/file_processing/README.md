# File Processing
Difficulty: Easy
Estimated Interview Time: 30 min
Prerequisites: Java I/O

## Problem Statement
Implement CSV parser (streaming) and word counter (wc clone).

## Requirements
- CsvParser: Iterator<String[]>, quoted fields with escaping, custom delimiter
- WordCounter: count lines, words, characters from Reader

## Implementation Notes
- CsvParser implements Iterator and AutoCloseable
- Quoted fields support escaped quotes ("" → ")
- WordCounter uses character-level whitespace detection

## Test Strategy
- Unit tests for parsing, quoting, edge cases
- Unicode handling

## Edge Cases
- Empty file, blank lines
- Quoted field containing delimiter
- Escaped quotes within quoted field

## Failure Cases
- NoSuchElementException on empty iterator
- IOException on reader failure

## Complexity
- Time: O(n) per file
- Space: O(f) where f = max field count in a row

## Progression Path
1. Simple CSV → 2. Quoted fields → 3. Custom delimiter → 4. Streaming

## Common Interview Follow-ups
- How would you handle multi-line quoted fields?
- How would you handle encoding detection?
- Streaming vs loading entirely in memory?

## Possible Production Improvements
- Apache Commons CSV-style configuration
- Multi-line quoted fields
- Header mapping to objects
