#include "lru_cache.hpp"
#include <cassert>
#include <iostream>

void test_lru_basic() {
    LRUCache<int, int> cache(3);
    cache.put(1, 10);
    cache.put(2, 20);
    cache.put(3, 30);
    assert(cache.get(1) == 10);
    assert(cache.get(2) == 20);
    assert(cache.get(3) == 30);

    cache.put(4, 40);
    assert(cache.get(1) == 0);
    assert(cache.get(4) == 40);
}

void test_lru_update() {
    LRUCache<int, int> cache(2);
    cache.put(1, 10);
    cache.put(1, 100);
    assert(cache.get(1) == 100);
}

int main() {
    test_lru_basic();
    test_lru_update();
    std::cout << "PASS: lru_cache" << std::endl;
    return 0;
}
