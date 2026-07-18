# Thread Pool
Difficulty: Hard
Estimated Interview Time: 50 min
Prerequisites: Threads, BlockingQueue, CompletableFuture

## Problem Statement
Implement a custom thread pool without using Executors.

## Requirements
- Fixed number of worker threads
- Task queue (BlockingQueue)
- submit(Callable) -> Future via CompletableFuture
- shutdown() and awaitTermination()

## Implementation Notes
- Workers poll from a LinkedBlockingQueue
- Tasks wrap Callable + CompletableFuture
- Poll with timeout enables graceful shutdown

## Test Strategy
- Submit single/multiple tasks, verify results
- Exception propagation
- Shutdown behavior (rejects, waits)
- Stress test: 1000 tasks with 8 workers

## Edge Cases
- Shutdown with pending tasks
- Interrupted worker threads
- Task exception handling

## Failure Cases
- IllegalStateException after shutdown
- ExecutionException from failed tasks

## Complexity
- Time: O(1) submit
- Space: O(numWorkers + queueSize)

## Progression Path
1. Single-threaded executor → 2. Fixed pool → 3. Shutdown hooks → 4. Dynamic scaling

## Common Interview Follow-ups
- How would you implement a dynamically resizing pool?
- What is the difference with Executors.newFixedThreadPool?
- How would you add a rejection policy?

## Possible Production Improvements
- Core/max pool size with scaling
- Keep-alive and idle timeout
- Rejection policies (abort, discard, caller-runs)
- Prestart core threads
