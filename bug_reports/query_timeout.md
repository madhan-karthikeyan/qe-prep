# SQL query timeout with missing index on JOIN column

**Severity:** Major
**Priority:** P1
**Environment:** PostgreSQL 14+, 10M+ rows in orders table, 5M+ rows in customers table
**Component:** Reporting service — Order history query

## Summary

A JOIN query between `orders` and `customers` on `customer_id` takes 30+ seconds on 10M rows, causing connection pool exhaustion and cascading timeouts across the reporting service. The `customer_id` foreign key column in `orders` has no index.

## Steps to Reproduce

1. Create tables with 10M orders and 5M customers
2. Run the following query:
   ```sql
   SELECT o.id, o.total, o.created_at, c.name, c.email
   FROM orders o
   JOIN customers c ON o.customer_id = c.id
   WHERE o.created_at >= '2025-01-01'
   ORDER BY o.created_at DESC
   LIMIT 100;
   ```
3. Note execution time

## Expected Behavior

Query completes in <100ms with proper indexing.

## Actual Behavior

- Query takes 35-45 seconds
- Sequential scan on `orders` (10M rows)
- Hash JOIN builds hash table from 5M customers
- Connection pool (20 connections) is exhausted within 2 seconds under moderate load
- Downstream services receive `SQLException: timeout`

## Logs / Screenshots

```
EXPLAIN ANALYZE output:
------------------------------------------------------------------------------
 Limit  (cost=452031.20..452031.45 rows=100) (actual time=35240.12..35240.15)
   ->  Sort  (cost=452031.20..452562.08 rows=212352) (actual time=35240.10..35240.12)
         Sort Key: o.created_at DESC
         Sort Method: top-N heapsort  Memory: 57kB
         ->  Hash Join  (cost=83453.40..446861.50 rows=212352) (actual time=1201.45..34892.30)
               Hash Cond: (o.customer_id = c.id)
               ->  Seq Scan on orders o  (cost=0.00..321542.40 rows=2123520)
                     Filter: (created_at >= '2025-01-01'::date)
               ->  Hash  (cost=52142.00..52142.00 rows=2500000) (actual time=1198.20..1198.20)
                     ->  Seq Scan on customers c  (cost=0.00..52142.00 rows=2500000)
```

Key observations:
- `Seq Scan on orders`: Full table scan of 10M rows
- `Hash Join` cost dominates (34.8s actual)
- No index on `orders.customer_id`

## Root Cause Analysis

The `customer_id` column in `orders` was defined as a foreign key constraint but was **not indexed**:

```sql
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES customers(id),  -- FK constraint, but NO index!
    total DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT NOW()
);
```

PostgreSQL does **not** automatically index foreign key columns. Without an index, every JOIN between `orders` and `customers` requires a full sequential scan of the `orders` table, plus a hash build over `customers`.

The sequential scan of 10M rows at ~300MB/s takes ~30 seconds on standard disk I/O. The hash join also requires building an in-memory hash table of 5M customer rows.

## Fix

Add a composite index covering the join column and filter column:

```sql
CREATE INDEX idx_orders_customer_id_created_at 
ON orders(customer_id, created_at DESC);

-- Also index the FK column alone (for other queries)
CREATE INDEX idx_orders_customer_id 
ON orders(customer_id);
```

After adding the index:

```
EXPLAIN ANALYZE output:
------------------------------------------------------------------------------
 Limit  (cost=1.42..482.34 rows=100) (actual time=0.35..2.14 rows=100)
   ->  Nested Loop  (cost=1.42..1022345.12 rows=212352) (actual time=0.35..2.10)
         ->  Index Scan Backward using idx_orders_created_at on orders o
               (cost=0.43..872345.10 rows=2123520) (actual time=0.20..1.50)
               Index Cond: (created_at >= '2025-01-01')
         ->  Index Scan using customers_pkey on customers c
               (cost=0.42..0.68 rows=1) (actual time=0.01..0.01)
               Index Cond: (id = o.customer_id)
```

Execution time: **2.1ms** (down from 35,240ms) — a **16,000× improvement**.

## Regression Tests

### 1. Query Performance Benchmark in CI

```python
def test_order_history_query_performance():
    # Connect to test DB with real-size data
    conn = psycopg2.connect(TEST_DATABASE_URL)
    cur = conn.cursor()
    
    # Warm cache (run once)
    cur.execute("""
        SELECT o.id, o.total, o.created_at, c.name, c.email
        FROM orders o
        JOIN customers c ON o.customer_id = c.id
        WHERE o.created_at >= '2025-01-01'
        ORDER BY o.created_at DESC
        LIMIT 100
    """)
    
    # Measure execution time
    start = time.perf_counter()
    cur.execute("""
        SELECT o.id, o.total, o.created_at, c.name, c.email
        FROM orders o
        JOIN customers c ON o.customer_id = c.id
        WHERE o.created_at >= '2025-01-01'
        ORDER BY o.created_at DESC
        LIMIT 100
    """)
    elapsed = time.perf_counter() - start
    
    # Assert: query completes in <50ms
    assert elapsed < 0.05, f"Query took {elapsed:.3f}s (threshold: 0.05s)"
    
    # Verify result set
    rows = cur.fetchall()
    assert len(rows) == 100 or len(rows) <= 100
```

### 2. Missing Index Detection (Static Analysis)

```python
def test_foreign_keys_have_indexes():
    conn = psycopg2.connect(TEST_DATABASE_URL)
    cur = conn.cursor()
    
    # Find FK columns without indexes
    cur.execute("""
        SELECT
            con.conrelid::regclass AS table_name,
            con.conkey AS fk_columns,
            array_agg(ind.indkey) AS indexed_columns
        FROM pg_constraint con
        LEFT JOIN pg_index ind 
            ON ind.indrelid = con.conrelid
            AND con.conkey && ind.indkey
        WHERE con.contype = 'f'
        GROUP BY con.conrelid, con.conkey
        HAVING array_agg(ind.indkey) IS NULL
           OR NOT (con.conkey <@ ANY(array_agg(ind.indkey)))
    """)
    
    unindexed = cur.fetchall()
    assert len(unindexed) == 0, f"FK columns without indexes: {unindexed}"
```

### 3. CI Integration

Add to CI pipeline:
- Run query performance benchmark as part of integration tests
- Fail build if query time exceeds 50ms baseline
- Run FK index check on every schema migration
- Alert if `EXPLAIN` shows sequential scans on tables > 1M rows
