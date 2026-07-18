# Producer-Consumer
Difficulty: Medium
Estimated Interview Time: 45 min
Prerequisites: ReentrantLock, Condition, concurrency

## Problem Statement
Implement a bounded blocking queue and a producer-consumer demo.

## Requirements
- Array-based blocking queue with put/take/offer/poll
- ReentrantLock + Condition for synchronization
- Multiple producers and consumers

## Implementation Notes
- Ring buffer with ReentrantLock and two Conditions (notEmpty, notFull)
- offer/poll with timeout support
- ProducerConsumerDemo orchestrates N producers and N consumers

## Test Strategy
- FIFO ordering
- Blocking behavior verification
- Timeout variants
- Stress test: 10 producers + 10 consumers, 100 items each

## Edge Cases
- Full queue blocking
- Empty queue blocking
- Interrupted thread handling

## Failure Cases
- InterruptedException propagation
- Invalid capacity

## Complexity
- Time: O(1) put/take
- Space: O(capacity)

## Progression Path
1. Simple blocking queue → 2. Timeout variants → 3. Multiple producers/consumers

## Common Interview Follow-ups
- What happens if a thread is interrupted while waiting?
- How would you implement a priority blocking queue?
- Difference between wait/notify and Lock/Condition?

## Possible Production Improvements
- LinkedBlockingQueue-style dynamic sizing
- TransferQueue semantics
- Backpressure with metrics
