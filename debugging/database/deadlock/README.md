## Symptoms
The program hangs indefinitely. Both threads are stuck waiting for locks held by the other, as shown by the "DEADLOCK DETECTED" output when a timeout is used.

## Root Cause
Transaction A acquires lock A then lock B; Transaction B acquires lock B then lock A. When both hold their first lock and wait for the other's, neither can proceed — a classic deadly embrace.

## Fix
Both transactions acquire locks in the same global order (A then B). This eliminates the circular wait condition.

## Prevention
- Always acquire multiple locks in a consistent, documented order.
- Use `with` statements to ensure locks are released even on exceptions.
- Consider `threading.Lock.acquire(timeout=...)` and backoff strategies.
- Use lock hierarchies or deadlock detection mechanisms.
