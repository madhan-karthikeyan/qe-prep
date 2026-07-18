import time
from collections import deque
from threading import Lock


class SlidingWindowLog:
    def __init__(self, window_size: float, max_requests: int) -> None:
        if window_size <= 0:
            raise ValueError("window_size must be positive")
        if max_requests <= 0:
            raise ValueError("max_requests must be positive")

        self._window_size = window_size
        self._max_requests = max_requests
        self._timestamps: deque[float] = deque()
        self._lock = Lock()

    def allow_request(self) -> bool:
        now = time.monotonic()
        with self._lock:
            cutoff = now - self._window_size
            while self._timestamps and self._timestamps[0] < cutoff:
                self._timestamps.popleft()

            if len(self._timestamps) < self._max_requests:
                self._timestamps.append(now)
                return True
            return False

    @property
    def request_count(self) -> int:
        now = time.monotonic()
        cutoff = now - self._window_size
        with self._lock:
            while self._timestamps and self._timestamps[0] < cutoff:
                self._timestamps.popleft()
            return len(self._timestamps)

    def reset(self) -> None:
        with self._lock:
            self._timestamps.clear()
