from __future__ import annotations

import threading
from collections import deque
from typing import Generic, TypeVar

T = TypeVar("T")


class BlockingQueue(Generic[T]):
    def __init__(self, capacity: int) -> None:
        if capacity < 1:
            raise ValueError("capacity must be >= 1")
        self._capacity = capacity
        self._queue: deque[T] = deque()
        self._lock = threading.Lock()
        self._not_full = threading.Condition(self._lock)
        self._not_empty = threading.Condition(self._lock)

    def put(self, item: T, timeout: float | None = None) -> bool:
        with self._not_full:
            while self._size() >= self._capacity:
                if not self._not_full.wait(timeout=timeout):
                    return False
            self._queue.append(item)
            self._not_empty.notify()
            return True

    def get(self, timeout: float | None = None) -> T | None:
        with self._not_empty:
            while self._size() == 0:
                if not self._not_empty.wait(timeout=timeout):
                    return None
            item = self._queue.popleft()
            self._not_full.notify()
            return item

    def _size(self) -> int:
        return len(self._queue)

    @property
    def size(self) -> int:
        with self._lock:
            return len(self._queue)

    @property
    def capacity(self) -> int:
        return self._capacity
