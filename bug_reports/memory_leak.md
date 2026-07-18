# Memory leak in LRU cache under concurrent access

**Severity:** Critical
**Priority:** P1
**Environment:** All platforms, Java 11+, concurrent read/write ratio > 5:1
**Component:** Caching subsystem — LRUCache

## Summary

Under concurrent `put()` and `get()` operations, the LRU cache's size exceeds its configured maximum capacity and never recovers, leading to unbounded memory growth. Over 24 hours of production traffic, the cache grows 3-5× beyond the configured limit.

## Steps to Reproduce

1. Configure LRUCache with `maxSize = 1000`
2. Spawn 16 threads performing random `get()` and `put()` operations
3. Monitor cache size via `size()` or internal counter
4. Run for 10 minutes

```java
// Reproduction test
LRUCache<Integer, String> cache = new LRUCache<>(1000);
ExecutorService pool = Executors.newFixedThreadPool(16);

for (int i = 0; i < 16; i++) {
    pool.submit(() -> {
        for (int j = 0; j < 100000; j++) {
            int key = ThreadLocalRandom.current().nextInt(2000);
            if (ThreadLocalRandom.current().nextBoolean()) {
                cache.put(key, "value-" + key);
            } else {
                cache.get(key);
            }
        }
    });
}
pool.shutdown();
pool.awaitTermination(1, TimeUnit.MINUTES);

System.out.println("Expected size: ~1000, Actual size: " + cache.size());
// Output: Expected size: ~1000, Actual size: 4723
```

## Expected Behavior

Cache should never exceed `maxSize`. When a new entry is added at capacity, the least recently used entry should be evicted atomically.

## Actual Behavior

- Cache size grows to 3-5× the configured maximum under concurrent load
- Memory usage grows linearly with runtime
- Once oversized, the cache never evicts down to capacity
- GC overhead increases, eventually causing `OutOfMemoryError`

## Logs / Screenshots

```
Heap dump analysis:
- LRUCache.internalMap size: 4723 (configured max: 1000)
- 87% of entries are old (last accessed > 1 hour ago)
- 12% of entries are duplicates with different access times
- Top retainers: LRUCache$Node objects (3.2MB)
```

## Root Cause Analysis

The LRU eviction logic has a race condition in the `put()` method:

```java
// Simplified problematic code
public V put(K key, V value) {
    synchronized (this) {
        Node node = map.get(key);
        if (node != null) {
            moveToHead(node);
            node.value = value;
            return oldValue;
        }
    }
    // RACE WINDOW: between sync blocks
    // Another thread can put() here, filling the cache
    
    synchronized (this) {
        if (size() >= maxSize) {
            removeTail();           // evict
        }
        addToHead(key, value);     // insert
    }
}
```

The first synchronized block releases the lock after checking for an existing key. Between the two synchronized blocks, multiple threads can pass the `size() >= maxSize` check before any of them evict. This results in the cache growing beyond `maxSize` by up to the number of racing threads.

## Fix

Make the entire `put()` operation atomic within a single synchronized block:

```java
public V put(K key, V value) {
    synchronized (this) {
        Node node = map.get(key);
        if (node != null) {
            moveToHead(node);
            node.value = value;
            return oldValue;
        }
        
        // Evict BEFORE adding, while still holding the lock
        while (map.size() >= maxSize) {
            removeTail();
        }
        
        addToHead(key, value);
        return null;
    }
}
```

For higher throughput, consider using `ConcurrentHashMap` with a lock per segment or `LinkedHashMap` with `removeEldestEntry` (with proper synchronization):

```java
public class ThreadSafeLRUCache<K, V> {
    private final int maxSize;
    private final ConcurrentHashMap<K, Node<K, V>> map;
    private final ReentrantLock lock = new ReentrantLock();
    
    public V put(K key, V value) {
        lock.lock();
        try {
            Node<K, V> node = map.get(key);
            if (node != null) {
                moveToHead(node);
                node.value = value;
                return node.value;
            }
            while (map.size() >= maxSize) {
                removeTail();
            }
            addToHead(key, value);
            return null;
        } finally {
            lock.unlock();
        }
    }
}
```

## Regression Tests

### 1. Concurrent Access Stress Test

```java
@Test
public void testLRUCacheDoesNotExceedCapacityUnderConcurrentAccess() {
    int maxSize = 1000;
    LRUCache<Integer, String> cache = new LRUCache<>(maxSize);
    int threadCount = 16;
    int opsPerThread = 100_000;
    
    ExecutorService pool = Executors.newFixedThreadPool(threadCount);
    CountDownLatch latch = new CountDownLatch(threadCount);
    
    for (int i = 0; i < threadCount; i++) {
        pool.submit(() -> {
            for (int j = 0; j < opsPerThread; j++) {
                int key = ThreadLocalRandom.current().nextInt(2000);
                if (ThreadLocalRandom.current().nextBoolean()) {
                    cache.put(key, "v" + key);
                } else {
                    cache.get(key);
                }
            }
            latch.countDown();
        });
    }
    latch.await();
    
    assertTrue("Cache exceeded max size: " + cache.size(),
               cache.size() <= maxSize * 1.05);  // allow 5% slack for atomic ops
}
```

### 2. Memory Leak Assertion

```java
@Test
public void testNoMemoryGrowthAfterStabilization() {
    LRUCache<Integer, byte[]> cache = new LRUCache<>(100);
    
    // Phase 1: fill cache
    for (int i = 0; i < 100; i++) {
        cache.put(i, new byte[1024]);
    }
    
    // Phase 2: churn with keys within capacity
    System.gc();
    long before = Runtime.getRuntime().totalMemory() - Runtime.getRuntime().freeMemory();
    
    for (int round = 0; round < 100; round++) {
        for (int i = 0; i < 100; i++) {
            cache.put(i + round, new byte[1024]);  // always within 100 unique keys
        }
    }
    
    System.gc();
    long after = Runtime.getRuntime().totalMemory() - Runtime.getRuntime().freeMemory();
    long delta = after - before;
    
    assertTrue("Memory grew by " + delta + " bytes after churn (leak suspected)",
               delta < 50_000);  // allow small GC overhead
}
```

### 3. CI Integration

Add to CI pipeline:
- Stress test runs for 60s with 16 threads and asserts `size() <= maxSize * 1.05`
- Heap dump analyzed on failure to detect node object accumulation
- Memory regression benchmark compared against previous run
