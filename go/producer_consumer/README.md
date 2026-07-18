# Producer-Consumer
Difficulty: Medium
Estimated Interview Time: 25 min
Prerequisites: channels, goroutines, context.Context

## Problem Statement
Implement a blocking queue and a producer-consumer demo with configurable producers and consumers.

## Requirements
- BlockingQueue: fixed capacity, Put blocks when full, Get blocks when empty
- Close() signals no more producers; consumers drain remaining items
- Context-based cancellation
- Demo with N producers and M consumers

## Implementation Notes
- BlockingQueue wraps a buffered channel internally
- A closed channel signals shutdown (separate from the data channel)
- Context enables timeout/cancellation for Put/Get
- Close is safe to call multiple times via careful design

## Test Strategy
- Basic Put/Get ordering (FIFO)
- Blocking behavior verification
- Context timeout
- Close and drain remaining items
- Multiple goroutines with stress test
- Race detection with -race flag

## Edge Cases
- Put to full queue blocks
- Get from empty queue blocks
- Close while producers are active
- Context cancelled during blocked Put/Get

## Failure Cases
- Put/Get after close returns ErrQueueClosed
- Context cancellation returns context.Canceled

## Complexity (Time + Space)
- Put/Get: O(1) amortized (channel operations)
- Space: O(capacity)

## Progression Path (Basic → Intermediate → Advanced → Production)
- Basic: Single producer, single consumer
- Intermediate: Multiple producers/consumers
- Advanced: Context support, graceful shutdown
- Production: Priority queue, rate limiting, backpressure signals

## Common Interview Follow-ups
- How would you implement a priority blocking queue?
- How would you add rate limiting?
- How would you handle slow consumers?

## Possible Production Improvements
- Add metrics (queue depth, throughput, latency)
- Implement work stealing
- Add backpressure signaling to producers
- Support multiple queues with routing
