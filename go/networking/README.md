# Networking
Difficulty: Medium
Estimated Interview Time: 40 min
Prerequisites: net, net/http, context, goroutines

## Problem Statement
Implement a TCP echo server/client, an HTTP client with retry, and a connection
pool.

## Requirements
- TCP echo with concurrent connections and graceful shutdown
- HTTP client with exponential backoff and jitter
- Thread-safe connection pool

## Implementation Notes
- TCP server per-connection goroutines
- Retry client retries on 5xx responses and transport errors
- Conn pool reuses idle connections

## Test Strategy
- Unit tests per component
- Integration: TCP echo round-trip
- Stress: 50 concurrent connections through pool

## Edge Cases
- Server shutdown mid-connection
- Pool exhaustion
- Timeout on idle connections

## Failure Cases
- Connection refused
- Pool at capacity
- TCP write errors

## Complexity (Time + Space)
- TCP echo: O(n) per message, O(1) per connection
- Retry client: O(retries) time
- Conn pool: O(1) acquire/release

## Progression Path
- Add TLS support
- Add connection health checks
- Add circuit breaker pattern

## Common Interview Follow-ups
- Graceful shutdown patterns
- Connection leak prevention
- Backoff strategy comparison

## Possible Production Improvements
- Use exponential backoff with full jitter
- Add metrics and tracing
- Integrate with service discovery
