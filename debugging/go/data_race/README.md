## Symptoms
Running without `-race` may appear to work correctly. With `-race`, the Go race detector reports concurrent unsafe reads and writes to the map. Without the detector, the program can crash with `fatal error: concurrent map writes`.

## Root Cause
Go maps are not safe for concurrent access. Multiple goroutines write to (and read from) the same map without synchronization, causing a data race.

## Fix
Wrap map accesses with a `sync.RWMutex`. Use `Lock`/`Unlock` for writes and `RLock`/`RUnlock` for reads to allow concurrent reads while ensuring exclusive writes.

## Prevention
- Always use `-race` during development and in CI.
- Use `sync.Map` for simple cases, or protect maps with a mutex.
- Document thread-safety guarantees in API comments.
