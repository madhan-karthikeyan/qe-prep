import queue
import threading

import pytest

from producer_consumer.implementation import ProducerConsumer


def test_all_items_consumed():
    items = list(range(20))
    consumed: list[int] = []
    lock = threading.Lock()

    def producer(q: queue.Queue):
        for item in items:
            q.put(item)

    def consumer(item: int):
        with lock:
            consumed.append(item)

    pc = ProducerConsumer(num_producers=1, num_consumers=2)
    pc.run(producer, consumer)

    assert sorted(consumed) == items


def test_multiple_producers():
    items = list(range(100))
    produced: list[int] = []
    consumed: list[int] = []
    prod_lock = threading.Lock()
    cons_lock = threading.Lock()

    def producer(q: queue.Queue):
        while True:
            with prod_lock:
                if not items:
                    break
                item = items.pop(0)
            q.put(item)

    def consumer(item: int):
        with cons_lock:
            consumed.append(item)

    pc = ProducerConsumer(num_producers=2, num_consumers=3)
    pc.run(producer, consumer)

    assert sorted(consumed) == list(range(100))


def test_shutdown_consumes_no_extra_items():
    consumed: list[int] = []
    lock = threading.Lock()

    def producer(q: queue.Queue):
        for i in range(10):
            q.put(i)

    def consumer(item: int):
        with lock:
            consumed.append(item)

    pc = ProducerConsumer(num_producers=1, num_consumers=2)
    pc.run(producer, consumer)

    assert sorted(consumed) == list(range(10))
    assert pc._queue.empty()


def test_empty_producer():
    consumed: list[int] = []

    def producer(q: queue.Queue):
        pass

    def consumer(item: int):
        consumed.append(item)

    pc = ProducerConsumer(num_producers=1, num_consumers=1)
    pc.run(producer, consumer)

    assert consumed == []


def test_producer_exception_does_not_block():
    consumed: list[int] = []

    def producer(q: queue.Queue):
        raise RuntimeError("unexpected")
        q.put(1)

    def consumer(item: int):
        consumed.append(item)

    pc = ProducerConsumer(num_producers=1, num_consumers=1)
    pc.run(producer, consumer)

    assert consumed == []
