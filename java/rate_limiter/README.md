# Rate Limiter
Difficulty: Medium
Estimated Interview Time: 30 min
Prerequisites: Concurrency, ReentrantLock, Token Bucket + Sliding Window algorithms

## Problem Statement
Implement two rate-limiting algorithms: Token Bucket and Sliding Window Log. Both must be thread-safe and support concurrent access.

## Requirements
- TokenBucket: burst capacity, refill rate, allowRequest()
- SlidingWindowLog: fixed window, max requests, timestamp-based eviction
- Thread-safe using ReentrantLock
- Configurable parameters with validation

## Implementation Notes
- Uses System.nanoTime() for high-resolution timing
- ReentrantLock for fine-grained locking (not synchronized)
- TokenBucket refills lazily on request
- SlidingWindowLog prunes expired timestamps on each check

## Test Strategy
- Unit tests for burst, refill, edge cases
- Concurrent stress tests with CountDownLatch + ExecutorService
- @Timeout on blocking tests

## Edge Cases
- Zero/negative burst capacity
- Zero/negative refill rate
- Token count overflow after long idle period
- Rapid concurrent requests near window boundary

## Failure Cases
- Construction with invalid parameters → IllegalArgumentException
- Requests beyond capacity → deny (return false)

## Complexity (Time + Space)
- TokenBucket: O(1) per request, O(1) space
- SlidingWindowLog: O(W) per request worst-case (pruning), O(W) space where W = window size

## Progression Path
Start with TokenBucket, then implement SlidingWindowLog. Add concurrent stress tests last.

## Common Interview Follow-ups
- Distributed rate limiting with Redis
- Adaptive rate limiting based on system load
- Multi-tier rate limits (per-user + global)

## Possible Production Improvements
- Switch to a background thread for refill
- Use LongAdder for lock-free counting
- Add metrics collection for monitoring
