package com.qe.test.ds;

import java.util.LinkedList;
import java.util.NoSuchElementException;

public class Queue<E> {
    private final LinkedList<E> elements;

    public Queue() {
        this.elements = new LinkedList<>();
    }

    public void enqueue(E item) {
        if (item == null) throw new IllegalArgumentException("item must not be null");
        elements.addLast(item);
    }

    public E dequeue() {
        if (isEmpty()) throw new NoSuchElementException("queue is empty");
        return elements.removeFirst();
    }

    public E peek() {
        if (isEmpty()) throw new NoSuchElementException("queue is empty");
        return elements.getFirst();
    }

    public boolean isEmpty() {
        return elements.isEmpty();
    }

    public int size() {
        return elements.size();
    }
}
