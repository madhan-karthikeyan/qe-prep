# Producer-Consumer

Difficulty: Medium
Estimated Interview Time: 30 min
Prerequisites: Threading, synchronization, Condition variables

## Problem Statement

Implement a thread-safe bounded blocking queue and a producer-consumer demo with multiple producers and consumers.

## Requirements

BlockingQueue:
- Fixed capacity
- Blocking put when full, blocking get when empty
- Optional timeout on put/get
- Thread-safe using Condition variables

Producer-Consumer Demo:
- Multiple producers and consumers running concurrently
- All items produced are eventually consumed
- No deadlock or data loss

## Implementation Notes

- BlockingQueue uses two Condition variables (_not_full, _not_empty) on a shared lock
- put() waits on _not_full, signals _not_empty
- get() waits on _not_empty, signals _not_full
- Timeout returns False (put) or None (get) on expiry

## Test Strategy
- Unit: put/get, timeout on full/empty, FIFO ordering
- Stress: 10 producers + 10 consumers, 100k items, verify all items produced and consumed

## Edge Cases

- Put on full queue blocks until consumer drains
- Get on empty queue blocks until producer adds
- Timeout expiry returns correct sentinel values
- Multiple consumers can contend for same item

## Failure Cases

- Capacity < 1 (ValueError)
- Deadlock if notify is not called (avoided by using Condition)
- Spurious wakeup (handled by while-loop waiting pattern in Condition)

## Complexity
- Time: O(1) per put/get (amortized, ignoring contention)
- Space: O(capacity)

## Progression Path
Basic → Timeout → Multiple producers/consumers → Zero-copy ring buffer

## Common Interview Follow-ups

- How would you implement a zero-copy ring buffer?
- What happens with spurious wakeups?
- How would you handle priority ordering?
- How is this different from queue.Queue?

## Possible Production Improvements

- Use a lock-free ring buffer for lower contention
- Support batch put/get for throughput
- Priority queue variant
- Backpressure signaling to producers
- Metrics (queue depth, wait times)
