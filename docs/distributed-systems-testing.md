# Distributed Systems Testing

## Network Partition Testing

Simulate network splits between nodes to verify the system handles partitions gracefully.

```python
# Using toxiproxy to simulate network failures
from toxiproxy import Toxiproxy, ToxicDirection

proxy = Toxiproxy()
proxy.add_proxy("db_proxy", "localhost:15432", "db:5432")

# Cut connection for 5 seconds
with proxy["db_proxy"].disconnect() as toxic:
    time.sleep(5)
    # System should still serve reads from cache
    # System should queue writes for retry
```

**What to verify:**
- System does not lose data during partition
- System recovers when partition heals
- Read replicas still serve stale data (if acceptable)
- Writes are queued or rejected gracefully

## Chaos Engineering Principles

1. **Define steady state** — what "normal" looks like (latency p99 < 200ms, error rate < 0.1%)
2. **Hypothesize** — "If we kill the auth service, new users can't register but existing users still work."
3. **Introduce variables** — Kill processes, inject latency, fill disks
4. **Measure** — Compare against steady state
5. **Rollback** — Stop the experiment if the blast radius exceeds expectations

**Tools:** Chaos Monkey, Litmus, Gremlin, Toxiproxy, `tc` (traffic control).

```bash
# Inject 500ms latency on eth0 (Linux)
tc qdisc add dev eth0 root netem delay 500ms

# Remove
tc qdisc del dev eth0 root
```

## Jepsen-Style Testing

Jepsen tests distributed systems by constructing concurrent operations and checking that results are linearizable/consistent.

**Core ideas:**
1. Generate concurrent client operations
2. Introduce faults (partitions, crashes, clock skew)
3. Record all operations and their order
4. Check history against consistency model (linearizable, sequential, etc.)

**Simplified Jepsen-style check in Python:**

```python
def test_linearizable_register():
    history = []
    def client(worker_id):
        for i in range(10):
            key = f"k_{worker_id}"
            old = db.read(key)
            db.write(key, worker_id * 100 + i)
            history.append((worker_id, old, new))
    # Run concurrently with faults
    inject_partition(duration=2)
    assert is_linearizable(history)  # Check consistency
```

## Testing Consensus Implementations

Consensus (Raft, Paxos) is notoriously hard to test.

**Approach:**
1. Unit test each state transition (follower → candidate → leader)
2. Test election timeouts with clock mocking
3. Test leader failure during log replication
4. Test network partitions during commit
5. Test membership changes (add/remove nodes)
6. Test with Jepsen-style fault injection

**What can go wrong:**
- Split brain (two leaders) — rare but catastrophic
- Log divergence after leader crash
- Committed log entries that are then overwritten
- Slow nodes causing throughput collapse

## Failure Injection

| Failure | How | What to Check |
|---------|-----|---------------|
| **Process crash** | `kill -9` or `os.exit()` | Data durability on restart |
| **OOM** | Set memory limits, allocate | Graceful degradation |
| **Disk full** | Fill disk with `dd` | Logging, error messages |
| **Clock skew** | NTP disable + manual set | Lease expiry, TTL behavior |
| **Slow network** | `tc netem delay 500ms` | Timeouts, retry amplification |
| **Packet loss** | `tc netem loss 10%` | Retry logic, idempotency |
| **DNS failure** | Block port 53 | Service discovery fallback |

## Testing Retry Logic

Retries are common and commonly buggy.

```python
# Exponential backoff with jitter
def retry_with_backoff(fn, max_retries=3):
    for attempt in range(max_retries):
        try:
            return fn()
        except TransientError:
            if attempt == max_retries - 1:
                raise
            time.sleep(2 ** attempt + random.uniform(0, 1))
```

**Test:**
```python
@responses.activate
def test_retry_eventually_succeeds():
    call_count = 0
    def handler(_):
        nonlocal call_count
        call_count += 1
        if call_count < 3:
            return (500, {}, "Server Error")
        return (200, {}, "OK")

    responses.add_callback(responses.GET, "http://api.example.com/data", callback=handler)

    result = retry_with_backoff(lambda: requests.get("http://api.example.com/data"))
    assert call_count == 3
    assert result.status_code == 200
```

**What to verify:**
- Maximum retry count is respected
- Backoff increases exponentially
- Jitter prevents thundering herd
- Retry doesn't happen on non-transient errors (4xx)
- Idempotency — retrying the same request has no side effects

## Timeout and Latency Testing

```python
# Test timeout behavior
@responses.activate
def test_timeout_triggers_fallback():
    responses.add(
        responses.GET, "http://api.example.com/slow",
        body=Exception("timeout"),
    )
    result = service.fetch_with_fallback()
    assert result.source == "cache"  # falls back to cache
```

**Key patterns:**
- Always set timeouts on network calls (connect + read)
- Use circuit breakers for downstream failures
- Test timeout values: too short causes spurious failures, too long hangs the system
- Test that timeout errors are distinguishable from other failures

## State Verification in Distributed Systems

Checking state in a distributed system is harder than in a monolithic app.

```python
def test_eventual_consistency():
    # Write to primary
    primary.write("key1", "value1")
    # Wait for replication
    time.sleep(1)
    # Read from all replicas
    results = [replica.read("key1") for replica in replicas]
    # Eventually all replicas should have the value
    assert all(r == "value1" for r in results), f"Got {results}"
```

**Better approach:** Instead of sleeping, poll with a timeout:

```python
def wait_for_condition(check_fn, timeout=10):
    deadline = time.time() + timeout
    while time.time() < deadline:
        if check_fn():
            return True
        time.sleep(0.1)
    return False

def test_eventual_consistency():
    primary.write("key1", "value1")
    assert wait_for_condition(
        lambda: all(r.read("key1") == "value1" for r in replicas)
    )
```

**State verification checklist:**
- No data loss after node failure
- No duplicate records after retries
- Consistent across replicas (eventually)
- Correct ordering of events
- No zombie records after deletion
