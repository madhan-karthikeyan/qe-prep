#include "circular_queue.hpp"
#include <cassert>
#include <iostream>

void test_cq() {
    CircularQueue<int> q(3);
    assert(q.is_empty());
    assert(!q.is_full());

    q.enqueue(1);
    q.enqueue(2);
    q.enqueue(3);
    assert(q.is_full());

    try {
        q.enqueue(4);
        assert(false);
    } catch (const std::overflow_error &) {}

    assert(q.dequeue() == 1);
    assert(!q.is_full());
    q.enqueue(4);

    assert(q.dequeue() == 2);
    assert(q.dequeue() == 3);
    assert(q.dequeue() == 4);
    assert(q.is_empty());
}

int main() {
    test_cq();
    std::cout << "PASS: circular_queue" << std::endl;
    return 0;
}
