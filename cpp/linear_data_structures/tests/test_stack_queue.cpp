#include "stack.hpp"
#include "queue.hpp"
#include <cassert>
#include <iostream>

void test_stack() {
    Stack<int> s;
    assert(s.is_empty());
    s.push(1);
    s.push(2);
    s.push(3);
    assert(!s.is_empty());
    assert(s.peek() == 3);
    assert(s.pop() == 3);
    assert(s.pop() == 2);
    assert(s.pop() == 1);
    assert(s.is_empty());
}

void test_queue() {
    Queue<int> q;
    assert(q.is_empty());
    q.enqueue(1);
    q.enqueue(2);
    q.enqueue(3);
    assert(!q.is_empty());
    assert(q.peek() == 1);
    assert(q.dequeue() == 1);
    assert(q.dequeue() == 2);
    assert(q.dequeue() == 3);
    assert(q.is_empty());
}

int main() {
    test_stack();
    test_queue();
    std::cout << "PASS: stack_queue" << std::endl;
    return 0;
}
