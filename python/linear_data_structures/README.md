# Linear Data Structures

Difficulty: Easy
Estimated Interview Time: 15 min
Prerequisites: Basic data structures

## Problem Statement

Implement generic Stack, Queue, and CircularQueue (fixed-size circular buffer) data structures using idiomatic Python.

## Requirements

Stack (list-based):
- push, pop, peek, is_empty, size

Queue (collections.deque-based):
- enqueue, dequeue, peek, is_empty, size

CircularQueue (fixed-size array-based):
- enqueue, dequeue, peek, is_full, is_empty, size

## Implementation Notes

- Stack uses list append/pop for O(1) amortized
- Queue uses deque popleft/append for O(1)
- CircularQueue uses modulo arithmetic for wrap-around
- All are generic via TypeVar

## Test Strategy
- Unit: push/pop, peek, empty/full, FIFO order, wrap-around (circular queue), error cases

## Edge Cases

- Pop/peek on empty Stack, Queue, CircularQueue (IndexError)
- Enqueue on full CircularQueue (IndexError)
- CircularQueue wrap-around after fill + partial drain + refill

## Failure Cases

- CircularQueue capacity < 1 (ValueError)
- All empty-structure operations raise IndexError consistently

## Complexity
- Stack: O(1) all operations
- Queue: O(1) all operations
- CircularQueue: O(1) all operations

## Progression Path
Basic → Bounded Queue → Thread-safe versions → Lock-free versions

## Common Interview Follow-ups

- How would you make these thread-safe?
- When would you use a circular buffer over a deque?
- How would you implement a min-stack?
- How would you implement a queue using two stacks?

## Possible Production Improvements

- Thread-safe variants with locks
- Lock-free circular buffer with atomics
- Bounded blocking queue (see producer_consumer module)
- Capacity-doubling dynamic circular buffer
