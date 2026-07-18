package com.qe.test.concurrent;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.Executors;
import java.util.concurrent.atomic.AtomicInteger;

public class ProducerConsumerDemo {
    private final BlockingQueue<Integer> queue;
    private final int producerCount;
    private final int consumerCount;
    private final int itemsPerProducer;

    public ProducerConsumerDemo(int queueCapacity, int producerCount, int consumerCount, int itemsPerProducer) {
        this.queue = new BlockingQueue<>(queueCapacity);
        this.producerCount = producerCount;
        this.consumerCount = consumerCount;
        this.itemsPerProducer = itemsPerProducer;
    }

    public record Result(int produced, int consumed, long elapsedMs) { }

    public Result run() throws InterruptedException {
        var executor = Executors.newFixedThreadPool(producerCount + consumerCount);
        var producedCounter = new AtomicInteger(0);
        var consumedCounter = new AtomicInteger(0);
        var startLatch = new CountDownLatch(1);
        var doneLatch = new CountDownLatch(producerCount + consumerCount);

        for (int i = 0; i < producerCount; i++) {
            int id = i;
            executor.submit(() -> {
                try {
                    startLatch.await();
                    for (int j = 0; j < itemsPerProducer; j++) {
                        queue.put(id * itemsPerProducer + j);
                        producedCounter.incrementAndGet();
                    }
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                } finally {
                    doneLatch.countDown();
                }
            });
        }

        for (int i = 0; i < consumerCount; i++) {
            executor.submit(() -> {
                try {
                    startLatch.await();
                    while (true) {
                        var item = queue.poll(100, java.util.concurrent.TimeUnit.MILLISECONDS);
                        if (item == null) break;
                        consumedCounter.incrementAndGet();
                    }
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                } finally {
                    doneLatch.countDown();
                }
            });
        }

        long start = System.currentTimeMillis();
        startLatch.countDown();
        doneLatch.await();
        long elapsed = System.currentTimeMillis() - start;
        executor.shutdown();
        return new Result(producedCounter.get(), consumedCounter.get(), elapsed);
    }
}
