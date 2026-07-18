package com.qe.test.concurrent;

import static org.junit.jupiter.api.Assertions.*;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

@DisplayName("BlockingQueue")
class BlockingQueueTest {

    @Test
    @DisplayName("put and take in FIFO order")
    void putTake() throws Exception {
        var queue = new BlockingQueue<Integer>(5);
        queue.put(1);
        queue.put(2);
        queue.put(3);
        assertEquals(1, queue.take());
        assertEquals(2, queue.take());
        assertEquals(3, queue.take());
    }

    @Test
    @DisplayName("take blocks until item available")
    void takeBlocks() throws Exception {
        var queue = new BlockingQueue<Integer>(5);
        var result = new Integer[1];
        var latch = new CountDownLatch(1);
        var executor = Executors.newSingleThreadExecutor();
        executor.submit(() -> {
            try {
                result[0] = queue.take();
                latch.countDown();
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        });
        Thread.sleep(50);
        assertNull(result[0]);
        queue.put(42);
        latch.await(1, TimeUnit.SECONDS);
        assertEquals(42, result[0]);
        executor.shutdown();
    }

    @Test
    @DisplayName("offer with timeout returns false on full")
    void offerTimeout() throws Exception {
        var queue = new BlockingQueue<Integer>(2);
        queue.put(1);
        queue.put(2);
        assertFalse(queue.offer(3, 100, TimeUnit.MILLISECONDS));
    }

    @Test
    @DisplayName("poll with timeout returns null when empty")
    void pollTimeout() throws Exception {
        var queue = new BlockingQueue<Integer>(5);
        assertNull(queue.poll(100, TimeUnit.MILLISECONDS));
    }

    @Test
    @DisplayName("stress test with multiple producers and consumers")
    void stressTest() throws Exception {
        int producers = 10;
        int consumers = 10;
        int itemsPerProducer = 100;
        var queue = new BlockingQueue<Integer>(50);
        var produced = new AtomicInteger(0);
        var consumed = new AtomicInteger(0);
        var doneLatch = new CountDownLatch(producers + consumers);
        var executor = Executors.newFixedThreadPool(producers + consumers);

        for (int i = 0; i < producers; i++) {
            executor.submit(() -> {
                try {
                    for (int j = 0; j < itemsPerProducer; j++) {
                        queue.put(1);
                        produced.incrementAndGet();
                    }
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                } finally {
                    doneLatch.countDown();
                }
            });
        }

        for (int i = 0; i < consumers; i++) {
            executor.submit(() -> {
                try {
                    while (true) {
                        var item = queue.poll(500, TimeUnit.MILLISECONDS);
                        if (item == null) break;
                        consumed.incrementAndGet();
                    }
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                } finally {
                    doneLatch.countDown();
                }
            });
        }

        doneLatch.await();
        executor.shutdown();
        assertEquals(producers * itemsPerProducer, produced.get());
        assertEquals(produced.get(), consumed.get(),
                "All produced items should be consumed");
    }

    @Test
    @DisplayName("capacity must be positive")
    void capacityValidation() {
        assertThrows(IllegalArgumentException.class, () -> new BlockingQueue<>(0));
        assertThrows(IllegalArgumentException.class, () -> new BlockingQueue<>(-1));
    }

    @Test
    @DisplayName("null item rejected")
    void nullItem() {
        var queue = new BlockingQueue<Integer>(5);
        assertThrows(IllegalArgumentException.class, () -> queue.put(null));
    }
}
