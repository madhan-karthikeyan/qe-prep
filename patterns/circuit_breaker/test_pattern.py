import time

import pytest

from circuit_breaker.implementation import CircuitBreaker, CircuitBreakerOpenError, CircuitBreakerState


def test_closed_state_on_init():
    cb = CircuitBreaker(failure_threshold=3, timeout=1.0)
    assert cb.state is CircuitBreakerState.CLOSED


def test_opens_after_threshold_failures():
    cb = CircuitBreaker(failure_threshold=3, timeout=30.0)

    def fail():
        raise ValueError("fail")

    for _ in range(3):
        with pytest.raises(ValueError):
            cb.call(fail)

    assert cb.state is CircuitBreakerState.OPEN


def test_raises_open_error_when_open():
    cb = CircuitBreaker(failure_threshold=1, timeout=30.0)

    def fail():
        raise ValueError("fail")

    with pytest.raises(ValueError):
        cb.call(fail)

    with pytest.raises(CircuitBreakerOpenError):
        cb.call(fail)


def test_half_open_after_timeout():
    cb = CircuitBreaker(failure_threshold=1, timeout=0.1)

    def fail():
        raise ValueError("fail")

    with pytest.raises(ValueError):
        cb.call(fail)
    assert cb.state is CircuitBreakerState.OPEN

    time.sleep(0.15)
    assert cb.state is CircuitBreakerState.HALF_OPEN


def test_closes_on_success_in_half_open():
    cb = CircuitBreaker(failure_threshold=1, success_threshold=2, timeout=0.1)

    def fail():
        raise ValueError("fail")

    with pytest.raises(ValueError):
        cb.call(fail)
    time.sleep(0.15)

    def ok():
        return "ok"

    assert cb.call(ok) == "ok"
    assert cb.state is CircuitBreakerState.HALF_OPEN

    assert cb.call(ok) == "ok"
    assert cb.state is CircuitBreakerState.CLOSED


def test_failure_in_half_open_goes_back_to_open():
    cb = CircuitBreaker(failure_threshold=2, success_threshold=1, timeout=0.1)

    def fail():
        raise ValueError("fail")

    for _ in range(2):
        with pytest.raises(ValueError):
            cb.call(fail)

    assert cb.state is CircuitBreakerState.OPEN
    time.sleep(0.15)
    assert cb.state is CircuitBreakerState.HALF_OPEN

    with pytest.raises(ValueError):
        cb.call(fail)
    assert cb.state is CircuitBreakerState.OPEN


def test_success_resets_failure_count_in_closed():
    cb = CircuitBreaker(failure_threshold=3, timeout=30.0)

    def fail():
        raise ValueError("fail")

    for _ in range(2):
        with pytest.raises(ValueError):
            cb.call(fail)

    def ok():
        return "ok"

    cb.call(ok)
    with pytest.raises(ValueError):
        cb.call(fail)
    with pytest.raises(ValueError):
        cb.call(fail)

    assert cb.state is CircuitBreakerState.CLOSED
