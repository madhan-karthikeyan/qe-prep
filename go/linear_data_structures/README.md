# Linear Data Structures
Difficulty: Easy
Estimated Interview Time: 20 min
Prerequisites: slices, linked lists, generics

## Problem Statement
Implement generic stack (LIFO), queue (FIFO), and circular queue data structures.

## Requirements
- Stack: Push, Pop, Peek, IsEmpty, Size — slice-backed
- Queue: Enqueue, Dequeue, Peek, IsEmpty, Size — container/list-backed
- CircularQueue: Enqueue, Dequeue, Peek, IsFull, IsEmpty, Size — fixed buffer
- All generic (any type)

## Implementation Notes
- Stack uses Go slice for O(1) amortized push/pop
- Queue uses container/list for O(1) enqueue/dequeue
- CircularQueue uses a fixed-size array with head/tail pointers
- All types use Go generics (go 1.18+)

## Test Strategy
- Table-driven tests for each operation
- Edge cases: empty pop, empty peek, full enqueue
- Wrap-around behavior for circular queue

## Edge Cases
- Pop from empty stack/queue
- Peek from empty structure
- Enqueue to full circular buffer
- Capacity 1 circular queue

## Failure Cases
- N/A (all operations return bool for failure)

## Complexity (Time + Space)
- Stack: all O(1) amortized; space O(n)
- Queue: all O(1); space O(n)
- CircularQueue: all O(1); space O(capacity)

## Progression Path (Basic → Intermediate → Advanced → Production)
- Basic: Single-type implementations
- Intermediate: Generic (any type) support
- Advanced: Thread-safe variants
- Production: Lock-free queues, ring buffers

## Common Interview Follow-ups
- When would you use a slice vs linked list for a queue?
- How would you make these thread-safe?
- What are the trade-offs of circular buffers?

## Possible Production Improvements
- Add mutex for thread-safe variants
- Implement a lock-free queue using atomics
- Add resize capability to circular queue
