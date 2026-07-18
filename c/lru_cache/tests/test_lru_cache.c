#include "lru_cache.h"
#include <assert.h>
#include <stdio.h>

void test_lru_basic() {
    lru_cache_t *cache = lru_create(3);
    assert(cache);

    lru_put(cache, 1, 10);
    lru_put(cache, 2, 20);
    lru_put(cache, 3, 30);
    assert(lru_get(cache, 1) == 10);
    assert(lru_get(cache, 2) == 20);
    assert(lru_get(cache, 3) == 30);

    lru_put(cache, 4, 40);
    assert(lru_get(cache, 1) == -1);
    assert(lru_get(cache, 4) == 40);

    lru_free(cache);
}

void test_lru_update() {
    lru_cache_t *cache = lru_create(2);
    lru_put(cache, 1, 10);
    lru_put(cache, 1, 100);
    assert(lru_get(cache, 1) == 100);
    lru_free(cache);
}

int main() {
    test_lru_basic();
    test_lru_update();
    printf("PASS: lru_cache\n");
    return 0;
}
