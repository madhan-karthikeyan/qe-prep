# Rate Limiter

## Intent
Control the rate of requests to protect downstream resources.

## Structure
- `RateLimiter` — token bucket implementation
- `allow_request()` — returns bool; consumes one token if available
- Refills tokens at `rate` per second, capped at `burst`

## When to Use
- API gateways, database connection throttling, external service calls
- Preventing resource exhaustion under load spikes

## When NOT to Use
- When you need distributed rate limiting across processes (use Redis)
- When you need per-user/per-IP rate limiting (extend the key concept)

## Implementation Notes
- Token bucket with lazy refill (computed on `allow_request`)
- Thread-safe via `threading.Lock`
- Uses `time.monotonic()` for wall-clock insensitive timing

## Tradeoffs
- **Pros** — Simple, efficient, allows bursts up to capacity
- **Cons** — Single-process only; burst size must be tuned

## Language Notes
- **Go** — `golang.org/x/time/rate` (official token bucket); `sync.Mutex` for custom
- **Java** — `Guava RateLimiter` (SmoothBursty / SmoothWarmingUp); `Semaphore` for simple cases

## Related Patterns
- [circuit_breaker](../circuit_breaker) — complementary: rate limiter throttles, breaker stops
- [retry](../retry) — rate limiter should be applied before retry to avoid amplification
