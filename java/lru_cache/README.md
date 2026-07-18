# LRU Cache
Difficulty: Medium
Estimated Interview Time: 45 min
Prerequisites: HashMap, doubly-linked list

## Problem Statement
Implement a generic LRU (Least Recently Used) cache with O(1) get/put.

## Requirements
- Generic types <K, V>
- O(1) get and put
- Configurable capacity
- Optional thread-safe mode

## Implementation Notes
- HashMap + doubly-linked list (manual implementation)
- Head/tail sentinel nodes
- ReentrantLock for thread-safe mode

## Test Strategy
- Capacity eviction
- Access-order promotion
- Update existing key preserves size
- Concurrent access test with CountDownLatch

## Edge Cases
- Cache size 1
- Access after eviction
- Update of existing key

## Failure Cases
- Invalid capacity (non-positive)
- Null key

## Complexity
- Time: O(1) get/put
- Space: O(capacity)

## Progression Path
1. Simple HashMap cache → 2. Bounded with eviction → 3. Access-order

## Common Interview Follow-ups
- How would you add TTL expiry?
- What if you need a max memory-based eviction?
- How to implement LFU instead?

## Possible Production Improvements
- TTL with background sweeper thread
- Weak/soft references for values
- JMX monitoring
