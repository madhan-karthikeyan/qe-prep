# Pub/Sub

## Intent
Decouple message producers from consumers via topic-based broadcast.

## Structure
- `PubSub` — central broker with topic routing
- `subscribe(topic, callback)` — returns subscription ID
- `unsubscribe(subscription_id)` — removes subscriber
- `publish(topic, message)` — async delivery to all subscribers

## When to Use
- Event-driven architectures, notification systems
- Decoupling components that should not know about each other

## When NOT to Use
- When message ordering must be guaranteed
- When exactly-once delivery is required
- When back-pressure is needed (use a queue-based pattern)

## Implementation Notes
- Async delivery via `ThreadPoolExecutor` — non-blocking publish
- Subscriptions stored per-topic in a dict of dicts
- Thread-safe using `threading.Lock` for subscribe/unsubscribe/publish

## Tradeoffs
- **Pros** — Loose coupling, easy to extend, async by default
- **Cons** — No delivery guarantees, no persistence, no ordering

## Language Notes
- **Go** — Idiomatic: channels + goroutines; `sync.RWMutex` for subscriber map
- **Java** — `java.util.concurrent.ConcurrentHashMap`; `ExecutorService` for async dispatch

## Related Patterns
- [producer_consumer](../producer_consumer) — queue-based point-to-point vs topic-based broadcast
- [rate_limiter](../rate_limiter) — protect subscribers from publish storms
