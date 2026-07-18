from __future__ import annotations

import threading
from typing import Generic, TypeVar, cast

KT = TypeVar("KT")
VT = TypeVar("VT")


class _Node(Generic[KT, VT]):
    __slots__ = ("key", "value", "prev", "next")

    def __init__(self, key: KT, value: VT) -> None:
        self.key = key
        self.value = value
        self.prev: _Node[KT, VT] | None = None
        self.next: _Node[KT, VT] | None = None


class LRUCache(Generic[KT, VT]):
    def __init__(self, capacity: int, thread_safe: bool = False) -> None:
        if capacity < 1:
            raise ValueError("capacity must be >= 1")
        self._capacity = capacity
        self._cache: dict[KT, _Node[KT, VT]] = {}
        self._head = _Node(cast(KT, None), cast(VT, None))
        self._tail = _Node(cast(KT, None), cast(VT, None))
        self._head.next = self._tail
        self._tail.prev = self._head
        self._lock = threading.Lock() if thread_safe else None

    def _acquire(self) -> None:
        if self._lock is not None:
            self._lock.acquire()

    def _release(self) -> None:
        if self._lock is not None:
            self._lock.release()

    def _remove_node(self, node: _Node[KT, VT]) -> None:
        prev_node = node.prev
        next_node = node.next
        if prev_node is not None:
            prev_node.next = next_node
        if next_node is not None:
            next_node.prev = prev_node

    def _add_to_front(self, node: _Node[KT, VT]) -> None:
        node.prev = self._head
        node.next = self._head.next
        if self._head.next is not None:
            self._head.next.prev = node
        self._head.next = node

    def get(self, key: KT) -> VT | None:
        self._acquire()
        try:
            node = self._cache.get(key)
            if node is None:
                return None
            self._remove_node(node)
            self._add_to_front(node)
            return node.value
        finally:
            self._release()

    def put(self, key: KT, value: VT) -> None:
        self._acquire()
        try:
            if key in self._cache:
                node = self._cache[key]
                node.value = value
                self._remove_node(node)
                self._add_to_front(node)
                return
            if len(self._cache) >= self._capacity:
                lru = self._tail.prev
                if lru is not None and lru is not self._head:
                    self._remove_node(lru)
                    del self._cache[lru.key]
            new_node = _Node(key, value)
            self._cache[key] = new_node
            self._add_to_front(new_node)
        finally:
            self._release()

    @property
    def capacity(self) -> int:
        return self._capacity

    def __len__(self) -> int:
        self._acquire()
        try:
            return len(self._cache)
        finally:
            self._release()
