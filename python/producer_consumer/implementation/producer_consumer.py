from __future__ import annotations

import threading

from producer_consumer.implementation.blocking_queue import BlockingQueue


def producer(
    q: BlockingQueue[int],
    start: int,
    count: int,
    results: list[int],
) -> None:
    for i in range(start, start + count):
        q.put(i)
        results.append(i)


def consumer(
    q: BlockingQueue[int],
    count: int,
    results: list[int],
) -> None:
    for _ in range(count):
        item = q.get()
        if item is not None:
            results.append(item)


def run_demo(
    num_producers: int = 3,
    num_consumers: int = 3,
    items_per_producer: int = 10,
    queue_capacity: int = 5,
) -> None:
    q: BlockingQueue[int] = BlockingQueue(queue_capacity)
    produced: list[int] = []
    consumed: list[int] = []
    produced_lock = threading.Lock()
    consumed_lock = threading.Lock()

    def safe_producer(start: int, count: int) -> None:
        local: list[int] = []
        for i in range(start, start + count):
            q.put(i)
            local.append(i)
        with produced_lock:
            produced.extend(local)

    def safe_consumer(count: int) -> None:
        local: list[int] = []
        for _ in range(count):
            item = q.get()
            if item is not None:
                local.append(item)
        with consumed_lock:
            consumed.extend(local)

    total_items = num_producers * items_per_producer

    producers = [
        threading.Thread(
            target=safe_producer,
            args=(p * items_per_producer, items_per_producer),
        )
        for p in range(num_producers)
    ]
    consumers = [
        threading.Thread(target=safe_consumer, args=(total_items // num_consumers,))
        for _ in range(num_consumers)
    ]

    for t in producers:
        t.start()
    for t in consumers:
        t.start()
    for t in producers:
        t.join()
    for t in consumers:
        t.join()

    assert sorted(produced) == sorted(consumed), (
        f"Produced {sorted(produced)} != Consumed {sorted(consumed)}"
    )
    print(f"Demo: produced {len(produced)} items, consumed {len(consumed)} items.")


if __name__ == "__main__":
    run_demo()
