# Rate Limiter
Difficulty: Medium
Estimated Interview Time: 30 min
Prerequisites: goroutines, sync.Mutex, time

## Problem Statement
Implement a token bucket and sliding window log rate limiter.

## Requirements
- Token bucket with burst capacity and refill rate
- Sliding window log with max requests per window
- Both thread-safe via sync.Mutex

## Implementation Notes
- Token bucket uses continuous time-based refill
- Sliding window evicts expired entries on each Allow call
- Both support concurrent access

## Test Strategy
- Unit tests for burst/refill limits
- Stress tests with 100 concurrent goroutines

## Edge Cases
- Empty bucket at start
- Rapid concurrent access
- Window boundary conditions

## Failure Cases
- Requests denied when over limit
- Token exhaustion

## Complexity (Time + Space)
- Token bucket: O(1) time, O(1) space per Allow call
- Sliding window: O(n) time (n = window size), O(n) space

## Progression Path
- Add AllowN for bulk operations
- Add distributed rate limiter with Redis

## Common Interview Follow-ups
- Distributed rate limiting
- Adaptive rate limiting
- GCRA algorithm

## Possible Production Improvements
- Use sync.Pool for timestamp slices
- Add metrics collection
- Support hierarchical rate limits
