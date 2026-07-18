# Retry

## Intent
Automatically retry a failed operation with configurable backoff and jitter to handle transient failures.

## Structure
- `Retry` — configurable retry executor
- `RetryError` — raised when all attempts exhausted

## When to Use
- Network calls, I/O operations, distributed service calls
- When failures are transient and self-correcting

## When NOT to Use
- Idempotency is not guaranteed on the callee side
- Failures are deterministic (e.g., invalid input)

## Implementation Notes
- Exponential backoff: `base * factor^(attempt-1)`, capped at `max_delay`
- Full jitter: `random(0, delay)` to spread retry storms
- `on_retry` callback enables logging/metrics without coupling

## Tradeoffs
- **Pros** — Simple, battle-tested, low overhead
- **Cons** — Blocks the caller during sleep; consider async for high concurrency

## Language Notes
- **Go** — `time.Sleep` in a loop; `math/rand` for jitter; `context` for cancellation
- **Java** — Thread.sleep in a loop; `ThreadLocalRandom` for jitter; `spring-retry` for declarative approach

## Related Patterns
- [exponential_backoff](../exponential_backoff) — pure delay calculation, reusable by Retry
- [circuit_breaker](../circuit_breaker) — prevents calls when system is known to be down
