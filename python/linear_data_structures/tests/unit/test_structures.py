from __future__ import annotations

from typing import TypeVar

import pytest

from linear_data_structures.implementation.circular_queue import CircularQueue
from linear_data_structures.implementation.queue import Queue
from linear_data_structures.implementation.stack import Stack

T = TypeVar("T")


class TestStack:
    def test_push_pop(self) -> None:
        s = Stack[int]()
        s.push(1)
        s.push(2)
        assert s.pop() == 2
        assert s.pop() == 1

    def test_peek(self) -> None:
        s = Stack[int]()
        s.push(10)
        assert s.peek() == 10
        assert s.size == 1

    def test_is_empty(self) -> None:
        s = Stack[int]()
        assert s.is_empty
        s.push(1)
        assert not s.is_empty

    def test_size(self) -> None:
        s = Stack[int]()
        assert s.size == 0
        s.push(1)
        s.push(2)
        assert s.size == 2

    def test_pop_empty(self) -> None:
        s = Stack[int]()
        with pytest.raises(IndexError):
            s.pop()

    def test_peek_empty(self) -> None:
        s = Stack[int]()
        with pytest.raises(IndexError):
            s.peek()


class TestQueue:
    def test_enqueue_dequeue(self) -> None:
        q = Queue[int]()
        q.enqueue(1)
        q.enqueue(2)
        assert q.dequeue() == 1
        assert q.dequeue() == 2

    def test_peek(self) -> None:
        q = Queue[int]()
        q.enqueue(42)
        assert q.peek() == 42
        assert q.size == 1

    def test_fifo_order(self) -> None:
        q = Queue[str]()
        for c in "abc":
            q.enqueue(c)
        result = "".join(q.dequeue() for _ in range(q.size))
        assert result == "abc"

    def test_dequeue_empty(self) -> None:
        q = Queue[int]()
        with pytest.raises(IndexError):
            q.dequeue()

    def test_peek_empty(self) -> None:
        q = Queue[int]()
        with pytest.raises(IndexError):
            q.peek()


class TestCircularQueue:
    def test_enqueue_dequeue(self) -> None:
        cq = CircularQueue[int](3)
        cq.enqueue(1)
        cq.enqueue(2)
        assert cq.dequeue() == 1
        assert cq.dequeue() == 2

    def test_wrap_around(self) -> None:
        cq = CircularQueue[int](3)
        cq.enqueue(1)
        cq.enqueue(2)
        cq.enqueue(3)
        cq.dequeue()
        cq.dequeue()
        cq.enqueue(4)
        cq.enqueue(5)
        assert cq.dequeue() == 3
        assert cq.dequeue() == 4
        assert cq.dequeue() == 5

    def test_is_full(self) -> None:
        cq = CircularQueue[int](2)
        assert not cq.is_full
        cq.enqueue(1)
        cq.enqueue(2)
        assert cq.is_full

    def test_is_empty(self) -> None:
        cq = CircularQueue[int](3)
        assert cq.is_empty
        cq.enqueue(1)
        assert not cq.is_empty

    def test_enqueue_full(self) -> None:
        cq = CircularQueue[int](1)
        cq.enqueue(1)
        with pytest.raises(IndexError):
            cq.enqueue(2)

    def test_dequeue_empty(self) -> None:
        cq = CircularQueue[int](3)
        with pytest.raises(IndexError):
            cq.dequeue()

    def test_peek_empty(self) -> None:
        cq = CircularQueue[int](3)
        with pytest.raises(IndexError):
            cq.peek()

    def test_peek(self) -> None:
        cq = CircularQueue[int](3)
        cq.enqueue(99)
        assert cq.peek() == 99
        assert cq.size == 1

    def test_size(self) -> None:
        cq = CircularQueue[int](5)
        assert cq.size == 0
        cq.enqueue(1)
        cq.enqueue(2)
        assert cq.size == 2
