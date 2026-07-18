import time
import unittest

from rate_limiter.implementation.sliding_window import SlidingWindowLog
from rate_limiter.implementation.token_bucket import TokenBucket


class TestTokenBucket(unittest.TestCase):
    def test_allow_request_under_capacity(self) -> None:
        bucket = TokenBucket(burst_capacity=10, refill_rate=5, refill_interval=1)
        for _ in range(10):
            self.assertTrue(bucket.allow_request())

    def test_deny_when_exhausted(self) -> None:
        bucket = TokenBucket(burst_capacity=3, refill_rate=1, refill_interval=60)
        for _ in range(3):
            self.assertTrue(bucket.allow_request())
        self.assertFalse(bucket.allow_request())

    def test_refill_over_time(self) -> None:
        bucket = TokenBucket(burst_capacity=5, refill_rate=5, refill_interval=0.01)
        for _ in range(5):
            bucket.allow_request()
        self.assertFalse(bucket.allow_request())
        time.sleep(0.02)
        self.assertTrue(bucket.allow_request())

    def test_consume_multiple_tokens(self) -> None:
        bucket = TokenBucket(burst_capacity=10, refill_rate=5, refill_interval=1)
        self.assertTrue(bucket.allow_request(5))
        self.assertTrue(bucket.allow_request(5))
        self.assertFalse(bucket.allow_request(1))

    def test_invalid_construction(self) -> None:
        with self.assertRaises(ValueError):
            TokenBucket(burst_capacity=0, refill_rate=1, refill_interval=1)
        with self.assertRaises(ValueError):
            TokenBucket(burst_capacity=1, refill_rate=0, refill_interval=1)
        with self.assertRaises(ValueError):
            TokenBucket(burst_capacity=1, refill_rate=1, refill_interval=0)

    def test_invalid_tokens_argument(self) -> None:
        bucket = TokenBucket(burst_capacity=10, refill_rate=5, refill_interval=1)
        with self.assertRaises(ValueError):
            bucket.allow_request(0)

    def test_available_tokens_property(self) -> None:
        bucket = TokenBucket(burst_capacity=10, refill_rate=5, refill_interval=1)
        self.assertEqual(bucket.available_tokens, 10)
        bucket.allow_request(3)
        self.assertEqual(bucket.available_tokens, 7)

    def test_reset(self) -> None:
        bucket = TokenBucket(burst_capacity=5, refill_rate=1, refill_interval=60)
        bucket.allow_request(5)
        self.assertEqual(bucket.available_tokens, 0)
        bucket.reset()
        self.assertEqual(bucket.available_tokens, 5)

    def test_burst_then_refill_partial(self) -> None:
        bucket = TokenBucket(burst_capacity=100, refill_rate=100, refill_interval=0.05)
        for _ in range(100):
            self.assertTrue(bucket.allow_request())
        self.assertFalse(bucket.allow_request())
        time.sleep(0.03)
        self.assertFalse(bucket.allow_request(100))
        time.sleep(0.03)
        self.assertTrue(bucket.allow_request())


class TestSlidingWindowLog(unittest.TestCase):
    def test_allow_under_limit(self) -> None:
        window = SlidingWindowLog(window_size=10, max_requests=5)
        for _ in range(5):
            self.assertTrue(window.allow_request())

    def test_deny_when_exceeded(self) -> None:
        window = SlidingWindowLog(window_size=10, max_requests=3)
        for _ in range(3):
            self.assertTrue(window.allow_request())
        self.assertFalse(window.allow_request())

    def test_window_slides(self) -> None:
        window = SlidingWindowLog(window_size=0.05, max_requests=2)
        self.assertTrue(window.allow_request())
        self.assertTrue(window.allow_request())
        self.assertFalse(window.allow_request())
        time.sleep(0.06)
        self.assertTrue(window.allow_request())

    def test_request_count_property(self) -> None:
        window = SlidingWindowLog(window_size=10, max_requests=10)
        self.assertEqual(window.request_count, 0)
        window.allow_request()
        self.assertEqual(window.request_count, 1)
        window.allow_request()
        self.assertEqual(window.request_count, 2)

    def test_reset(self) -> None:
        window = SlidingWindowLog(window_size=10, max_requests=5)
        window.allow_request()
        window.allow_request()
        window.reset()
        self.assertEqual(window.request_count, 0)

    def test_invalid_construction(self) -> None:
        with self.assertRaises(ValueError):
            SlidingWindowLog(window_size=0, max_requests=5)
        with self.assertRaises(ValueError):
            SlidingWindowLog(window_size=10, max_requests=0)

    def test_stale_timestamps_removed(self) -> None:
        window = SlidingWindowLog(window_size=0.02, max_requests=5)
        window.allow_request()
        time.sleep(0.03)
        self.assertEqual(window.request_count, 0)


if __name__ == "__main__":
    unittest.main()
