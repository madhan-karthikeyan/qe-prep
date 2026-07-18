#include "producer_consumer.hpp"
#include <cassert>
#include <iostream>
#include <thread>

void test_pc() {
    BlockingQueue<int> q(5);
    std::thread prod[3], cons[3];

    for (int i = 0; i < 3; i++) {
        prod[i] = std::thread([&q, i]() { q.put(i + 1); });
    }
    for (int i = 0; i < 3; i++) {
        cons[i] = std::thread([&q]() { q.get(); });
    }

    for (int i = 0; i < 3; i++) {
        prod[i].join();
        cons[i].join();
    }
}

int main() {
    test_pc();
    std::cout << "PASS: producer_consumer" << std::endl;
    return 0;
}
