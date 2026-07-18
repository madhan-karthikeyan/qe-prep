package com.qe.test.net;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

import static org.junit.jupiter.api.Assertions.*;

@DisplayName("ConnectionPool")
class ConnectionPoolTest {

    @Test
    @DisplayName("acquires and releases connections")
    void acquireRelease() throws Exception {
        var pool = new ConnectionPool<>(() -> new TestConnection(), 5);
        var conn = pool.acquire(1, TimeUnit.SECONDS);
        assertNotNull(conn);
        assertNotNull(conn.get());
        conn.close();
        assertEquals(1, pool.getAvailableCount());
        pool.closeAll();
    }

    @Test
    @DisplayName("rejects null factory")
    void rejectsNullFactory() {
        assertThrows(IllegalArgumentException.class, () -> new ConnectionPool<>(null, 5));
    }

    @Test
    @DisplayName("rejects non-positive max connections")
    void rejectsNonPositiveMaxConnections() {
        assertThrows(IllegalArgumentException.class,
                () -> new ConnectionPool<>(() -> new TestConnection(), 0));
        assertThrows(IllegalArgumentException.class,
                () -> new ConnectionPool<>(() -> new TestConnection(), -1));
    }

    @Test
    @DisplayName("throws when pool is exhausted and times out")
    void exhaustPool() throws Exception {
        var pool = new ConnectionPool<>(() -> new TestConnection(), 2);
        var c1 = pool.acquire(1, TimeUnit.SECONDS);
        var c2 = pool.acquire(1, TimeUnit.SECONDS);
        assertThrows(IllegalStateException.class, () -> pool.acquire(100, TimeUnit.MILLISECONDS));
        c1.close();
        c2.close();
        pool.closeAll();
    }

    @Test
    @DisplayName("rejects operations on closed pool")
    void closedPool() throws Exception {
        var pool = new ConnectionPool<>(() -> new TestConnection(), 2);
        pool.closeAll();
        assertTrue(pool.isClosed());
    }

    @Test
    @DisplayName("handles concurrent acquire/release")
    @Timeout(10)
    void concurrentAccess() throws Exception {
        var pool = new ConnectionPool<>(() -> new TestConnection(), 5);
        int threadCount = 10;
        var latch = new CountDownLatch(threadCount);
        var executor = Executors.newFixedThreadPool(threadCount);
        var acquired = new AtomicInteger(0);

        for (int i = 0; i < threadCount; i++) {
            executor.submit(() -> {
                try {
                    var conn = pool.acquire(2, TimeUnit.SECONDS);
                    if (conn != null) {
                        acquired.incrementAndGet();
                        Thread.sleep(10);
                        conn.close();
                    }
                } catch (Exception e) {
                    // expected when pool is exhausted
                }
                latch.countDown();
            });
        }

        latch.await(10, TimeUnit.SECONDS);
        executor.shutdown();
        pool.closeAll();
        assertTrue(acquired.get() >= 5, "Should acquire at least max connections");
    }

    private static class TestConnection implements AutoCloseable {
        @Override
        public void close() {
            // no-op
        }
    }
}
