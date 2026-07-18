# Networking
Difficulty: Medium
Estimated Interview Time: 35 min
Prerequisites: socket programming, threading, HTTP

## Problem Statement
Implement a TCP echo server/client, an HTTP client with retry, and a connection pool.

## Requirements
- TCP echo server handles multiple concurrent clients with threading
- TCP echo client connects, sends, and receives
- HTTP client with configurable retries, exponential backoff, jitter, and timeout
- Connection pool with acquire/release, max connections, timeout, and health check

## Implementation Notes
- Uses `socket` for TCP, `http.client` for HTTP, `threading` for concurrency
- Connection pool evicts stale idle connections
- HTTP client retries on 429 and 5xx status codes and connection errors

## Test Strategy (Unit/Integration/Stress)
- Unit: echo round-trip, binary data, large payloads; HTTP client construction; pool acquire/release/reuse
- Integration: echo server + client end-to-end, pool echo round-trip
- Stress: 50 concurrent connections through pool

## Edge Cases
- Empty payload, binary data, large payloads
- Connection pool exhaustion
- Server shutdown during active connections
- Timeout handling

## Failure Cases
- Connection refused / timeout
- Pool capacity exceeded
- Unhealthy connections detected by health check

## Complexity
- TCP Echo: O(1) per client, O(n) threads
- HTTP Client: O(r) per request (r = retries), O(1) space
- Connection Pool: O(1) acquire/release, O(s) eviction (s = stale connections)

## Progression Path
- Basic: single-client echo server
- Intermediate: multi-client server with threads
- Advanced: HTTP client with exponential backoff
- Production: async I/O, connection pooling with keep-alive

## Common Interview Follow-ups
- How would you implement async versions?
- How would you detect dead connections?
- How would you implement circuit breaker?

## Possible Production Improvements
- Use asyncio for better scalability
- Add TLS/SSL support
- Implement circuit breaker pattern
- Add metrics for pool utilization
