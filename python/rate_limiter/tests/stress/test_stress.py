import threading
import unittest

from rate_limiter.implementation.sliding_window import SlidingWindowLog
from rate_limiter.implementation.token_bucket import TokenBucket


class TestTokenBucketStress(unittest.TestCase):
    def test_100_concurrent_requesters(self) -> None:
        bucket = TokenBucket(burst_capacity=500, refill_rate=500, refill_interval=1)
        allowed = threading.Event()
        results: list[bool] = []
        lock = threading.Lock()

        def requester() -> None:
            allowed.wait()
            for _ in range(10):
                result = bucket.allow_request()
                with lock:
                    results.append(result)

        threads = [threading.Thread(target=requester) for _ in range(100)]
        for t in threads:
            t.start()

        allowed.set()
        for t in threads:
            t.join()

        allowed_count = sum(1 for r in results if r)
        self.assertLessEqual(allowed_count, 500)


class TestSlidingWindowStress(unittest.TestCase):
    def test_100_concurrent_requesters(self) -> None:
        window = SlidingWindowLog(window_size=1, max_requests=200)
        allowed = threading.Event()
        results: list[bool] = []
        lock = threading.Lock()

        def requester() -> None:
            allowed.wait()
            for _ in range(5):
                result = window.allow_request()
                with lock:
                    results.append(result)

        threads = [threading.Thread(target=requester) for _ in range(100)]
        for t in threads:
            t.start()

        allowed.set()
        for t in threads:
            t.join()

        allowed_count = sum(1 for r in results if r)
        self.assertLessEqual(allowed_count, 200)


if __name__ == "__main__":
    unittest.main()
