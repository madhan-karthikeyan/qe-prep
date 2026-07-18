import threading
import time
from enum import Enum, auto
from typing import Any, Callable, Optional


class CircuitBreakerState(Enum):
    CLOSED = auto()
    OPEN = auto()
    HALF_OPEN = auto()


class CircuitBreakerOpenError(Exception):
    """Raised when the circuit breaker is OPEN and rejects a call."""


class CircuitBreaker:
    """Thread-safe circuit breaker with automatic state transitions.

    Args:
        failure_threshold: Failures before circuit opens.
        success_threshold: Successes in HALF_OPEN to close.
        timeout: Seconds before OPEN transitions to HALF_OPEN.
    """

    def __init__(
        self,
        failure_threshold: int = 5,
        success_threshold: int = 2,
        timeout: float = 30.0,
    ) -> None:
        self._failure_threshold = failure_threshold
        self._success_threshold = success_threshold
        self._timeout = timeout

        self._state = CircuitBreakerState.CLOSED
        self._failure_count = 0
        self._success_count = 0
        self._last_failure_time: float = 0.0
        self._lock = threading.Lock()

    @property
    def state(self) -> CircuitBreakerState:
        self._try_transition()
        return self._state

    def call(self, func: Callable[..., Any], *args, **kwargs) -> Any:
        """Execute *func* through the circuit breaker.

        Raises CircuitBreakerOpenError when the circuit is OPEN.
        """
        self._try_transition()

        with self._lock:
            if self._state is CircuitBreakerState.OPEN:
                raise CircuitBreakerOpenError("Circuit breaker is OPEN")

        try:
            result = func(*args, **kwargs)
        except Exception as exc:
            self._record_failure()
            raise exc

        self._record_success()
        return result

    def _try_transition(self) -> None:
        with self._lock:
            if (
                self._state is CircuitBreakerState.OPEN
                and time.monotonic() - self._last_failure_time >= self._timeout
            ):
                self._state = CircuitBreakerState.HALF_OPEN
                self._success_count = 0

    def _record_failure(self) -> None:
        with self._lock:
            self._last_failure_time = time.monotonic()
            if self._state is CircuitBreakerState.HALF_OPEN:
                self._state = CircuitBreakerState.OPEN
                self._failure_count = 0
                return
            self._failure_count += 1
            if self._failure_count >= self._failure_threshold:
                self._state = CircuitBreakerState.OPEN
                self._failure_count = 0

    def _record_success(self) -> None:
        with self._lock:
            if self._state is CircuitBreakerState.HALF_OPEN:
                self._success_count += 1
                if self._success_count >= self._success_threshold:
                    self._state = CircuitBreakerState.CLOSED
                    self._failure_count = 0
                    self._success_count = 0
            elif self._state is CircuitBreakerState.CLOSED:
                self._failure_count = 0
