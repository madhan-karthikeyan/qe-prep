package com.qe.test.cache;

import static org.junit.jupiter.api.Assertions.*;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.atomic.AtomicInteger;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

@DisplayName("LRUCache")
class LRUCacheTest {

    @Test
    @DisplayName("get returns value after put")
    void getAfterPut() {
        var cache = new LRUCache<String, Integer>(3);
        cache.put("a", 1);
        assertEquals(1, cache.get("a"));
    }

    @Test
    @DisplayName("get returns null for missing key")
    void getMissing() {
        var cache = new LRUCache<String, Integer>(3);
        assertNull(cache.get("missing"));
    }

    @Test
    @DisplayName("evicts least recently used when at capacity")
    void eviction() {
        var cache = new LRUCache<String, Integer>(3);
        cache.put("a", 1);
        cache.put("b", 2);
        cache.put("c", 3);
        cache.put("d", 4);
        assertNull(cache.get("a"));
        assertNotNull(cache.get("b"));
        assertNotNull(cache.get("c"));
        assertNotNull(cache.get("d"));
    }

    @Test
    @DisplayName("access promotes item to most recent")
    void accessPromotes() {
        var cache = new LRUCache<String, Integer>(3);
        cache.put("a", 1);
        cache.put("b", 2);
        cache.put("c", 3);
        cache.get("a"); // promotes 'a'
        cache.put("d", 4);
        assertNotNull(cache.get("a"));
        assertNull(cache.get("b")); // 'b' should be evicted
    }

    @Test
    @DisplayName("update existing key does not change size")
    void updateExisting() {
        var cache = new LRUCache<String, Integer>(3);
        cache.put("a", 1);
        cache.put("a", 99);
        assertEquals(1, cache.size());
        assertEquals(99, cache.get("a"));
    }

    @Test
    @DisplayName("capacity must be positive")
    void capacityValidation() {
        assertThrows(IllegalArgumentException.class, () -> new LRUCache<>(0));
        assertThrows(IllegalArgumentException.class, () -> new LRUCache<>(-1));
    }

    @Test
    @DisplayName("null key rejected")
    void nullKey() {
        var cache = new LRUCache<String, Integer>(3);
        assertThrows(NullPointerException.class, () -> cache.put(null, 1));
    }

    @Test
    @DisplayName("thread-safe mode handles concurrent access")
    void concurrentAccess() throws Exception {
        int threadCount = 10;
        int opsPerThread = 100;
        var cache = new LRUCache<Integer, Integer>(50, true);
        var latch = new CountDownLatch(threadCount);
        var executor = Executors.newFixedThreadPool(threadCount);

        for (int t = 0; t < threadCount; t++) {
            int base = t * opsPerThread;
            executor.submit(() -> {
                try {
                    for (int i = 0; i < opsPerThread; i++) {
                        cache.put(base + i, base + i);
                        cache.get(base + i);
                    }
                } finally {
                    latch.countDown();
                }
            });
        }
        latch.await();
        executor.shutdown();
        assertTrue(cache.size() <= 50);
    }

    @Test
    @DisplayName("forEach iterates entries")
    void forEachTest() {
        var cache = new LRUCache<String, Integer>(3);
        cache.put("a", 1);
        cache.put("b", 2);
        var count = new AtomicInteger(0);
        cache.forEach((k, v) -> count.incrementAndGet());
        assertEquals(2, count.get());
    }
}
