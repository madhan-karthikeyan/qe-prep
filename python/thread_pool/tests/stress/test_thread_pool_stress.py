from __future__ import annotations

from thread_pool.implementation.thread_pool import ThreadPool


def test_1000_tasks_8_workers() -> None:
    pool = ThreadPool(8)
    n_tasks = 1000

    futures = [pool.submit(lambda i=i: i * 2) for i in range(n_tasks)]
    results = [f.result() for f in futures]

    assert len(results) == n_tasks
    assert results == [i * 2 for i in range(n_tasks)]
    pool.shutdown()


def test_all_tasks_complete() -> None:
    pool = ThreadPool(4)
    n = 500

    completed = [False] * n

    def mark(i: int) -> int:
        completed[i] = True
        return i

    futures = [pool.submit(mark, i) for i in range(n)]
    for f in futures:
        f.result()

    assert all(completed)
    pool.shutdown()
