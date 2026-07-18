package com.qe.test.net;

import java.util.concurrent.BlockingQueue;
import java.util.concurrent.LinkedBlockingQueue;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.locks.ReentrantLock;

public class ConnectionPool<T extends AutoCloseable> {
    private final BlockingQueue<PooledConnection<T>> pool;
    private final ConnectionFactory<T> factory;
    private final int maxConnections;
    private final AtomicInteger createdCount = new AtomicInteger(0);
    private final ReentrantLock createLock = new ReentrantLock();
    private volatile boolean closed = false;

    @FunctionalInterface
    public interface ConnectionFactory<T extends AutoCloseable> {
        T create() throws Exception;
    }

    public ConnectionPool(ConnectionFactory<T> factory, int maxConnections) {
        if (factory == null) {
            throw new IllegalArgumentException("factory must not be null");
        }
        if (maxConnections <= 0) {
            throw new IllegalArgumentException("maxConnections must be positive: " + maxConnections);
        }
        this.factory = factory;
        this.maxConnections = maxConnections;
        this.pool = new LinkedBlockingQueue<>(maxConnections);
    }

    public PooledConnection<T> acquire(long timeout, TimeUnit unit) throws Exception {
        if (closed) {
            throw new IllegalStateException("Connection pool is closed");
        }
        PooledConnection<T> conn = pool.poll();
        if (conn != null) {
            return conn;
        }
        if (createdCount.get() < maxConnections) {
            createLock.lock();
            try {
                if (createdCount.get() < maxConnections) {
                    T raw = factory.create();
                    createdCount.incrementAndGet();
                    return new PooledConnection<>(raw, this);
                }
            } finally {
                createLock.unlock();
            }
        }
        conn = pool.poll(timeout, unit);
        if (conn == null) {
            throw new IllegalStateException("Timeout waiting for connection from pool");
        }
        return conn;
    }

    public void release(PooledConnection<T> conn) {
        if (closed) {
            closeConnection(conn);
            return;
        }
        if (!pool.offer(conn)) {
            closeConnection(conn);
        }
    }

    public void closeAll() {
        closed = true;
        PooledConnection<T> conn;
        while ((conn = pool.poll()) != null) {
            closeConnection(conn);
        }
    }

    private void closeConnection(PooledConnection<T> conn) {
        try {
            conn.actualClose();
        } catch (Exception e) {
            System.err.println("Error closing connection: " + e.getMessage());
        }
    }

    public int getAvailableCount() {
        return pool.size();
    }

    public int getCreatedCount() {
        return createdCount.get();
    }

    public boolean isClosed() {
        return closed;
    }

    public static class PooledConnection<T extends AutoCloseable> implements AutoCloseable {
        private final T delegate;
        private final ConnectionPool<T> pool;
        private volatile boolean released = false;

        PooledConnection(T delegate, ConnectionPool<T> pool) {
            this.delegate = delegate;
            this.pool = pool;
        }

        public T get() {
            if (released) {
                throw new IllegalStateException("Connection has been released");
            }
            return delegate;
        }

        @Override
        public void close() {
            if (!released) {
                released = true;
                pool.release(this);
            }
        }

        void actualClose() throws Exception {
            delegate.close();
        }
    }
}
