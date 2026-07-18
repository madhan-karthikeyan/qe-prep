import pytest
from lru_cache.implementation.lru_cache import LRUCache


@pytest.mark.parametrize("capacity", [100, 1000, 10000])
def test_lru_put(benchmark, capacity):
    cache = LRUCache[int, int](capacity)

    def puts():
        for i in range(capacity):
            cache.put(i, i)

    benchmark(puts)


@pytest.mark.parametrize("capacity", [100, 1000, 10000])
def test_lru_get(benchmark, capacity):
    cache = LRUCache[int, int](capacity)
    for i in range(capacity):
        cache.put(i, i)

    def gets():
        for i in range(capacity):
            cache.get(i)

    benchmark(gets)


@pytest.mark.parametrize("capacity", [100, 1000, 10000])
def test_lru_mixed(benchmark, capacity):
    cache = LRUCache[int, int](capacity)

    def mixed():
        for i in range(capacity * 2):
            cache.put(i, i)
            cache.get(i // 2)

    benchmark(mixed)


@pytest.mark.parametrize("capacity", [100, 1000, 10000])
def test_lru_eviction(benchmark, capacity):
    cache = LRUCache[int, int](capacity)
    for i in range(capacity):
        cache.put(i, i)

    def evictions():
        for i in range(capacity):
            cache.put(i + capacity, i + capacity)

    benchmark(evictions)
