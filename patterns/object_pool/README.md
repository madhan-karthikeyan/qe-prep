# Object Pool

## Intent
Reuse expensive-to-create objects by maintaining a pool of ready instances.

## Structure
- `ObjectPool` — thread-safe pool with acquire/release lifecycle
- `factory` — creates new objects when pool is empty
- `validator` — optional check on release; invalid objects are discarded

## When to Use
- Database connections, socket connections, large buffers
- When object creation is expensive and reuse is safe

## When NOT to Use
- Lightweight objects (just create new ones)
- Objects with thread-unsafe state that isn't reset between uses

## Implementation Notes
- Blocking `acquire` with configurable timeout
- Pool grows on demand up to `max_size`; blocks once saturated
- Validator ensures objects are in a usable state on release
- `queue.Queue` provides thread-safe FIFO semantics

## Tradeoffs
- **Pros** — Reduces allocation overhead, predictable resource usage
- **Cons** — Object reset complexity, tuning pool size, stale connections

## Language Notes
- **Go** — `sync.Pool` (GC-friendly, no size limit); or `buffered channel` for fixed pools
- **Java** — `commons-pool2` (GenericObjectPool); connection pools in JDBC/HikariCP

## Related Patterns
- [worker_pool](../worker_pool) — pools threads instead of objects
- [rate_limiter](../rate_limiter) — limits resource consumption rate
