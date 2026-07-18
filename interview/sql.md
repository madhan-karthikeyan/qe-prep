# SQL Interview Guide — QE Engineer

## Overview

SQL is a core QE skill — for data validation, test setup, and analyzing test results. Expect questions on joins, aggregations, window functions, query optimization, and index usage.

## Top 25 SQL Questions

### Joins (Difficulty: ★☆☆ Easy)

1. **What's the difference between INNER JOIN and LEFT JOIN?**
   - **INNER**: Only matching rows from both tables.
   - **LEFT**: All rows from left table, NULLs where right has no match.

2. **Explain CROSS JOIN with an example.**
   - **Expected**: Cartesian product of two tables. No `ON` clause. Useful for generating combinations (e.g., `products × stores` for inventory).

3. **What is a self-join? When would you use it?**
   - **Expected**: Joining a table to itself with aliases. Common for hierarchical data (e.g., employees with manager_id).

4. **Difference between NATURAL JOIN and USING?**
   - **NATURAL**: Joins on all columns with same name (dangerous — implicit).
   - **USING**: Explicitly names join columns. Prefer USING over NATURAL.

5. **Write a query to find employees who have never been assigned a project.**
   - **Expected**: `SELECT * FROM employees e LEFT JOIN assignments a ON e.id = a.emp_id WHERE a.emp_id IS NULL`

### Problem (Difficulty: ★★☆ Medium)

6. **Find duplicate emails in a users table.**
   ```sql
   SELECT email, COUNT(*) 
   FROM users 
   GROUP BY email 
   HAVING COUNT(*) > 1;
   ```

7. **Find the Nth highest salary (without LIMIT/OFFSET).**
   ```sql
   SELECT DISTINCT salary 
   FROM employees e1 
   WHERE N = (SELECT COUNT(DISTINCT salary) FROM employees e2 WHERE e2.salary >= e1.salary);
   ```

8. **Running total of sales by date.**
   ```sql
   SELECT order_date, amount,
          SUM(amount) OVER (ORDER BY order_date) AS running_total
   FROM orders;
   ```

9. **Find employees who earn more than their manager.**
   ```sql
   SELECT e.name 
   FROM employees e 
   JOIN employees m ON e.manager_id = m.id 
   WHERE e.salary > m.salary;
   ```

10. **Department-wise average salary (with average above company average).**
    ```sql
    SELECT department_id, AVG(salary) AS dept_avg 
    FROM employees 
    GROUP BY department_id 
    HAVING AVG(salary) > (SELECT AVG(salary) FROM employees);
    ```

### Indexes (Difficulty: ★★☆ Medium)

11. **How does a B-tree index work?**
    - **Expected**: Balanced tree structure. Root → branch → leaf (sorted). Search, insert, delete in O(log N). Supports range queries, ORDER BY, GROUP BY.

12. **What's the difference between B-tree and hash index?**
    - **B-tree**: Supports range queries (`>`, `<`, `BETWEEN`), ORDER BY. Hash: O(1) equality lookups only, no range support.

13. **When does an index NOT get used?**
    - Function on indexed column (`WHERE YEAR(date) = 2024`)
    - Leading wildcard (`LIKE '%abc'`)
    - Column type mismatch (string vs int)
    - Low cardinality (boolean column)

14. **What is a composite index? Column order matters — why?**
    - **Expected**: Index on multiple columns. Leftmost prefix rule: a query must use the first column to benefit. For `INDEX(a,b,c)`, queries on `a`, `(a,b)`, `(a,b,c)` use the index; `(b,c)` alone does not.

15. **Explain covering index vs included columns.**
    - **Expected**: Covering index contains ALL columns needed by a query (no table lookup). Included columns (SQL Server) store extra data at leaf level without affecting index order.

### Query Optimization (Difficulty: ★★☆ Medium–★★★ Hard)

16. **How would you diagnose a slow query?**
    - **Expected**: `EXPLAIN ANALYZE` (or `EXPLAIN`), check seq scans vs index scans, examine rows estimate vs actual, look for sorts/hash joins on large tables, check for missing indexes.

17. **What's the difference between `EXPLAIN` and `EXPLAIN ANALYZE`?**
    - **EXPLAIN**: Shows query plan with cost estimates.
    - **EXPLAIN ANALYZE**: Executes the query and shows actual vs estimated rows, timings. Slower but more accurate.

18. **How can you optimize a query with pagination (`LIMIT/OFFSET`)?**
    - **Problem**: Large OFFSET scans all skipped rows.
    - **Fix**: Keyset pagination — `WHERE id > last_seen_id ORDER BY id LIMIT 20`. Consistent and fast.

19. **What causes a "Using temporary; Using filesort" in MySQL EXPLAIN?**
    - **Expected**: Temporary — query needs to build temp table (GROUP BY without index, subquery). Filesort — cannot sort using index (ORDER BY on non-indexed column or mixed ASC/DESC).

20. **When would you denormalize a schema?**
    - **Expected**: Read-heavy workload, complex joins are slow, reporting/analytics queries need pre-joined data. Tradeoff: data redundancy, update complexity, storage cost.

### Window Functions (Difficulty: ★★★ Hard)

21. **Explain ROW_NUMBER(), RANK(), DENSE_RANK() differences.**
    ```sql
    -- ROW_NUMBER: 1,2,3,4 (unique, gaps for ties)
    -- RANK:       1,1,3,4 (same rank for ties, gap)
    -- DENSE_RANK: 1,1,2,3 (same rank for ties, no gap)
    ```

22. **Write a query to get the top 3 products by revenue per category.**
    ```sql
    SELECT category_id, product_id, revenue
    FROM (
      SELECT *, ROW_NUMBER() OVER (PARTITION BY category_id ORDER BY revenue DESC) AS rn
      FROM products
    ) ranked
    WHERE rn <= 3;
    ```

23. **Difference between ROWS and RANGE in window frame?**
    - **ROWS**: Physical rows (preceding N rows).
    - **RANGE**: Logical rows (rows within value range). `RANGE BETWEEN 5 PRECEDING AND 5 FOLLOWING`.

24. **Write a query to find gaps in a sequence (e.g., missing order IDs).**
    ```sql
    SELECT id + 1 AS gap_start, next_id - 1 AS gap_end
    FROM (
      SELECT id, LEAD(id) OVER (ORDER BY id) AS next_id FROM orders
    ) gaps
    WHERE next_id - id > 1;
    ```

25. **Calculate month-over-month revenue change.**
    ```sql
    SELECT month, revenue,
           LAG(revenue) OVER (ORDER BY month) AS prev_month_revenue,
           revenue - LAG(revenue) OVER (ORDER BY month) AS change
    FROM monthly_revenue;
    ```

---

## How to Approach SQL Problems

| Step | Action |
|------|--------|
| 1 | Clarify schema — tables, columns, relationships |
| 2 | Write query step by step — start with FROM + JOIN, then WHERE, GROUP BY, HAVING, SELECT, ORDER BY |
| 3 | Test with sample data mentally |
| 4 | Check edge cases: NULLs, empty tables, duplicates |
| 5 | Discuss performance — indexes, query plan |

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Forgetting NULL handling in joins | `COALESCE` or explicit NULL checks |
| Using `DISTINCT` as a crutch for bad joins | Fix the join condition |
| Not considering performance | Mention indexes, `EXPLAIN`, keyset pagination |
| Mixing WHERE and HAVING | WHERE filters rows, HAVING filters groups |
| Window function without PARTITION BY | Need partition to scope the window |

## Difficulty Levels

| Topic | Difficulty |
|-------|-----------|
| Basic joins (INNER, LEFT, RIGHT) | ★☆☆ |
| Aggregations + GROUP BY | ★★☆ |
| Subqueries (correlated vs non) | ★★☆ |
| Indexes + optimization | ★★☆ |
| Window functions | ★★★ |
| Query plan analysis | ★★★ |
