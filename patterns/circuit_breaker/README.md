# Circuit Breaker

## Intent
Prevent repeated calls to a failing service, giving it time to recover, and resume gracefully.

## Structure
- `CircuitBreaker` — thread-safe state machine wrapping function calls
- `CircuitBreakerState` — CLOSED, OPEN, HALF_OPEN
- `CircuitBreakerOpenError` — raised when calls are rejected

## When to Use
- Remote service calls, databases, external APIs
- Protecting callers from cascading failures

## When NOT to Use
- Local, deterministic operations that never fail transiently
- When fail-fast is preferred over degraded behavior

## Implementation Notes
- Thread-safe via `threading.Lock`
- Timeout uses `time.monotonic()` for wall-clock measurement
- Success in CLOSED resets failure count (continuous success heuristic)
- Failure in HALF_OPEN trips back to OPEN immediately

## Tradeoffs
- **Pros** — Prevents cascading failures, self-healing, clear state model
- **Cons** — Adds latency from timeout; tuning thresholds is application-specific

## Language Notes
- **Go** — `sync.Mutex` + `time.Now`; or use `github.com/sony/gobreaker`
- **Java** — `java.util.concurrent.locks.ReentrantLock`; or Netflix Hystrix / Resiliance4j

## Related Patterns
- [retry](../retry) — often paired: retry before tripping the breaker
- [exponential_backoff](../exponential_backoff) — used for half-open probe timing
