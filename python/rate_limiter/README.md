# Rate Limiter
Difficulty: Medium
Estimated Interview Time: 30 min
Prerequisites: threading, time-series algorithms

## Problem Statement
Implement two rate-limiting strategies: a token bucket and a sliding window log, both thread-safe.

## Requirements
- Token bucket with configurable burst capacity, refill rate, and refill interval
- Sliding window log with configurable window size and max requests
- Both must be thread-safe (Lock)
- `allow_request()` returns bool
- Token bucket supports consuming multiple tokens at once

## Implementation Notes
- Uses `time.monotonic()` for wall-clock-independent timing
- Token bucket refill is lazy — calculated on demand
- Sliding window uses `collections.deque` for O(1) pops from left

## Test Strategy (Unit/Stress)
- Unit: burst then exhaust, refill after interval, multi-token consumption, invalid args
- Stress: 100 concurrent requesters verifying rate enforcement

## Edge Cases
- Zero or negative constructor parameters
- Token bucket: requesting zero tokens
- Sliding window: timestamps that have fallen out of the window

## Failure Cases
- Exhausted capacity → return False
- Concurrent access causing race conditions (mitigated by Lock)

## Complexity
- Token Bucket: O(1) per request, O(1) space
- Sliding Window: O(n) per request in worst case (pruning), O(w) space where w = max requests

## Progression Path
- Basic: single-threaded token bucket
- Intermediate: thread-safe sliding window
- Advanced: distributed rate limiter (Redis)
- Production: multi-tier rate limiting, adaptive rate limiting

## Common Interview Follow-ups
- How would you make this distributed?
- How would you handle clock skew?
- How would you implement rate limiting per user?

## Possible Production Improvements
- Use Redis or Memcached for distributed state
- Add metrics/monitoring
- Support hierarchical rate limits (global + per-user)
