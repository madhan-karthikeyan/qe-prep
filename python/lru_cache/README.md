# LRU Cache

Difficulty: Medium
Estimated Interview Time: 25 min
Prerequisites: Doubly-linked list, hash map, O(1) operations

## Problem Statement

Implement an LRU (Least Recently Used) cache with O(1) get and put operations using a dictionary combined with a manually-implemented doubly-linked list.

## Requirements

- O(1) get and put
- Dict + doubly-linked list (NOT collections.OrderedDict)
- Configurable capacity
- Thread-safe option via threading.Lock
- get(key) returns value or None
- put(key, value) evicts LRU item when at capacity
- Accessing an item makes it most recently used

## Implementation Notes

- _Node stores key, value, prev, next pointers
- Sentinel head/tail nodes simplify boundary conditions
- _remove_node and _add_to_front helpers keep code DRY
- Thread safety via optional lock (acquire/release around operations)

## Test Strategy
- Unit: get missing, put/get, eviction, update existing, access renewal, None values, capacity=1, len, invalid capacity
- Thread safety: 10 threads doing 1000 ops each on capacity=100 cache

## Edge Cases

- Capacity of 1 ensures immediate eviction on second put
- None as a stored value (distinct from "missing")
- Repeatedly accessing the same item keeps it alive
- Putting an existing key updates its value and moves it to front

## Failure Cases

- Capacity < 1 (raises ValueError)
- Concurrent modification without thread_safe=True (data race possible)
- Storing unhashable keys (raises TypeError from dict)

## Complexity
- Time: O(1) get and put
- Space: O(capacity)

## Progression Path
Basic → Thread-safe → TTL-based expiry → Distributed LRU

## Common Interview Follow-ups

- How would you add TTL (time-to-live) expiry?
- How would you implement an LFU cache?
- How would you make this work across multiple processes?
- What if we needed O(1) contains check?

## Possible Production Improvements

- TTL-based automatic eviction of stale entries
- Memory-aware capacity (limit by byte size rather than count)
- Weak reference support
- Statistics (hit rate, miss rate)
