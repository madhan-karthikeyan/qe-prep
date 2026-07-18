from __future__ import annotations

import threading

from producer_consumer.implementation.blocking_queue import BlockingQueue

_SENTINEL = -1


def test_10_producers_10_consumers_10k_items() -> None:
    q = BlockingQueue[int](50)
    n_producers = 10
    n_consumers = 10
    items_per_producer = 1_000
    total = n_producers * items_per_producer

    produced = [0] * total
    consumed = [0] * total
    prod_counter: list[int] = [0]
    cons_counter: list[int] = [0]
    prod_lock = threading.Lock()

    def producer_worker(pid: int) -> None:
        start = pid * items_per_producer
        for i in range(start, start + items_per_producer):
            q.put(i)
            with prod_lock:
                produced[i] = 1
                prod_counter[0] += 1

    def consumer_worker() -> None:
        while True:
            item = q.get()
            if item == _SENTINEL:
                q.put(_SENTINEL)
                break
            consumed[item] = 1
            cons_counter[0] += 1

    prod_threads = [
        threading.Thread(target=producer_worker, args=(pid,))
        for pid in range(n_producers)
    ]
    cons_threads = [threading.Thread(target=consumer_worker) for _ in range(n_consumers)]

    for t in prod_threads:
        t.start()
    for t in cons_threads:
        t.start()
    for t in prod_threads:
        t.join(timeout=30)

    assert prod_counter[0] == total, f"Produced {prod_counter[0]} != {total}"

    q.put(_SENTINEL)

    for t in cons_threads:
        t.join(timeout=30)

    assert cons_counter[0] == total, f"Consumed {cons_counter[0]} != {total}"
    assert all(consumed), "Not all items consumed"
