# Linear Data Structures
Difficulty: Easy
Estimated Interview Time: 25 min
Prerequisites: Java Collections

## Problem Statement
Implement generic Stack, Queue, and CircularQueue.

## Requirements
- Stack<E>: push, pop, peek (ArrayList-based)
- Queue<E>: enqueue, dequeue, peek (LinkedList-based)
- CircularQueue<E>: fixed-size ring buffer with wrap

## Implementation Notes
- Stack uses ArrayList for dynamic growth
- Queue uses LinkedList for O(1) head/tail operations
- CircularQueue uses Object[] ring buffer with modular arithmetic

## Test Strategy
- LIFO/FIFO ordering verification
- Empty state exceptions
- Full buffer behavior for CircularQueue
- Wrap-around test

## Edge Cases
- Empty pop/dequeue
- Full enqueue (CircularQueue)
- Null item rejection

## Failure Cases
- EmptyStackException / NoSuchElementException
- IllegalStateException on full queue

## Complexity
- Stack: O(1) push/pop/peek (amortized)
- Queue: O(1) enqueue/dequeue/peek
- CircularQueue: O(1) all operations

## Progression Path
1. Stack → 2. Queue → 3. CircularQueue → 4. Deque

## Common Interview Follow-ups
- Implement a deque with both ends
- How would you make these thread-safe?
- Dynamic resizing ring buffer?

## Possible Production Improvements
- ArrayDeque-based stack
- Thread-safe versions
- Bounded blocking queue variant
