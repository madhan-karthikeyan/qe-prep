# Database Utilities
Difficulty: Hard
Estimated Interview Time: 45 min
Prerequisites: SQL, transactions, in-memory data structures

## Problem Statement
Implement an in-memory database and generic CRUD operations on top of it. Support table operations and transactions with commit/rollback.

## Requirements
- InMemoryDatabase: create/drop table, insert, select, update, delete
- Transactions with begin/commit/rollback using undo log
- Thread-safe with read/write locks
- CrudOperations: typed CRUD helper with findById, updateById, deleteById

## Implementation Notes
- Uses ConcurrentHashMap for table storage
- ReentrantReadWriteLock for concurrent reads, exclusive writes
- Transaction undo log captures state for rollback
- CRUD operations delegate to InMemoryDatabase

## Test Strategy
- CRUD operations on tables
- Transaction commit persists changes
- Transaction rollback reverts changes
- Nested transaction rejection
- State tracking (isInTransaction)

## Edge Cases
- Duplicate table creation
- Operations on non-existent tables
- Commit/rollback without active transaction
- Null column definitions

## Failure Cases
- Duplicate table → IllegalArgumentException
- Non-existent table → IllegalArgumentException
- No active transaction → IllegalStateException
- Null parameters → NullPointerException

## Complexity (Time + Space)
- Insert: O(1) amortized
- Select with predicate: O(n)
- Update/Delete with predicate: O(n)
- Space: O(rows) for data, O(operations) for undo log

## Progression Path
Start with table creation and simple CRUD, then add transactions, then add query predicates.

## Common Interview Follow-ups
- Index support for faster queries
- JOIN operations between tables
- Multi-version concurrency control (MVCC)

## Possible Production Improvements
- B-tree indexes for O(log n) lookups
- Query optimizer for predicate pushdown
- Persistence to disk with WAL
