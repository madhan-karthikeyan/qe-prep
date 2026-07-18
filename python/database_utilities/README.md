# Database Utilities
Difficulty: Medium
Estimated Interview Time: 30 min
Prerequisites: SQL, sqlite3, context managers

## Problem Statement
Build an in-memory/SQLite database wrapper with a generic CRUD utility, transaction support, and context manager integration.

## Requirements
- connect(), close(), execute(), executemany(), fetchone(), fetchall()
- Transaction support (begin/commit/rollback) + context manager
- Connection as context manager
- Generic CRUD with parameterized queries (no SQL injection)
- SQLite file-based persistence

## Implementation Notes
- Wraps sqlite3 with proper connection management
- CRUD uses parameterized queries throughout
- Transaction context manager auto-commits or rolls back

## Test Strategy (Unit/Integration)
- Unit: CRUD operations, transactions, rollback, context managers
- Integration: SQLite file-based persistence across connections

## Edge Cases
- Connecting twice
- fetchone with no results
- Table existence checks
- Transaction rollback on exception

## Failure Cases
- Invalid SQL → sqlite3.OperationalError
- Missing table → sqlite3.OperationalError
- Connection closed → RuntimeError

## Complexity
- All operations O(1) overhead over SQLite

## Progression Path
- Basic: raw SQL execution
- Intermediate: wrapper with connection management
- Advanced: CRUD with ORM-like interface
- Production: connection pooling, migration system

## Common Interview Follow-ups
- How would you implement lazy connections?
- How would you add connection pooling?
- How would you handle concurrent writes?

## Possible Production Improvements
- Connection pooling with max connections
- Schema migration system
- Query logging and metrics
- Read replica support
