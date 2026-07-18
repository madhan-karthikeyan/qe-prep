package com.qe.test.ds;

import static org.junit.jupiter.api.Assertions.*;

import java.util.NoSuchElementException;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

@DisplayName("Queue")
class QueueTest {

    @Test
    @DisplayName("enqueue and dequeue FIFO order")
    void enqueueDequeue() {
        var queue = new Queue<Integer>();
        queue.enqueue(1);
        queue.enqueue(2);
        queue.enqueue(3);
        assertEquals(1, queue.dequeue());
        assertEquals(2, queue.dequeue());
        assertEquals(3, queue.dequeue());
    }

    @Test
    @DisplayName("peek returns front without removing")
    void peek() {
        var queue = new Queue<String>();
        queue.enqueue("a");
        queue.enqueue("b");
        assertEquals("a", queue.peek());
        assertEquals(2, queue.size());
    }

    @Test
    @DisplayName("dequeue on empty throws NoSuchElementException")
    void dequeueEmpty() {
        var queue = new Queue<Integer>();
        assertThrows(NoSuchElementException.class, queue::dequeue);
    }

    @Test
    @DisplayName("isEmpty returns correct state")
    void isEmpty() {
        var queue = new Queue<Integer>();
        assertTrue(queue.isEmpty());
        queue.enqueue(1);
        assertFalse(queue.isEmpty());
        queue.dequeue();
        assertTrue(queue.isEmpty());
    }

    @Test
    @DisplayName("null item rejected")
    void nullItem() {
        var queue = new Queue<Integer>();
        assertThrows(IllegalArgumentException.class, () -> queue.enqueue(null));
    }
}
