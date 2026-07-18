package com.qe.test.concurrent;

import java.util.concurrent.TimeUnit;
import java.util.concurrent.locks.Condition;
import java.util.concurrent.locks.ReentrantLock;

public class BlockingQueue<T> {
    private final Object[] buffer;
    private int head;
    private int tail;
    private int count;
    private final ReentrantLock lock;
    private final Condition notEmpty;
    private final Condition notFull;

    public BlockingQueue(int capacity) {
        if (capacity <= 0) throw new IllegalArgumentException("capacity must be positive");
        this.buffer = new Object[capacity];
        this.lock = new ReentrantLock();
        this.notEmpty = lock.newCondition();
        this.notFull = lock.newCondition();
    }

    public void put(T item) throws InterruptedException {
        if (item == null) throw new IllegalArgumentException("item must not be null");
        lock.lockInterruptibly();
        try {
            while (isFull()) {
                notFull.await();
            }
            enqueue(item);
        } finally {
            lock.unlock();
        }
    }

    public boolean offer(T item, long timeout, TimeUnit unit) throws InterruptedException {
        if (item == null) throw new IllegalArgumentException("item must not be null");
        long nanos = unit.toNanos(timeout);
        lock.lockInterruptibly();
        try {
            while (isFull()) {
                if (nanos <= 0) return false;
                nanos = notFull.awaitNanos(nanos);
            }
            enqueue(item);
            return true;
        } finally {
            lock.unlock();
        }
    }

    @SuppressWarnings("unchecked")
    public T take() throws InterruptedException {
        lock.lockInterruptibly();
        try {
            while (isEmpty()) {
                notEmpty.await();
            }
            return dequeue();
        } finally {
            lock.unlock();
        }
    }

    @SuppressWarnings("unchecked")
    public T poll(long timeout, TimeUnit unit) throws InterruptedException {
        long nanos = unit.toNanos(timeout);
        lock.lockInterruptibly();
        try {
            while (isEmpty()) {
                if (nanos <= 0) return null;
                nanos = notEmpty.awaitNanos(nanos);
            }
            return dequeue();
        } finally {
            lock.unlock();
        }
    }

    public int size() {
        lock.lock();
        try {
            return count;
        } finally {
            lock.unlock();
        }
    }

    public int capacity() {
        return buffer.length;
    }

    private void enqueue(T item) {
        buffer[tail] = item;
        tail = (tail + 1) % buffer.length;
        count++;
        notEmpty.signal();
    }

    @SuppressWarnings("unchecked")
    private T dequeue() {
        var item = (T) buffer[head];
        buffer[head] = null;
        head = (head + 1) % buffer.length;
        count--;
        notFull.signal();
        return item;
    }

    private boolean isEmpty() {
        return count == 0;
    }

    private boolean isFull() {
        return count == buffer.length;
    }
}
