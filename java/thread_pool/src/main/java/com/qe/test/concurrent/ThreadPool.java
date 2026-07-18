package com.qe.test.concurrent;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.Callable;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Future;
import java.util.concurrent.LinkedBlockingQueue;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;

public class ThreadPool implements AutoCloseable {
    private final List<Worker> workers;
    private final BlockingQueue<Task<?>> taskQueue;
    private final AtomicBoolean shutdown;
    private final AtomicBoolean terminated;

    public ThreadPool(int numWorkers) {
        if (numWorkers <= 0) throw new IllegalArgumentException("numWorkers must be positive");
        this.workers = new ArrayList<>(numWorkers);
        this.taskQueue = new LinkedBlockingQueue<>();
        this.shutdown = new AtomicBoolean(false);
        this.terminated = new AtomicBoolean(false);
        for (int i = 0; i < numWorkers; i++) {
            var worker = new Worker();
            workers.add(worker);
            worker.start();
        }
    }

    public <T> Future<T> submit(Callable<T> task) {
        if (task == null) throw new IllegalArgumentException("task must not be null");
        if (shutdown.get()) throw new IllegalStateException("pool is shut down");
        var future = new CompletableFuture<T>();
        taskQueue.add(new Task<>(task, future));
        return future;
    }

    public void shutdown() {
        shutdown.set(true);
    }

    public boolean awaitTermination(long timeout, TimeUnit unit) throws InterruptedException {
        shutdown();
        long deadline = System.nanoTime() + unit.toNanos(timeout);
        for (var worker : workers) {
            long remaining = deadline - System.nanoTime();
            if (remaining <= 0) return false;
            worker.join(TimeUnit.NANOSECONDS.toMillis(remaining));
        }
        for (var worker : workers) {
            if (worker.isAlive()) return false;
        }
        terminated.set(true);
        return true;
    }

    public boolean isShutdown() { return shutdown.get(); }
    public boolean isTerminated() { return terminated.get(); }

    @Override
    public void close() throws Exception {
        shutdown();
        awaitTermination(5, TimeUnit.SECONDS);
    }

    private record Task<T>(Callable<T> callable, CompletableFuture<T> future) { }

    private class Worker extends Thread {
        @Override
        public void run() {
            while (!shutdown.get() || !taskQueue.isEmpty()) {
                try {
                    var task = taskQueue.poll(100, TimeUnit.MILLISECONDS);
                    if (task == null) continue;
                    executeTask(task);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    break;
                }
            }
        }

        private <T> void executeTask(Task<T> task) {
            try {
                T result = task.callable().call();
                task.future().complete(result);
            } catch (Exception e) {
                task.future().completeExceptionally(e);
            }
        }
    }
}
