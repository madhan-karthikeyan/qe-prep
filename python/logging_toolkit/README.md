# Logging Toolkit

Difficulty: Easy
Estimated Interview Time: 20 min
Prerequisites: threading, file I/O

## Problem Statement

Implement a thread-safe logging library with multiple output targets, log levels, log rotation, and configurable formatting/filtering.

## Requirements

- Multiple output targets (stdout via StringIO, file via RotatingFileHandler)
- Log levels: DEBUG, INFO, WARN, ERROR
- Thread-safe writes using threading.Lock
- Log rotation by file size with configurable backup count
- Configurable format template (timestamp, level, message)
- Filter by minimum level and/or regex pattern

## Implementation Notes

- LogRecord dataclass holds individual log entries with formatting support
- LogFilter supports both minimum-level and regex-pattern filtering
- RotatingFileHandler handles size-based rotation with numbered backups
- Logger holds a list of handlers and a threading.Lock for safe concurrent writes

## Test Strategy
- Unit: LogRecord formatting, LogFilter acceptance, Logger output to StringIO, multiple handlers, handler removal
- Integration: File rotation with temp files, multi-threaded writes, rotation restoration
- Stress: 8 threads writing 1250 entries each (10k total), verify no duplicates

## Edge Cases

- Logging when no handlers are attached
- Rotating when file is empty or missing
- Thread contention on the same handler
- Very rapid successive log writes during rotation

## Failure Cases

- Disk full during rotation (write will raise, propagated to caller)
- Invalid regex pattern in LogFilter
- Removing a handler that was never added

## Complexity
- Time: O(1) per log call (amortized; rotation is O(n) but infrequent)
- Space: O(1) per Record, O(k) for k backup files

## Progression Path
Basic → Add JSON structured logging → Add async logging → Production-grade with queue-based async writer

## Common Interview Follow-ups

- How would you make logging asynchronous?
- How would you add structured (JSON) output?
- How would you handle disk-full scenarios?
- How would you implement hierarchical loggers?

## Possible Production Improvements

- Use queue-based async writer to avoid blocking application threads
- Support structured logging (JSON, logfmt)
- Add log sampling/rate-limiting
- Add remote log shipping (syslog, network)
