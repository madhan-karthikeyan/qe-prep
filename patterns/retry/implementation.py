import random
import time
from typing import Callable, Collection, Optional, Type


class RetryError(Exception):
    """Raised when all retry attempts are exhausted."""


class Retry:
    """Configurable retry with exponential backoff, jitter, and retryable exceptions.

    Args:
        max_attempts: Maximum number of retry attempts.
        base_delay: Initial delay in seconds.
        max_delay: Maximum delay cap in seconds.
        factor: Exponential factor (default 2.0).
        retryable_exceptions: Collection of exception types to retry on.
        on_retry: Optional callback(retry_state) called before each retry.
        jitter: If True, apply full jitter (randomize delay).
    """

    def __init__(
        self,
        max_attempts: int = 3,
        base_delay: float = 1.0,
        max_delay: float = 60.0,
        factor: float = 2.0,
        retryable_exceptions: Optional[Collection[Type[Exception]]] = None,
        on_retry: Optional[Callable[..., None]] = None,
        jitter: bool = True,
    ) -> None:
        self.max_attempts = max_attempts
        self.base_delay = base_delay
        self.max_delay = max_delay
        self.factor = factor
        self.retryable_exceptions = retryable_exceptions or (Exception,)
        self.on_retry = on_retry
        self.jitter = jitter

    def call(self, func: Callable, *args, **kwargs):
        """Execute *func* with retry logic.

        Raises RetryError if all attempts fail.
        """
        last_exception: Optional[Exception] = None
        for attempt in range(1, self.max_attempts + 1):
            try:
                return func(*args, **kwargs)
            except self.retryable_exceptions as exc:
                last_exception = exc
                if attempt == self.max_attempts:
                    break
                delay = self._compute_delay(attempt)
                if self.on_retry:
                    self.on_retry(attempt=attempt, delay=delay, exception=exc)
                time.sleep(delay)
        raise RetryError(f"Failed after {self.max_attempts} attempts") from last_exception

    def _compute_delay(self, attempt: int) -> float:
        delay = min(self.base_delay * (self.factor ** (attempt - 1)), self.max_delay)
        if self.jitter:
            delay = random.uniform(0, delay)
        return delay
