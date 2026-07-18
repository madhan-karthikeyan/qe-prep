# Logging Toolkit
Difficulty: Easy
Estimated Interview Time: 20 min
Prerequisites: io.Writer, sync.Mutex, iota

## Problem Statement
Implement a thread-safe logging toolkit with multiple output targets, log levels, configurable formatting, file rotation, and filtering.

## Requirements
- Logger: multiple io.Writer targets, DEBUG/INFO/WARN/ERROR levels, fmt.Sprintf-style formatting
- RotatingFileWriter: size-based rotation with timestamp archiving
- Filter: minimum level and regex pattern filtering
- Thread safety throughout

## Implementation Notes
- Levels use iota for compact definition
- Mutex protects shared writers and internal state
- RotatingFileWriter implements io.WriteCloser for composability

## Test Strategy
- Table-driven tests for logger formatting
- Concurrent write tests with goroutines
- File rotation integrity checks
- Filter matching for level and regex

## Edge Cases
- Write after Close returns error
- Zero writers (default to io.Discard)
- Concurrent writes during rotation

## Failure Cases
- Disk full during rotation
- Invalid regex pattern

## Complexity (Time + Space)
- Logger: O(1) per write (amortized)
- Rotation: O(1) write, O(n) rename where n is file size
- Filter: O(m) regex match where m is message length

## Progression Path (Basic → Intermediate → Advanced → Production)
- Basic: Single-writer logger with levels
- Intermediate: Multi-writer with rotation
- Advanced: Async/batched writes, structured JSON logging
- Production: Log shipping, compression, retention policies

## Common Interview Follow-ups
- How would you add structured (JSON) logging?
- How would you implement async writes?
- How would you handle backpressure?

## Possible Production Improvements
- Add a TTL-based retention policy
- Compress archived files
- Structured (JSON) output format
