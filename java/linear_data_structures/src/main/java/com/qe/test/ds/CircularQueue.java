package com.qe.test.ds;

import java.util.NoSuchElementException;

public class CircularQueue<E> {
    private final Object[] buffer;
    private int head;
    private int tail;
    private int count;

    @SuppressWarnings("unchecked")
    public CircularQueue(int capacity) {
        if (capacity <= 0) throw new IllegalArgumentException("capacity must be positive");
        this.buffer = new Object[capacity];
        this.head = 0;
        this.tail = 0;
        this.count = 0;
    }

    public void enqueue(E item) {
        if (item == null) throw new IllegalArgumentException("item must not be null");
        if (isFull()) throw new IllegalStateException("queue is full");
        buffer[tail] = item;
        tail = (tail + 1) % buffer.length;
        count++;
    }

    @SuppressWarnings("unchecked")
    public E dequeue() {
        if (isEmpty()) throw new NoSuchElementException("queue is empty");
        var item = (E) buffer[head];
        buffer[head] = null;
        head = (head + 1) % buffer.length;
        count--;
        return item;
    }

    @SuppressWarnings("unchecked")
    public E peek() {
        if (isEmpty()) throw new NoSuchElementException("queue is empty");
        return (E) buffer[head];
    }

    public boolean isEmpty() {
        return count == 0;
    }

    public boolean isFull() {
        return count == buffer.length;
    }

    public int size() {
        return count;
    }

    public int capacity() {
        return buffer.length;
    }
}
