# Database Utilities
Difficulty: Hard
Estimated Interview Time: 45 min
Prerequisites: SQL basics, data structures, map, slice, mutex

## Problem Statement
Implement an in-memory database with SQL-like operations and generic CRUD
utilities.

## Requirements
- In-memory map-based storage
- Execute: CREATE TABLE, INSERT, UPDATE, DELETE, DROP TABLE
- Query: SELECT with WHERE conditions
- Transaction: Begin, Commit, Rollback
- Generic CRUD: Create, Read, Update, Delete

## Implementation Notes
- SQL tokenizer and simple parser
- Parameterized queries with "?" placeholders
- Thread-safe via sync.RWMutex
- CRUD functions build SQL from map conditions

## Test Strategy
- Unit tests for CRUD operations
- Transaction commit and rollback
- Edge cases: empty tables, missing keys

## Edge Cases
- Empty result sets
- NULL values
- Multiple conditions
- String escaping in queries

## Failure Cases
- Table not found
- Duplicate table creation
- Unterminated string literals

## Complexity (Time + Space)
- Insert: O(1) amortized
- Select with WHERE: O(n) where n = rows
- Delete/Update: O(n)
- Space: O(rows * columns)

## Progression Path
- Add support for JOINs
- Add indexes for faster lookups
- Add GROUP BY and aggregation

## Common Interview Follow-ups
- SQL injection prevention
- Transaction isolation levels
- Query optimization

## Possible Production Improvements
- Use SQLite via database/sql
- Add connection pooling
- Implement WAL for durability
