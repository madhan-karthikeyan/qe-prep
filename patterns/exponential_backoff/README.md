# Exponential Backoff

## Intent
Calculate delay values for retry strategies using exponential backoff with pluggable jitter.

## Structure
- `BackoffCalculator` — computes delay(attempt) with configurable params
- `full_jitter` — random in [0, delay)
- `equal_jitter` — random around delay/2
- `decorrelated_jitter` — smooth backoff using previous delay

## When to Use
- As a building block for retry, circuit breaker, or rate-limiting logic
- When you need consistent backoff math across a system

## When NOT to Use
- When you need the full retry orchestrator (use the [retry](../retry) pattern)

## Implementation Notes
- Three jitter strategies map to common production patterns (AWS, Google, Netflix)
- `decorrelated_jitter` requires tracking previous_delay; exposed as separate function
- All delays are capped at `max_delay`

## Tradeoffs
- **Pros** — Composable, testable in isolation, no side effects
- **Cons** — Single-purpose; callers must manage loop/state

## Language Notes
- **Go** — `time.Duration` math; `math/rand` for jitter; can be a pure function
- **Java** — `Duration` API; `ThreadLocalRandom` for thread-safe jitter

## Related Patterns
- [retry](../retry) — orchestrates retries using backoff delays
- [circuit_breaker](../circuit_breaker) — uses backoff for half-open timeout
