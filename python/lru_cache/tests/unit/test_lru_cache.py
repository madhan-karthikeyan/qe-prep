from __future__ import annotations

import threading

import pytest

from lru_cache.implementation.lru_cache import LRUCache


class TestLRUCache:
    def test_get_missing(self) -> None:
        cache = LRUCache[str, int](3)
        assert cache.get("a") is None

    def test_put_and_get(self) -> None:
        cache = LRUCache[str, int](2)
        cache.put("a", 1)
        assert cache.get("a") == 1

    def test_eviction_when_full(self) -> None:
        cache = LRUCache[str, int](2)
        cache.put("a", 1)
        cache.put("b", 2)
        cache.put("c", 3)
        assert cache.get("a") is None
        assert cache.get("b") == 2
        assert cache.get("c") == 3

    def test_update_existing(self) -> None:
        cache = LRUCache[str, int](2)
        cache.put("a", 1)
        cache.put("a", 99)
        assert cache.get("a") == 99
        assert len(cache) == 1

    def test_access_renews_lru(self) -> None:
        cache = LRUCache[str, int](2)
        cache.put("a", 1)
        cache.put("b", 2)
        cache.get("a")
        cache.put("c", 3)
        assert cache.get("b") is None
        assert cache.get("a") == 1
        assert cache.get("c") == 3

    def test_none_value(self) -> None:
        cache = LRUCache[str, object](2)
        cache.put("a", None)
        assert cache.get("a") is None

    def test_capacity_edge(self) -> None:
        cache = LRUCache[str, int](1)
        cache.put("a", 1)
        cache.put("b", 2)
        assert cache.get("a") is None
        assert cache.get("b") == 2

    def test_len(self) -> None:
        cache = LRUCache[str, int](5)
        assert len(cache) == 0
        cache.put("a", 1)
        assert len(cache) == 1

    def test_invalid_capacity(self) -> None:
        with pytest.raises(ValueError):
            LRUCache[str, int](0)

    def test_thread_safety(self) -> None:
        cache = LRUCache[int, int](100, thread_safe=True)
        n_threads = 10
        ops_per_thread = 1000
        barrier = threading.Barrier(n_threads)

        def worker() -> None:
            barrier.wait()
            for i in range(ops_per_thread):
                cache.put(i, i)
                cache.get(i)

        threads = [threading.Thread(target=worker) for _ in range(n_threads)]
        for t in threads:
            t.start()
        for t in threads:
            t.join()

        assert len(cache) <= 100
