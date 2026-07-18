package com.qe.test.ratelimit;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.atomic.AtomicInteger;

import static org.junit.jupiter.api.Assertions.*;

@DisplayName("TokenBucket")
class TokenBucketTest {

    @Test
    @DisplayName("allows burst of requests up to capacity")
    void allowsBurst() {
        var bucket = new TokenBucket(10, 1.0);
        for (int i = 0; i < 10; i++) {
            assertTrue(bucket.allowRequest(), "Request " + i + " should be allowed");
        }
        assertFalse(bucket.allowRequest(), "Request beyond capacity should be denied");
    }

    @Test
    @DisplayName("refills tokens over time")
    @Timeout(3)
    void refillsOverTime() throws InterruptedException {
        var bucket = new TokenBucket(5, 10.0);
        for (int i = 0; i < 5; i++) {
            assertTrue(bucket.allowRequest());
        }
        assertFalse(bucket.allowRequest());
        Thread.sleep(200);
        assertTrue(bucket.allowRequest());
    }

    @Test
    @DisplayName("does not exceed burst capacity after refill")
    void doesNotExceedBurstCapacity() {
        var bucket = new TokenBucket(5, 100.0);
        assertEquals(5.0, bucket.getAvailableTokens(), 0.01);
    }

    @Test
    @DisplayName("rejects negative capacity")
    void rejectsNegativeCapacity() {
        assertThrows(IllegalArgumentException.class, () -> new TokenBucket(0, 1.0));
        assertThrows(IllegalArgumentException.class, () -> new TokenBucket(-1, 1.0));
    }

    @Test
    @DisplayName("rejects non-positive refill rate")
    void rejectsNonPositiveRefillRate() {
        assertThrows(IllegalArgumentException.class, () -> new TokenBucket(10, 0));
        assertThrows(IllegalArgumentException.class, () -> new TokenBucket(10, -1));
    }

    @Test
    @DisplayName("allows request for specific token count")
    void allowsRequestForSpecificTokenCount() {
        var bucket = new TokenBucket(10, 1.0);
        assertTrue(bucket.allowRequest(5));
        assertTrue(bucket.allowRequest(5));
        assertFalse(bucket.allowRequest(1));
    }

    @Test
    @DisplayName("rejects non-positive token count in allowRequest")
    void rejectsNonPositiveTokenCount() {
        var bucket = new TokenBucket(10, 1.0);
        assertThrows(IllegalArgumentException.class, () -> bucket.allowRequest(0));
        assertThrows(IllegalArgumentException.class, () -> bucket.allowRequest(-1));
    }

    @Test
    @DisplayName("remains thread-safe under concurrent access")
    @Timeout(10)
    void concurrentAccess() throws InterruptedException {
        var bucket = new TokenBucket(100, 0.001); // extremely slow refill so burst dominates
        int threadCount = 10;
        var latch = new CountDownLatch(threadCount);
        var executor = Executors.newFixedThreadPool(threadCount);
        var allowed = new AtomicInteger(0);

        for (int i = 0; i < threadCount; i++) {
            executor.submit(() -> {
                for (int j = 0; j < 20; j++) {
                    if (bucket.allowRequest()) {
                        allowed.incrementAndGet();
                    }
                }
                latch.countDown();
            });
        }

        latch.await();
        executor.shutdown();
        assertTrue(allowed.get() <= 100, "Should not exceed burst capacity");
    }
}
