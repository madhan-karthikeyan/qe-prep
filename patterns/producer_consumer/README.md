# Producer-Consumer

## Intent
Decouple work generation from work processing using a thread-safe queue.

## Structure
- `ProducerConsumer` — orchestrates N producers and M consumers
- Poison pill sent per consumer to signal graceful shutdown

## When to Use
- Streaming data processing, log ingestion, job queues
- When producers and consumers have different throughput

## When NOT to Use
- Tightly coupled request-response flows
- When ordering must be guaranteed across partitions

## Implementation Notes
- Poison-pill pattern: one sentinel per consumer, placed after all producers finish
- Producers and consumers are daemon threads
- Exceptions in producers/consumers are caught to avoid hangs

## Tradeoffs
- **Pros** — Simple decoupling, configurable concurrency, graceful shutdown
- **Cons** — Single queue can become a bottleneck; no back-pressure beyond `maxsize`

## Language Notes
- **Go** — Idiomatic: goroutines + channels; `close(ch)` acts as broadcast poison pill
- **Java** — `BlockingQueue` + `ExecutorService`; `LinkedBlockingQueue` for bounded queues

## Related Patterns
- [worker_pool](../worker_pool) — fixed worker pool for task execution
- [pub_sub](../pub_sub) — one-to-many message delivery vs queue-based work distribution
