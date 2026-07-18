from __future__ import annotations

import time

import pytest

from thread_pool.implementation.thread_pool import ThreadPool


class TestThreadPool:
    def test_submit_and_get_result(self) -> None:
        pool = ThreadPool(2)
        future = pool.submit(lambda x: x * 2, 21)
        assert future.result() == 42
        pool.shutdown()

    def test_multiple_tasks(self) -> None:
        pool = ThreadPool(4)
        futures = [pool.submit(lambda a, b: a + b, i, i) for i in range(10)]
        results = [f.result() for f in futures]
        assert results == [i * 2 for i in range(10)]
        pool.shutdown()

    def test_shutdown_rejects_new_tasks(self) -> None:
        pool = ThreadPool(2)
        pool.shutdown()
        with pytest.raises(RuntimeError):
            pool.submit(lambda: 1)

    def test_exception_propagation(self) -> None:
        pool = ThreadPool(2)

        def failing() -> None:
            raise ValueError("oops")

        future = pool.submit(failing)
        with pytest.raises(ValueError):
            future.result()
        pool.shutdown()

    def test_shutdown_waits_for_pending(self) -> None:
        pool = ThreadPool(2)

        def slow() -> int:
            time.sleep(0.2)
            return 42

        future = pool.submit(slow)
        pool.shutdown(wait=True)
        assert future.result() == 42

    def test_invalid_worker_count(self) -> None:
        with pytest.raises(ValueError):
            ThreadPool(0)
