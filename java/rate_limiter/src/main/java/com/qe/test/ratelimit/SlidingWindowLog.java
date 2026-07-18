package com.qe.test.ratelimit;

import java.util.ArrayDeque;
import java.util.Deque;
import java.util.concurrent.locks.ReentrantLock;

public class SlidingWindowLog {
    private final long windowSizeNanos;
    private final int maxRequests;
    private final Deque<Long> timestamps = new ArrayDeque<>();
    private final ReentrantLock lock = new ReentrantLock();

    public SlidingWindowLog(long windowSize, TimeUnit unit, int maxRequests) {
        if (windowSize <= 0) {
            throw new IllegalArgumentException("windowSize must be positive: " + windowSize);
        }
        if (maxRequests <= 0) {
            throw new IllegalArgumentException("maxRequests must be positive: " + maxRequests);
        }
        this.windowSizeNanos = unit.toNanos(windowSize);
        this.maxRequests = maxRequests;
    }

    public boolean allowRequest() {
        lock.lock();
        try {
            long now = System.nanoTime();
            long cutoff = now - windowSizeNanos;
            while (!timestamps.isEmpty() && timestamps.peekFirst() < cutoff) {
                timestamps.removeFirst();
            }
            if (timestamps.size() < maxRequests) {
                timestamps.addLast(now);
                return true;
            }
            return false;
        } finally {
            lock.unlock();
        }
    }

    public int getRequestCount() {
        lock.lock();
        try {
            long now = System.nanoTime();
            long cutoff = now - windowSizeNanos;
            while (!timestamps.isEmpty() && timestamps.peekFirst() < cutoff) {
                timestamps.removeFirst();
            }
            return timestamps.size();
        } finally {
            lock.unlock();
        }
    }

    public long getWindowSizeNanos() {
        return windowSizeNanos;
    }

    public int getMaxRequests() {
        return maxRequests;
    }

    public enum TimeUnit {
        NANOSECONDS(1L),
        MILLISECONDS(1_000_000L),
        SECONDS(1_000_000_000L);

        private final long multiplier;

        TimeUnit(long multiplier) {
            this.multiplier = multiplier;
        }

        public long toNanos(long value) {
            return value * multiplier;
        }
    }
}
