import time
from unittest.mock import Mock

import pytest

from retry.implementation import Retry, RetryError


def test_retry_success_on_first_attempt():
    func = Mock(return_value=42)
    retry = Retry(max_attempts=3)
    assert retry.call(func) == 42
    assert func.call_count == 1


def test_retry_success_on_nth_attempt():
    call_count = 0

    def flaky():
        nonlocal call_count
        call_count += 1
        if call_count < 3:
            raise ValueError("not yet")
        return "ok"

    retry = Retry(max_attempts=5, base_delay=0.01, jitter=False)
    assert retry.call(flaky) == "ok"
    assert call_count == 3


def test_max_retries_exceeded():
    def always_fails():
        raise ValueError("boom")

    retry = Retry(max_attempts=3, base_delay=0.01, jitter=False)
    with pytest.raises(RetryError):
        retry.call(always_fails)


def test_on_retry_callback():
    callback = Mock()
    call_count = 0

    def flaky():
        nonlocal call_count
        call_count += 1
        if call_count < 3:
            raise ValueError("not yet")
        return "ok"

    retry = Retry(
        max_attempts=5, base_delay=0.01, jitter=False, on_retry=callback
    )
    retry.call(flaky)
    assert callback.call_count == 2


def test_jitter_range():
    delays = []
    original_sleep = time.sleep
    time.sleep = lambda d: delays.append(d)

    call_count = 0

    def flaky():
        nonlocal call_count
        call_count += 1
        if call_count < 5:
            raise ValueError("not yet")
        return "ok"

    try:
        retry = Retry(max_attempts=5, base_delay=1.0, factor=2.0, jitter=True)
        retry.call(flaky)
    finally:
        time.sleep = original_sleep

    for d in delays:
        assert 0 <= d <= 8.0  # max delay with base=1, factor=2, 4 attempts = 8


def test_non_retryable_exception():
    def raises_type_error():
        raise TypeError("not retryable")

    retry = Retry(
        max_attempts=3,
        retryable_exceptions=(ValueError,),
        base_delay=0.01,
    )
    with pytest.raises(TypeError):
        retry.call(raises_type_error)
