import threading
import time

import pytest

from rate_limiter.implementation import RateLimiter


def test_rate_enforcement():
    limiter = RateLimiter(rate=10, burst=10)

    allowed = sum(limiter.allow_request() for _ in range(20))
    assert allowed == 10  # burst consumed


def test_burst_behavior():
    limiter = RateLimiter(rate=100, burst=5)

    allowed = sum(limiter.allow_request() for _ in range(10))
    assert allowed == 5


def test_refill_over_time():
    limiter = RateLimiter(rate=10, burst=5)

    for _ in range(5):
        limiter.allow_request()

    assert not limiter.allow_request()

    time.sleep(0.15)  # ~1.5 tokens should refill
    assert limiter.allow_request()


def test_invalid_params():
    with pytest.raises(ValueError):
        RateLimiter(rate=0, burst=5)
    with pytest.raises(ValueError):
        RateLimiter(rate=5, burst=0)
    with pytest.raises(ValueError):
        RateLimiter(rate=-1, burst=5)


def test_concurrent_access():
    limiter = RateLimiter(rate=1000, burst=100)
    allowed_count = 0
    lock = threading.Lock()

    def worker():
        nonlocal allowed_count
        for _ in range(10):
            if limiter.allow_request():
                with lock:
                    allowed_count += 1

    threads = [threading.Thread(target=worker) for _ in range(10)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    assert allowed_count <= 100


def test_no_negative_tokens():
    limiter = RateLimiter(rate=1, burst=1)
    limiter.allow_request()
    assert not limiter.allow_request()
    assert not limiter.allow_request()  # still no tokens
