import threading
import time

import pytest

from worker_pool.implementation import WorkerPool


def test_all_tasks_complete():
    pool = WorkerPool(num_workers=3)
    pool.start()

    results: list[int] = []
    lock = threading.Lock()

    def task(n: int) -> int:
        return n * 2

    def cb(result: int):
        with lock:
            results.append(result)

    for i in range(20):
        pool.submit(task, i, result_callback=cb)

    pool.shutdown(wait=True)

    assert sorted(results) == [i * 2 for i in range(20)]


def test_shutdown_blocks_until_done():
    pool = WorkerPool(num_workers=2)
    pool.start()

    event = threading.Event()

    def slow_task():
        event.wait()
        return 42

    pool.submit(slow_task)

    start = time.monotonic()
    t = threading.Thread(target=lambda: (time.sleep(0.05), event.set()))
    t.start()

    pool.shutdown(wait=True)
    elapsed = time.monotonic() - start
    assert elapsed >= 0.05


def test_map_distributes_work():
    pool = WorkerPool(num_workers=4)
    pool.start()

    results: list[int] = []
    lock = threading.Lock()

    def cb(result: int):
        with lock:
            results.append(result)

    pool.map(lambda x: x**2, [1, 2, 3, 4, 5], result_callback=cb)
    pool.shutdown(wait=True)

    assert sorted(results) == [1, 4, 9, 16, 25]


def test_submit_after_shutdown_raises():
    pool = WorkerPool(num_workers=1)
    pool.start()
    pool.shutdown(wait=True)

    with pytest.raises(RuntimeError, match="shut down"):
        pool.submit(lambda: 1)


def test_result_callback_on_exception():
    pool = WorkerPool(num_workers=1)
    pool.start()

    exceptions: list[Exception] = []
    lock = threading.Lock()

    def task():
        raise ValueError("boom")

    def cb(result):
        with lock:
            exceptions.append(result)

    pool.submit(task, result_callback=cb)
    pool.shutdown(wait=True)

    assert len(exceptions) == 1
    assert isinstance(exceptions[0], ValueError)
