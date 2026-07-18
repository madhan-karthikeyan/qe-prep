# Concurrency Testing

## Race Condition Detection

A race condition occurs when the behavior of a program depends on the interleaving of operations across threads/goroutines.

### Python

```python
import threading

counter = 0
def increment():
    global counter
    for _ in range(100000):
        counter += 1  # RACE: read + write not atomic

threads = [threading.Thread(target=increment) for _ in range(10)]
for t in threads: t.start()
for t in threads: t.join()
print(counter)  # Expected: 1000000, Actual: ~950000
```

**Detection:** Use `threading.Lock` or `queue.Queue`. For tests, run under `pytest -x --count=10`.

### Go

```go
var counter int
for i := 0; i < 10; i++ {
    go func() {
        for j := 0; j < 100000; j++ {
            counter++ // RACE
        }
    }()
}
```

**Detection:** `go test -race`

### Java

```java
class Counter {
    private int count = 0;  // RACE: use AtomicInteger
    public void increment() { count++; }
}
```

**Detection:** Use `Thread` stress testing or specialized tools (vmlens, JCStress).

## Deadlock Detection

All threads are blocked waiting for resources held by each other.

### Python

```python
lock1 = threading.Lock()
lock2 = threading.Lock()

def thread_a():
    with lock1:
        with lock2: pass

def thread_b():
    with lock2:
        with lock1: pass  # ORDER swapped — deadlock possible
```

**Detection:** Thread dump / `threading.enumerate()`. Use locks in consistent order.

### Go

```go
var mu1, mu2 sync.Mutex

go func() { mu1.Lock(); mu2.Lock(); mu2.Unlock(); mu1.Unlock() }()
go func() { mu2.Lock(); mu1.Lock(); mu1.Unlock(); mu2.Unlock() }()
```

**Detection:** `go test -race` catches data races but not all deadlocks. Use `net/http/pprof` endpoint and examine goroutine stack traces: `go tool pprof http://localhost:6060/debug/pprof/goroutine`.

### Java

```java
synchronized(a) { synchronized(b) { ... } }  // Thread 1
synchronized(b) { synchronized(a) { ... } }  // Thread 2 — deadlock
```

**Detection:** `jstack <pid>` or `ThreadMXBean.findMonitorDeadlockedThreads()`.

## Stress Testing Patterns

Test under high concurrency to surface latent races.

```python
from concurrent.futures import ThreadPoolExecutor

def test_concurrent_cart_updates():
    cart = ShoppingCart()
    def add_item():
        cart.add("item", 1)
    with ThreadPoolExecutor(max_workers=20) as ex:
        futures = [ex.submit(add_item) for _ in range(100)]
        wait(futures)
    assert cart.count() == 100  # Fails if race present
```

```go
func TestConcurrentWrites(t *testing.T) {
    db := NewDB()
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            db.Write(Key(id), Value(id))
        }(i)
    }
    wg.Wait()
    // verify state
}
```

## Deterministic Testing with Controlled Scheduling

Instead of relying on random thread interleavings, control the scheduler to test specific orderings.

```python
import threading
import time

class ControlledScheduler:
    def __init__(self):
        self.events = {}
        self.barrier = threading.Barrier(2)

    def schedule(self, thread_id, checkpoint):
        """Ensure thread reaches checkpoint in order"""
        self.events[checkpoint] = thread_id
        self.barrier.wait()  # forces synchronization point
```

**Tools:**
- Python: Use barriers, events, and queues for deterministic interleaving
- Go: Use `sync.WaitGroup` and channels to enforce order
- Java: JCStress can test specific memory orderings

## Go Race Detector

```bash
# Build with race detection
go build -race
go test -race ./...

# Always run race-enabled tests in CI
go test -race -count=1 ./...
```

The race detector instruments memory accesses and reports conflicting access patterns. It adds CPU/memory overhead (~5-10x) but catches most data races.

**Limitations:** Only finds races that happen during the test run. Does not prove absence of races.

## Java ThreadMXBean

```java
ThreadMXBean threadMxBean = ManagementFactory.getThreadMXBean();

// Check for deadlocked threads
long[] deadlockedThreadIds = threadMxBean.findMonitorDeadlockedThreads();
if (deadlockedThreadIds != null) {
    for (long id : deadlockedThreadIds) {
        ThreadInfo info = threadMxBean.getThreadInfo(id);
        System.err.println("Deadlocked: " + info.getThreadName());
    }
}

// Thread CPU time
long cpuTime = threadMxBean.getThreadCpuTime(threadId);
```

## Python Threading Issues

**GIL:** Python's Global Interpreter Lock means true parallelism for CPU-bound work requires `multiprocessing`.

**Common bugs:**
- `threading.Condition` — forgotten `notify()` / `notify_all()` — thread sleeps forever
- `queue.Queue` — `get()` blocks forever if item never arrives; use `get(timeout=5)`
- Shared mutable state — plain `list`/`dict` modified from multiple threads (use `Lock` or `queue`)
- `threading.Event` — `wait(timeout)` returns `True` even if event was never set (check return value)
- `concurrent.futures` — uncaught exceptions in threads are swallowed; check `future.exception()`

## Common Concurrency Bugs

| Bug | Pattern | Fix |
|-----|---------|-----|
| **Check-then-act** | `if not in dict: dict[key] = value` | Use atomic `dict.setdefault()` or lock |
| **Lost update** | `counter += 1` | Atomic counter or mutex |
| **TOCTOU** | File exists check then open | Handle race in open (try/except) |
| **ABA problem** | Pointer reuse without versioning | Hazard pointers, epoch-based reclamation |
| **Double-checked locking** | `if cache is None: lock; if cache is None:` | Use atomic or `sync.Once` |
| **Thread starvation** | Low-priority threads never run | Bounded queues, fair locks |
| **Missing volatile** | Reader never sees writer's update | `volatile`/atomic in Java, `sync.Mutex` in Go |
