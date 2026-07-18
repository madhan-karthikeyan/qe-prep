package com.qe.test.ds;

import static org.junit.jupiter.api.Assertions.*;

import java.util.NoSuchElementException;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

@DisplayName("CircularQueue")
class CircularQueueTest {

    @Test
    @DisplayName("enqueue and dequeue in order, wrapping around")
    void enqueueDequeue() {
        var queue = new CircularQueue<Integer>(3);
        queue.enqueue(1);
        queue.enqueue(2);
        queue.enqueue(3);
        assertEquals(1, queue.dequeue());
        queue.enqueue(4);
        assertEquals(2, queue.dequeue());
        assertEquals(3, queue.dequeue());
        assertEquals(4, queue.dequeue());
        assertTrue(queue.isEmpty());
    }

    @Test
    @DisplayName("peek returns front without removing")
    void peek() {
        var queue = new CircularQueue<String>(3);
        queue.enqueue("a");
        queue.enqueue("b");
        assertEquals("a", queue.peek());
        assertEquals(2, queue.size());
    }

    @Test
    @DisplayName("enqueue on full throws IllegalStateException")
    void enqueueFull() {
        var queue = new CircularQueue<Integer>(2);
        queue.enqueue(1);
        queue.enqueue(2);
        assertThrows(IllegalStateException.class, () -> queue.enqueue(3));
    }

    @Test
    @DisplayName("dequeue on empty throws NoSuchElementException")
    void dequeueEmpty() {
        var queue = new CircularQueue<Integer>(3);
        assertThrows(NoSuchElementException.class, queue::dequeue);
    }

    @Test
    @DisplayName("wraps around correctly")
    void wrapAround() {
        var queue = new CircularQueue<Integer>(3);
        queue.enqueue(1);
        queue.enqueue(2);
        queue.enqueue(3);
        queue.dequeue();
        queue.dequeue();
        queue.enqueue(4);
        queue.enqueue(5);
        assertEquals(3, queue.dequeue());
        assertEquals(4, queue.dequeue());
        assertEquals(5, queue.dequeue());
        assertTrue(queue.isEmpty());
    }

    @Test
    @DisplayName("isFull returns true when at capacity")
    void isFull() {
        var queue = new CircularQueue<Integer>(2);
        assertFalse(queue.isFull());
        queue.enqueue(1);
        assertFalse(queue.isFull());
        queue.enqueue(2);
        assertTrue(queue.isFull());
    }

    @Test
    @DisplayName("capacity must be positive")
    void capacityValidation() {
        assertThrows(IllegalArgumentException.class, () -> new CircularQueue<>(0));
        assertThrows(IllegalArgumentException.class, () -> new CircularQueue<>(-1));
    }

    @Test
    @DisplayName("null item rejected")
    void nullItem() {
        var queue = new CircularQueue<Integer>(3);
        assertThrows(IllegalArgumentException.class, () -> queue.enqueue(null));
    }
}
