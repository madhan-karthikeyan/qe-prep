import queue
import threading

import pytest

from object_pool.implementation import ObjectPool


def test_acquire_release():
    pool = ObjectPool(factory=lambda: [], max_size=2, timeout=1.0)

    obj = pool.acquire()
    obj.append(1)
    pool.release(obj)

    obj2 = pool.acquire()
    assert obj2 == [1]
    pool.release(obj2)


def test_timeout_on_empty_pool():
    pool = ObjectPool(factory=lambda: [], max_size=1, timeout=0.1)

    obj = pool.acquire()
    with pytest.raises(queue.Empty):
        pool.acquire()

    pool.release(obj)


def test_reuse_objects():
    pool = ObjectPool(factory=lambda: [], max_size=2, timeout=1.0)
    objects = set()

    # Acquire both objects (pool grows to max_size)
    obj1 = pool.acquire()
    obj2 = pool.acquire()
    objects.update([id(obj1), id(obj2)])

    pool.release(obj1)
    pool.release(obj2)

    # Acquire again — should get back the same objects
    obj3 = pool.acquire()
    obj4 = pool.acquire()
    objects.update([id(obj3), id(obj4)])

    # Pool only ever created 2 objects
    assert len(objects) == 2


def test_validator_discards_invalid():
    def is_valid(obj):
        return obj.get("valid", False)

    pool = ObjectPool(factory=lambda: {"valid": True}, validator=is_valid, max_size=2)

    obj = pool.acquire()
    obj["valid"] = False
    pool.release(obj)  # discarded by validator

    # Pool should create a fresh object on next acquire
    obj2 = pool.acquire()
    assert obj2["valid"] is True
    pool.release(obj2)


def test_concurrent_access():
    pool = ObjectPool(factory=lambda: [], max_size=5, timeout=5.0)
    errors: list[Exception] = []
    lock = threading.Lock()

    def worker():
        for _ in range(10):
            try:
                obj = pool.acquire()
                pool.release(obj)
            except Exception as e:
                with lock:
                    errors.append(e)

    threads = [threading.Thread(target=worker) for _ in range(5)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    assert not errors
    assert pool.size <= 5


def test_acquire_creates_up_to_max_size():
    pool = ObjectPool(factory=lambda: [], max_size=3, timeout=0.1)
    objs = [pool.acquire() for _ in range(3)]

    assert pool.size == 3

    with pytest.raises(queue.Empty):
        pool.acquire()

    for obj in objs:
        pool.release(obj)
