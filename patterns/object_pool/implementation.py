import queue
import threading
import time
from typing import Any, Callable, Optional, TypeVar

T = TypeVar("T")


class ObjectPool:
    """Thread-safe pool of reusable objects.

    Args:
        factory: Callable that creates new objects.
        validator: Optional callable to validate objects on release.
        max_size: Maximum number of objects in the pool.
        timeout: Seconds to wait when blocking on acquire.
    """

    def __init__(
        self,
        factory: Callable[[], T],
        validator: Optional[Callable[[T], bool]] = None,
        max_size: int = 10,
        timeout: float = 5.0,
    ) -> None:
        self._factory = factory
        self._validator = validator
        self._max_size = max_size
        self._timeout = timeout
        self._pool: queue.Queue = queue.Queue()
        self._size = 0
        self._lock = threading.Lock()

    def acquire(self) -> T:
        """Acquire an object from the pool, blocking if empty.

        Raises queue.Empty if acquisition times out.
        """
        try:
            obj = self._pool.get(block=True, timeout=self._timeout)
            return obj
        except queue.Empty:
            with self._lock:
                if self._size < self._max_size:
                    obj = self._factory()
                    self._size += 1
                    return obj
            raise

    def release(self, obj: T) -> None:
        """Return an object to the pool.

        The object is validated (if validator is set) before being returned.
        Invalid objects are discarded.
        """
        if self._validator is not None and not self._validator(obj):
            with self._lock:
                self._size -= 1
            return

        try:
            self._pool.put(obj, block=False)
        except queue.Full:
            with self._lock:
                self._size -= 1

    @property
    def size(self) -> int:
        return self._size
