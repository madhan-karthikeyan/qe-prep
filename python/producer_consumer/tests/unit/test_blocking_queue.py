from __future__ import annotations

import threading

from producer_consumer.implementation.blocking_queue import BlockingQueue


class TestBlockingQueue:
    def test_put_and_get(self) -> None:
        q = BlockingQueue[int](3)
        q.put(1)
        q.put(2)
        assert q.get() == 1
        assert q.get() == 2

    def test_get_blocks_when_empty_with_timeout(self) -> None:
        q = BlockingQueue[int](3)
        result = q.get(timeout=0.1)
        assert result is None

    def test_put_blocks_when_full_with_timeout(self) -> None:
        q = BlockingQueue[int](2)
        q.put(1)
        q.put(2)
        result = q.put(3, timeout=0.1)
        assert not result

    def test_fifo_order(self) -> None:
        q = BlockingQueue[int](5)
        for i in range(5):
            q.put(i)
        for i in range(5):
            assert q.get() == i

    def test_multiple_producers_consumers(self) -> None:
        q = BlockingQueue[int](10)
        n = 100
        produced: list[int] = []
        consumed: list[int] = []
        lock = threading.Lock()

        def prod() -> None:
            for i in range(n):
                q.put(i)
                with lock:
                    produced.append(i)

        def cons() -> None:
            for _ in range(n):
                item = q.get()
                with lock:
                    consumed.append(item)

        threads = [threading.Thread(target=prod)] + [threading.Thread(target=cons)]
        for t in threads:
            t.start()
        for t in threads:
            t.join()

        assert sorted(produced) == sorted(consumed)

    def test_empty_queue_size(self) -> None:
        q = BlockingQueue[int](5)
        assert q.size == 0

    def test_non_empty_queue_size(self) -> None:
        q = BlockingQueue[int](5)
        q.put(1)
        q.put(2)
        assert q.size == 2
