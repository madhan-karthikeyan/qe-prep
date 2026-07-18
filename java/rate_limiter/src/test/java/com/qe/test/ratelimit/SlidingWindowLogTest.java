package com.qe.test.ratelimit;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.atomic.AtomicInteger;

import static org.junit.jupiter.api.Assertions.*;

@DisplayName("SlidingWindowLog")
class SlidingWindowLogTest {

    @Test
    @DisplayName("allows requests up to max within window")
    void allowsUpToMax() {
        var log = new SlidingWindowLog(1, SlidingWindowLog.TimeUnit.SECONDS, 5);
        for (int i = 0; i < 5; i++) {
            assertTrue(log.allowRequest(), "Request " + i + " should be allowed");
        }
        assertFalse(log.allowRequest(), "Request beyond max should be denied");
    }

    @Test
    @DisplayName("allows requests after window expires")
    @Timeout(3)
    void allowsAfterWindowExpires() throws InterruptedException {
        var log = new SlidingWindowLog(100, SlidingWindowLog.TimeUnit.MILLISECONDS, 2);
        assertTrue(log.allowRequest());
        assertTrue(log.allowRequest());
        assertFalse(log.allowRequest());
        Thread.sleep(150);
        assertTrue(log.allowRequest());
    }

    @Test
    @DisplayName("reports accurate request count")
    void reportsAccurateRequestCount() {
        var log = new SlidingWindowLog(1, SlidingWindowLog.TimeUnit.SECONDS, 10);
        assertEquals(0, log.getRequestCount());
        assertTrue(log.allowRequest());
        assertEquals(1, log.getRequestCount());
        assertTrue(log.allowRequest());
        assertEquals(2, log.getRequestCount());
    }

    @Test
    @DisplayName("rejects non-positive window size")
    void rejectsNonPositiveWindowSize() {
        assertThrows(IllegalArgumentException.class,
                () -> new SlidingWindowLog(0, SlidingWindowLog.TimeUnit.SECONDS, 5));
        assertThrows(IllegalArgumentException.class,
                () -> new SlidingWindowLog(-1, SlidingWindowLog.TimeUnit.SECONDS, 5));
    }

    @Test
    @DisplayName("rejects non-positive max requests")
    void rejectsNonPositiveMaxRequests() {
        assertThrows(IllegalArgumentException.class,
                () -> new SlidingWindowLog(1, SlidingWindowLog.TimeUnit.SECONDS, 0));
        assertThrows(IllegalArgumentException.class,
                () -> new SlidingWindowLog(1, SlidingWindowLog.TimeUnit.SECONDS, -1));
    }

    @Test
    @DisplayName("handles concurrent access safely")
    @Timeout(10)
    void concurrentAccess() throws InterruptedException {
        var log = new SlidingWindowLog(1, SlidingWindowLog.TimeUnit.SECONDS, 50);
        int threadCount = 5;
        var latch = new CountDownLatch(threadCount);
        var executor = Executors.newFixedThreadPool(threadCount);
        var allowed = new AtomicInteger(0);

        for (int i = 0; i < threadCount; i++) {
            executor.submit(() -> {
                for (int j = 0; j < 20; j++) {
                    if (log.allowRequest()) {
                        allowed.incrementAndGet();
                    }
                }
                latch.countDown();
            });
        }

        latch.await();
        executor.shutdown();
        assertTrue(allowed.get() <= 50, "Should not exceed max requests");
    }
}
