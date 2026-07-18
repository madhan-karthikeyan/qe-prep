from __future__ import annotations

from collections import deque
from typing import Generic, TypeVar

T = TypeVar("T")


class Queue(Generic[T]):
    def __init__(self) -> None:
        self._items: deque[T] = deque()

    def enqueue(self, item: T) -> None:
        self._items.append(item)

    def dequeue(self) -> T:
        if self.is_empty:
            raise IndexError("dequeue from empty queue")
        return self._items.popleft()

    def peek(self) -> T:
        if self.is_empty:
            raise IndexError("peek from empty queue")
        return self._items[0]

    @property
    def is_empty(self) -> bool:
        return len(self._items) == 0

    @property
    def size(self) -> int:
        return len(self._items)
