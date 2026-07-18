import threading
import time


class RateLimiter:
    """Token bucket rate limiter, thread-safe.

    Args:
        rate: Tokens added per second (steady-state request rate).
        burst: Maximum accumulated tokens (burst capacity).
    """

    def __init__(self, rate: float, burst: float) -> None:
        if rate <= 0 or burst <= 0:
            raise ValueError("rate and burst must be positive")
        self._rate = rate
        self._burst = burst
        self._tokens = burst
        self._last_refill = time.monotonic()
        self._lock = threading.Lock()

    def allow_request(self) -> bool:
        """Check if a request is allowed, consuming one token if so.

        Returns True if the request is permitted, False otherwise.
        """
        with self._lock:
            self._refill()
            if self._tokens >= 1.0:
                self._tokens -= 1.0
                return True
            return False

    def _refill(self) -> None:
        now = time.monotonic()
        elapsed = now - self._last_refill
        self._tokens = min(self._burst, self._tokens + elapsed * self._rate)
        self._last_refill = now
