# File Processing

Difficulty: Medium
Estimated Interview Time: 30 min
Prerequisites: File I/O, generators

## Problem Statement

Implement a CSV parser that handles quoted fields, escaped quotes, custom delimiters, and streaming. Implement a `wc`-style word count utility supporting stdin, file input, and multi-file totals.

## Requirements

CSV Parser:
- Quoted fields, escaped quotes, commas inside quotes
- Custom delimiter
- Header row detection (optional)
- Empty fields and trailing commas
- Streaming (yields rows one at a time)

Word Count:
- Lines, words, character count
- Stdin and file input
- Unicode-aware character counting
- Multiple file mode with totals row

## Implementation Notes

- CSV parser wraps stdlib `csv.reader` for correctness, adds streaming via generator
- Word count uses `len(line)` for character count (Unicode-aware in Python 3)
- Both support stdin via sys.stdin when no files provided

## Test Strategy
- Unit (CSV): simple, quoted, escaped quotes, custom delimiter, no header, empty fields, trailing comma, empty file, streaming
- Unit (wc): empty, single line, multi-line, unicode, blank lines
- Fuzz: random CSV strings, extremely long lines, special characters

## Edge Cases

- Empty file returns no rows
- File with only a header row returns one row
- Trailing commas produce empty-string fields
- Very long fields (10k+ chars)
- Null bytes and unicode in CSV

## Failure Cases

- Nonexistent file (word_count prints to stderr)
- Directory passed as file (word_count prints to stderr)
- Malformed CSV with unclosed quotes (stdlib raises, caller handles)

## Complexity
- Time: O(n) for parsing/counting (n = input size)
- Space: O(1) streaming, O(k) for k fields in a row

## Progression Path
Basic → Streaming parser → Memory-mapped for huge files → Parallel processing

## Common Interview Follow-ups

- How would you handle GB-sized CSV files?
- How would you process CSV files in parallel?
- How would you detect CSV dialect automatically?
- How would you handle malformed CSV rows gracefully?

## Possible Production Improvements

- Memory-mapped file reading for large files
- Parallel chunked processing with multiprocessing
- Automatic delimiter/quote detection
- Configurable error handling (skip bad rows)
