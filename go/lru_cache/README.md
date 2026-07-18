# LRU Cache
Difficulty: Medium
Estimated Interview Time: 25 min
Prerequisites: hash map, doubly-linked list, O(1) operations

## Problem Statement
Implement an LRU (Least Recently Used) cache with O(1) Get and Put operations.

## Requirements
- map[int]*Node + doubly-linked list
- O(1) Get and Put
- Configurable capacity
- Thread-safe via mutex

## Implementation Notes
- Hash map provides O(1) key lookup
- Doubly-linked list maintains access order
- Move-to-front on access; evict-from-tail when full
- Mutex wraps all public methods for concurrent safety

## Test Strategy
- Basic Put/Get and miss
- Eviction order verification
- Update (re-Put) existing key
- Concurrent access with multiple goroutines
- Zero-capacity edge case

## Edge Cases
- Zero capacity (never stores anything)
- Duplicate Put (update existing)
- Get on empty cache

## Failure Cases
- Nil map access (handled by constructor)

## Complexity (Time + Space)
- Get: O(1) time, O(1) space
- Put: O(1) time, O(1) space (excluding eviction)
- Space: O(capacity)

## Progression Path (Basic → Intermediate → Advanced → Production)
- Basic: Fixed-size cache with eviction
- Intermediate: Generics support, TTL
- Advanced: LFU variant, segmented LRU
- Production: Sharded cache, persistence, metrics

## Common Interview Follow-ups
- How would you add a TTL (time-to-live)?
- How would you implement an LFU variant?
- How would you shard the cache for higher throughput?

## Possible Production Improvements
- Generify to support any key/value types
- Add TTL with background expiration
- Add metrics (hit rate, eviction count)
- Shard by key hash for concurrent throughput
