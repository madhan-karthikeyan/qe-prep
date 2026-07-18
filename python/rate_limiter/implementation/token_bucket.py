import time
from threading import Lock


class TokenBucket:
    def __init__(
        self,
        burst_capacity: int,
        refill_rate: float,
        refill_interval: float = 1.0,
    ) -> None:
        if burst_capacity <= 0:
            raise ValueError("burst_capacity must be positive")
        if refill_rate <= 0:
            raise ValueError("refill_rate must be positive")
        if refill_interval <= 0:
            raise ValueError("refill_interval must be positive")

        self._capacity = burst_capacity
        self._refill_rate = refill_rate
        self._refill_interval = refill_interval
        self._tokens: float = float(burst_capacity)
        self._last_refill = time.monotonic()
        self._lock = Lock()

    def allow_request(self, tokens: int = 1) -> bool:
        if tokens <= 0:
            raise ValueError("tokens must be positive")

        with self._lock:
            self._refill()
            if self._tokens >= tokens:
                self._tokens -= tokens
                return True
            return False

    def _refill(self) -> None:
        now = time.monotonic()
        elapsed = now - self._last_refill
        if elapsed >= self._refill_interval:
            intervals_passed = elapsed / self._refill_interval
            added = intervals_passed * self._refill_rate
            self._tokens = min(self._capacity, self._tokens + added)
            self._last_refill = now

    @property
    def available_tokens(self) -> float:
        with self._lock:
            self._refill()
            return self._tokens

    @property
    def capacity(self) -> int:
        return self._capacity

    def reset(self) -> None:
        with self._lock:
            self._tokens = self._capacity
            self._last_refill = time.monotonic()
