import random
from typing import Callable


JitterFn = Callable[[float], float]


def full_jitter(delay: float) -> float:
    """Return a random value in [0, delay)."""
    return random.uniform(0, delay)


def equal_jitter(delay: float) -> float:
    """Return delay/2 + random(0, delay/2) for a spread around midpoint."""
    half = delay / 2
    return half + random.uniform(0, half)


def decorrelated_jitter(delay: float, previous_delay: float = 1.0) -> float:
    """Return random(previous_delay * 1, delay * 3) for smoother backoff."""
    low = previous_delay * 1.0
    high = delay * 3.0
    if low >= high:
        return delay * 2.0
    return random.uniform(low, high)


class BackoffCalculator:
    """Compute backoff delays with configurable strategy.

    Args:
        base_delay: Initial delay in seconds.
        multiplier: Exponential factor (default 2.0).
        max_delay: Maximum delay cap.
        jitter_fn: Jitter function (default full_jitter).
    """

    def __init__(
        self,
        base_delay: float = 1.0,
        multiplier: float = 2.0,
        max_delay: float = 60.0,
        jitter_fn: JitterFn = full_jitter,
    ) -> None:
        self.base_delay = base_delay
        self.multiplier = multiplier
        self.max_delay = max_delay
        self.jitter_fn = jitter_fn

    def delay(self, attempt: int) -> float:
        """Return the delay in seconds for the given attempt (1-indexed)."""
        raw = self.base_delay * (self.multiplier ** (attempt - 1))
        capped = min(raw, self.max_delay)
        return self.jitter_fn(capped)
