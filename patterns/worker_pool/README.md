# Worker Pool

## Intent
Manage a fixed set of worker threads to execute tasks concurrently from a shared queue.

## Structure
- `WorkerPool` — manages N daemon worker threads
- `submit(task)` — non-blocking task enqueue
- `map(func, iterable)` — distribute work items
- `shutdown(wait)` — graceful shutdown with optional drain

## When to Use
- CPU-bound or I/O-bound work that can be parallelized
- Backend task execution (e.g., request processing, batch jobs)

## When NOT to Use
- When you need dynamic scaling (use producer-consumer with variable workers)
- For async I/O (consider asyncio instead)

## Implementation Notes
- Workers poll with timeout for responsive shutdown
- Poison pill pattern: one sentinel per worker
- `result_callback` receives the return value (or exception) per task

## Tradeoffs
- **Pros** — Fixed resource usage, clean lifecycle, easy to reason about
- **Cons** — No dynamic scaling; blocked workers reduce throughput

## Language Notes
- **Go** — `sync.WaitGroup` + goroutines + channels; no poison pill needed with `close(ch)`
- **Java** — `ThreadPoolExecutor` with `LinkedBlockingQueue`; `submit(Callable)` returns `Future`

## Related Patterns
- [producer_consumer](../producer_consumer) — general queue decoupling
- [object_pool](../object_pool) — pools resources, not workers
