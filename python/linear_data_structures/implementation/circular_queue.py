from __future__ import annotations

from typing import Generic, TypeVar

T = TypeVar("T")


class CircularQueue(Generic[T]):
    def __init__(self, capacity: int) -> None:
        if capacity < 1:
            raise ValueError("capacity must be >= 1")
        self._buffer: list[T | None] = [None] * capacity
        self._capacity = capacity
        self._head = 0
        self._tail = 0
        self._count = 0

    def enqueue(self, item: T) -> None:
        if self.is_full:
            raise IndexError("enqueue on full circular queue")
        self._buffer[self._tail] = item
        self._tail = (self._tail + 1) % self._capacity
        self._count += 1

    def dequeue(self) -> T:
        if self.is_empty:
            raise IndexError("dequeue from empty circular queue")
        value = self._buffer[self._head]
        self._buffer[self._head] = None
        self._head = (self._head + 1) % self._capacity
        self._count -= 1
        return value  # type: ignore[return-value]

    def peek(self) -> T:
        if self.is_empty:
            raise IndexError("peek from empty circular queue")
        return self._buffer[self._head]  # type: ignore[return-value]

    @property
    def is_full(self) -> bool:
        return self._count == self._capacity

    @property
    def is_empty(self) -> bool:
        return self._count == 0

    @property
    def size(self) -> int:
        return self._count

    @property
    def capacity(self) -> int:
        return self._capacity
