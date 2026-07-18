# File Processing
Difficulty: Medium
Estimated Interview Time: 30 min
Prerequisites: io.Reader, strings, unicode

## Problem Statement
Implement a CSV parser and a wc-style word count utility.

## Requirements
- CSV: quoted fields, escaped quotes, configurable delimiter, header detection
- WordCount: lines, words, characters, bytes from io.Reader
- Support multiple files for wc

## Implementation Notes
- CSV parser is rune-aware for Unicode support
- Word counting uses unicode.IsSpace for proper Unicode whitespace detection
- Both modules accept io.Reader for testability

## Test Strategy
- Table-driven CSV tests for quoting, escaping, delimiters
- Fuzz test for CSV parser robustness
- Word count tests for Unicode, empty input, single/multi-line

## Edge Cases
- Empty fields in CSV
- Unterminated quotes
- Carriage return + newline line endings
- Mixed Unicode and ASCII in word count

## Failure Cases
- Unterminated quoted field
- Unexpected quote outside field

## Complexity (Time + Space)
- CSV: O(n) time, O(n) space for full parse; O(1) for streaming callback
- WordCount: O(n) time, O(1) space

## Progression Path (Basic → Intermediate → Advanced → Production)
- Basic: Simple split-by-delimiter CSV
- Intermediate: Quoted field support
- Advanced: Streaming with callback
- Production: Memory-mapped files, encoding detection

## Common Interview Follow-ups
- How would you handle very large files?
- What encoding issues might arise?
- How would you implement error recovery for malformed CSV?

## Possible Production Improvements
- Encoding detection (UTF-8 vs Latin-1)
- Memory-mapped file support for large files
- Parallel processing for multiple files
