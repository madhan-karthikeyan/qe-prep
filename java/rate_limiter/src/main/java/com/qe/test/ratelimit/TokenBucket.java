package com.qe.test.ratelimit;

import java.util.concurrent.locks.ReentrantLock;

public class TokenBucket {
    private final long burstCapacity;
    private final double refillRate;
    private double tokens;
    private long lastRefillTimestamp;
    private final ReentrantLock lock = new ReentrantLock();

    public TokenBucket(long burstCapacity, double refillRate) {
        if (burstCapacity <= 0) {
            throw new IllegalArgumentException("burstCapacity must be positive: " + burstCapacity);
        }
        if (refillRate <= 0) {
            throw new IllegalArgumentException("refillRate must be positive: " + refillRate);
        }
        this.burstCapacity = burstCapacity;
        this.refillRate = refillRate;
        this.tokens = burstCapacity;
        this.lastRefillTimestamp = System.nanoTime();
    }

    public boolean allowRequest() {
        lock.lock();
        try {
            refill();
            if (tokens >= 1.0) {
                tokens -= 1.0;
                return true;
            }
            return false;
        } finally {
            lock.unlock();
        }
    }

    public boolean allowRequest(long tokens) {
        if (tokens <= 0) {
            throw new IllegalArgumentException("tokens must be positive: " + tokens);
        }
        lock.lock();
        try {
            refill();
            if (this.tokens >= tokens) {
                this.tokens -= tokens;
                return true;
            }
            return false;
        } finally {
            lock.unlock();
        }
    }

    public double getAvailableTokens() {
        lock.lock();
        try {
            refill();
            return tokens;
        } finally {
            lock.unlock();
        }
    }

    private void refill() {
        long now = System.nanoTime();
        double elapsed = (now - lastRefillTimestamp) / 1_000_000_000.0;
        tokens = Math.min(burstCapacity, tokens + elapsed * refillRate);
        lastRefillTimestamp = now;
    }
}
