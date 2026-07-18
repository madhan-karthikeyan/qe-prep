package com.qe.test.ds;

import static org.junit.jupiter.api.Assertions.*;

import java.util.EmptyStackException;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

@DisplayName("Stack")
class StackTest {

    @Test
    @DisplayName("push and pop LIFO order")
    void pushPop() {
        var stack = new Stack<Integer>();
        stack.push(1);
        stack.push(2);
        stack.push(3);
        assertEquals(3, stack.pop());
        assertEquals(2, stack.pop());
        assertEquals(1, stack.pop());
    }

    @Test
    @DisplayName("peek returns top without removing")
    void peek() {
        var stack = new Stack<String>();
        stack.push("a");
        stack.push("b");
        assertEquals("b", stack.peek());
        assertEquals(2, stack.size());
    }

    @Test
    @DisplayName("pop on empty throws EmptyStackException")
    void popEmpty() {
        var stack = new Stack<Integer>();
        assertThrows(EmptyStackException.class, stack::pop);
    }

    @Test
    @DisplayName("peek on empty throws EmptyStackException")
    void peekEmpty() {
        var stack = new Stack<Integer>();
        assertThrows(EmptyStackException.class, stack::peek);
    }

    @Test
    @DisplayName("isEmpty returns true for empty stack")
    void isEmpty() {
        var stack = new Stack<Integer>();
        assertTrue(stack.isEmpty());
        stack.push(1);
        assertFalse(stack.isEmpty());
        stack.pop();
        assertTrue(stack.isEmpty());
    }

    @Test
    @DisplayName("null item rejected")
    void nullItem() {
        var stack = new Stack<Integer>();
        assertThrows(IllegalArgumentException.class, () -> stack.push(null));
    }
}
