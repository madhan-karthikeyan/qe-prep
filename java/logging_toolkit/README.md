# Logging Toolkit
Difficulty: Medium
Estimated Interview Time: 45 min
Prerequisites: Java I/O, concurrency basics

## Problem Statement
Implement a thread-safe logging toolkit with configurable output targets, log levels, file rotation, and filtering.

## Requirements
- Logger with stdout and file output
- Log levels: TRACE, DEBUG, INFO, WARN, ERROR
- Thread-safe via ReentrantLock
- RotatingFileWriter: size-based rotation with timestamp archives
- LogFilter: filter by level, regex, or custom predicate

## Implementation Notes
- Uses ReentrantLock for thread safety
- RotatingFileWriter archives on rotation with UTC timestamp suffix
- Records use Java 21 record type

## Test Strategy
- Unit tests per class
- Multi-threaded logging verification with CountDownLatch
- File rotation verification

## Edge Cases
- Null/blank logger name
- Closed writer access
- Rotation boundary at exact byte limit

## Failure Cases
- IOException during file writes
- Invalid constructor arguments

## Complexity
- Time: O(1) per log call
- Space: O(1) per log call

## Progression Path
1. Basic logging → 2. File output → 3. Rotation → 4. Filtering → 5. Async logging

## Common Interview Follow-ups
- How would you implement async logging?
- What if multiple loggers share the same file?
- How would you add JSON formatting?

## Possible Production Improvements
- Async logging with bounded queue
- JSON or structured logging
- Configurable via properties file
- Logback-style appender architecture
