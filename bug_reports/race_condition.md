# Race condition in rate limiter token refill

**Severity:** Critical
**Priority:** P0
**Environment:** All platforms, all languages, concurrent request rate > 1000/s
**Component:** API Gateway — Rate Limiter middleware

## Summary

Under concurrent load, the token bucket rate limiter allows up to 2× the configured request rate. A non-atomic check-and-refill operation lets multiple requests pass before the token count is decremented, defeating the rate limit entirely.

## Steps to Reproduce

1. Configure rate limiter: 100 requests/second (burst = 100)
2. Send 200 concurrent requests
3. Count how many succeed vs how many are rate-limited
4. Repeat 10 times

```python
import concurrent.futures
import requests

def send_request(i):
    resp = requests.get("http://localhost:8080/api/endpoint")
    return resp.status_code

with concurrent.futures.ThreadPoolExecutor(max_workers=50) as pool:
    results = list(pool.map(send_request, range(200)))

success = sum(1 for r in results if r == 200)
limited = sum(1 for r in results if r == 429)

print(f"Success: {success}, Rate-limited: {limited}")
# Expected: Success: ~100, Rate-limited: ~100
# Actual:   Success: ~180-200, Rate-limited: ~0-20
```

## Expected Behavior

No more than 100 requests should succeed in any 1-second window. Excess requests return HTTP 429.

## Actual Behavior

- 180-200 of 200 concurrent requests succeed
- The rate limiter allows nearly double the configured rate
- Behavior is non-deterministic (varies per run from 150-200 successes)
- Under sustained load, the effective rate can exceed 2× configured limit

## Logs / Screenshots

```
Rate limiter metrics (1s buckets):
Time 00:00  allowed=180  blocked=20  limit=100  over_by=80%
Time 00:01  allowed=195  blocked=5   limit=100  over_by=95%
Time 00:02  allowed=172  blocked=28  limit=100  over_by=72%

Thread dump during burst:
- Thread-12: executing refill() — computing new tokens
- Thread-13: executing allowRequest() — checking tokens > 0 → true
- Thread-14: executing allowRequest() — checking tokens > 0 → true (both see tokens > 0)
```

## Root Cause Analysis

The token bucket rate limiter has a **check-then-act** race condition in `allowRequest()`:

```python
# Problematic code — NOT thread-safe
class TokenBucketRateLimiter:
    def __init__(self, rate, burst):
        self.tokens = burst
        self.rate = rate
        self.last_refill = time.monotonic()
    
    def allow_request(self):
        self.refill()                             # Step 1: refill tokens
        if self.tokens > 0:                       # Step 2: check
            self.tokens -= 1                      # Step 3: consume (RACE!)
            return True
        return False
    
    def refill(self):
        now = time.monotonic()
        elapsed = now - self.last_refill
        new_tokens = elapsed * self.rate
        if new_tokens > 0:
            self.tokens = min(self.tokens + new_tokens, self.burst)
            self.last_refill = now
```

Race scenario:
1. Thread A calls `allow_request()`, `refill()` sets `tokens = 100`, then checks `tokens > 0` → True
2. Thread B calls `allow_request()` **before Thread A decrements**. Also sees `tokens = 100` → True
3. Both threads decrement: `tokens = 98` (instead of expected -1 or 99)
4. Repeat with many threads — all see `tokens > 0` before any decrement completes

## Fix

Protect the entire check-and-refill operation with a mutex:

```python
import threading

class ThreadSafeTokenBucketRateLimiter:
    def __init__(self, rate, burst):
        self.tokens = burst
        self.rate = rate
        self.burst = burst
        self.last_refill = time.monotonic()
        self.lock = threading.Lock()
    
    def allow_request(self):
        with self.lock:                      # Atomic: refill + check + consume
            self.refill()
            if self.tokens >= 1:
                self.tokens -= 1
                return True
            return False
    
    def refill(self):
        now = time.monotonic()
        elapsed = now - self.last_refill
        new_tokens = elapsed * self.rate
        if new_tokens > 0:
            self.tokens = min(self.tokens + new_tokens, self.burst)
            self.last_refill = now
```

For high-throughput scenarios (100k+ req/s), consider lock-free alternatives:
- **Atomic integers** with CAS (compare-and-swap): `token_count.decrementAndGet()` with retry
- **Striped locking**: Per-key locks for per-user rate limiting
- **Redis-based**: `INCR` + `EXPIRE` with Lua scripting for atomicity

## Regression Tests

### 1. Concurrent Request Stress Test

```python
def test_rate_limiter_enforces_limit_under_concurrent_load():
    limiter = ThreadSafeTokenBucketRateLimiter(rate=100, burst=100)
    num_threads = 50
    requests_per_thread = 10
    allowed = threading.Event()
    allowed_counter = 0
    lock = threading.Lock()
    
    def worker():
        nonlocal allowed_counter
        for _ in range(requests_per_thread):
            if limiter.allow_request():
                with lock:
                    allowed_counter += 1
    
    threads = [threading.Thread(target=worker) for _ in range(num_threads)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()
    
    # Window is 1s, rate is 100/s, so we should not exceed 100 + small tolerance
    assert allowed_counter <= 105, f"Allowed {allowed_counter} requests (limit: 100)"
```

### 2. Rate Verification Over Time

```python
def test_rate_limiter_averages_correct_rate_over_multiple_windows():
    limiter = ThreadSafeTokenBucketRateLimiter(rate=50, burst=50)
    results = []
    
    for second in range(10):
        count = 0
        deadline = time.monotonic() + 1.0
        while time.monotonic() < deadline:
            if limiter.allow_request():
                count += 1
        results.append(count)
    
    avg_rate = sum(results) / len(results)
    assert 45 <= avg_rate <= 55, f"Average rate {avg_rate:.1f} not within 10% of 50"
```

### 3. Burst Behavior Test

```python
def test_burst_allows_immediate_burst_then_enforces_limit():
    limiter = ThreadSafeTokenBucketRateLimiter(rate=10, burst=50)
    
    # Allow full burst
    allowed_burst = sum(1 for _ in range(100) if limiter.allow_request())
    assert allowed_burst <= 50, f"Burst allowed {allowed_burst} (max 50)"
    
    # After burst, rate should be limited to ~10/s
    time.sleep(1.0)
    count = sum(1 for _ in range(100) if limiter.allow_request())
    assert count <= 20, f"Post-burst allowed {count} in 1s (expected ~10)"
```

### 4. CI Integration

Add to CI pipeline:
- Run concurrent stress test with 100 threads × 100 requests each
- Assert allowed requests ≤ configured rate + 5% tolerance
- Run 10 iterations to detect non-deterministic failures
- Benchmark lock contention (allowRequest latency P99 < 1ms under load)
