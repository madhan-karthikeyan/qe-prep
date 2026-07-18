# Networking
Difficulty: Medium
Estimated Interview Time: 45 min
Prerequisites: TCP/IP, Java sockets, concurrency, HTTP basics

## Problem Statement
Implement a TCP echo server/client and supporting networking utilities including an HTTP client with retry logic and a generic connection pool.

## Requirements
- TcpEchoServer: thread-per-connection using virtual threads, graceful shutdown
- TcpEchoClient: connect, send/receive, timeout support
- HttpClient: retry with exponential backoff + jitter
- ConnectionPool: acquire/release, max connections, timeout

## Implementation Notes
- Uses Thread.ofVirtual() for lightweight server connections
- Connection pool uses BlockingQueue with bounded capacity
- HttpClient uses java.net.HttpURLConnection with configurable retry policy
- All components implement AutoCloseable

## Test Strategy
- Unit: validation, state checks
- Integration: start server, connect client, echo round-trip
- Stress: concurrent connection pool access

## Edge Cases
- Invalid port numbers (0-65535 required)
- Null/blank host names
- Connection pool exhaustion with timeout
- HTTP request failures with retry exhaustion

## Failure Cases
- Connection refused → IOException
- Pool timeout → IllegalStateException
- Invalid parameters → IllegalArgumentException

## Complexity (Time + Space)
- Echo server: O(1) per connection accept, O(n) threads for n clients
- ConnectionPool: O(1) acquire/release, O(maxConnections) space

## Progression Path
Build echo server first, then client. Add connection pool and HTTP client with retry logic.

## Common Interview Follow-ups
- Non-blocking I/O with Selector
- HTTP/2 support
- Connection health checking and eviction

## Possible Production Improvements
- Use java.nio.channels for scalable I/O
- Add TLS/SSL support
- Implement connection keep-alive and health checks
