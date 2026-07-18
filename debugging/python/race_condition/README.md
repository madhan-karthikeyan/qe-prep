## Symptoms
The counter's final value is different every run and is always less than the expected 100000.

## Root Cause
`counter.value += 1` is a read-modify-write operation. Without a lock, two threads can read the same value before either writes, causing one increment to be lost. This is a classic data race on shared mutable state.

## Fix
Wrap the increment with a `threading.Lock` to ensure mutual exclusion: only one thread can read and write `value` at a time.

## Prevention
- Always use `threading.Lock` (or `RLock`) when multiple threads access shared mutable state.
- Consider `threading.atomic` or higher-level concurrency primitives.
- Use immutable data structures where possible.
