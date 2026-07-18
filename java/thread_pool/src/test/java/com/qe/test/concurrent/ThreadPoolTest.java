package com.qe.test.concurrent;

import static org.junit.jupiter.api.Assertions.*;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

@DisplayName("ThreadPool")
class ThreadPoolTest {

    @Test
    @DisplayName("submits tasks and returns correct results")
    void submitAndGetResult() throws Exception {
        try (var pool = new ThreadPool(4)) {
            var future = pool.submit(() -> 42);
            assertEquals(42, future.get());
        }
    }

    @Test
    @DisplayName("handles multiple tasks")
    void multipleTasks() throws Exception {
        try (var pool = new ThreadPool(4)) {
            var futures = new java.util.ArrayList<java.util.concurrent.Future<Integer>>();
            for (int i = 0; i < 20; i++) {
                int val = i;
                futures.add(pool.submit(() -> val * 2));
            }
            for (int i = 0; i < 20; i++) {
                assertEquals(i * 2, futures.get(i).get());
            }
        }
    }

    @Test
    @DisplayName("propagates exception from task")
    void taskException() throws Exception {
        try (var pool = new ThreadPool(2)) {
            var future = pool.submit(() -> { throw new RuntimeException("fail"); });
            assertThrows(ExecutionException.class, future::get);
        }
    }

    @Test
    @DisplayName("shutdown prevents new tasks")
    void shutdownRejectsTasks() throws Exception {
        var pool = new ThreadPool(2);
        pool.shutdown();
        pool.awaitTermination(1, TimeUnit.SECONDS);
        assertThrows(IllegalStateException.class, () -> pool.submit(() -> 1));
    }

    @Test
    @DisplayName("awaitTermination returns true when all tasks done")
    void awaitTerminationAllDone() throws Exception {
        try (var pool = new ThreadPool(2)) {
            var f1 = pool.submit(() -> 1);
            var f2 = pool.submit(() -> 2);
            assertEquals(1, f1.get());
            assertEquals(2, f2.get());
            pool.shutdown();
            assertTrue(pool.awaitTermination(1, TimeUnit.SECONDS));
            assertTrue(pool.isTerminated());
        }
    }

    @Test
    @DisplayName("shutdown and awaitTermination with pending tasks")
    void shutdownWithPending() throws Exception {
        var pool = new ThreadPool(2);
        var latch = new CountDownLatch(1);
        var result = pool.submit(() -> {
            latch.await();
            return 42;
        });
        pool.shutdown();
        // awaitTermination should time out because latch is not released
        assertFalse(pool.awaitTermination(200, TimeUnit.MILLISECONDS));
        latch.countDown();
        assertTrue(pool.awaitTermination(5, TimeUnit.SECONDS));
        assertTrue(pool.isTerminated());
        assertEquals(42, result.get());
    }

    @Test
    @DisplayName("stress test with 1000 tasks and 8 workers")
    void stressTest() throws Exception {
        var pool = new ThreadPool(8);
        var futures = new java.util.ArrayList<java.util.concurrent.Future<Integer>>();
        for (int i = 0; i < 1000; i++) {
            int val = i;
            futures.add(pool.submit(() -> val * val));
        }
        for (int i = 0; i < 1000; i++) {
            assertEquals(i * i, futures.get(i).get());
        }
        pool.shutdown();
        assertTrue(pool.awaitTermination(5, TimeUnit.SECONDS));
    }

    @Test
    @DisplayName("zero workers rejected")
    void zeroWorkers() {
        assertThrows(IllegalArgumentException.class, () -> new ThreadPool(0));
    }
}
