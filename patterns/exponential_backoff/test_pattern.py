import random

import pytest

from exponential_backoff.implementation import (
    BackoffCalculator,
    decorrelated_jitter,
    equal_jitter,
    full_jitter,
)


class TestJitterFunctions:
    def test_full_jitter_range(self):
        for _ in range(100):
            d = full_jitter(10.0)
            assert 0 <= d < 10.0

    def test_equal_jitter_range(self):
        for _ in range(100):
            d = equal_jitter(10.0)
            assert 5.0 <= d <= 10.0

    def test_decorrelated_jitter(self):
        for _ in range(100):
            d = decorrelated_jitter(10.0, previous_delay=1.0)
            assert d > 0


class TestBackoffCalculator:
    def test_delay_values_no_jitter(self):
        calc = BackoffCalculator(
            base_delay=1.0, multiplier=2.0, max_delay=60.0,
            jitter_fn=lambda d: d,
        )
        assert calc.delay(1) == 1.0
        assert calc.delay(2) == 2.0
        assert calc.delay(3) == 4.0
        assert calc.delay(4) == 8.0

    def test_max_delay_cap(self):
        calc = BackoffCalculator(
            base_delay=1.0, multiplier=10.0, max_delay=50.0,
            jitter_fn=lambda d: d,
        )
        assert calc.delay(1) == 1.0
        assert calc.delay(2) == 10.0
        assert calc.delay(3) == 50.0
        assert calc.delay(4) == 50.0

    def test_jitter_range_with_full_jitter(self):
        random.seed(42)
        calc = BackoffCalculator(
            base_delay=10.0, multiplier=1.0, max_delay=10.0,
            jitter_fn=full_jitter,
        )
        values = [calc.delay(1) for _ in range(100)]
        assert all(0 <= v <= 10.0 for v in values)
        assert any(v < 9.0 for v in values)

    def test_custom_base_delay(self):
        calc = BackoffCalculator(
            base_delay=0.5, multiplier=3.0, max_delay=100.0,
            jitter_fn=lambda d: d,
        )
        assert calc.delay(1) == 0.5
        assert calc.delay(2) == 1.5
        assert calc.delay(3) == 4.5
