# Profiling

## Why Profiling Matters for QE

Profiling answers:
- Why is this slow?
- Why is memory growing?
- Why is CPU at 100%?

Without profiling, performance debugging is guesswork. A QE engineer uses profiling to:
1. **Reproduce and characterize** performance bugs with data, not opinion
2. **Find the bottleneck** before escalation — saves developers time
3. **Write performance regression tests** that catch slowdowns before release
4. **Validate fixes** — "CPU usage dropped from 80% to 20% after the fix"

## Python: cProfile, memory_profiler, snakeviz

### CPU Profiling

```python
# cProfile — built-in, no dependencies
import cProfile
import pstats

profiler = cProfile.Profile()
profiler.enable()
run_my_code()
profiler.disable()

stats = pstats.Stats(profiler).sort_stats("cumtime")
stats.print_stats(20)  # top 20 functions by cumulative time
```

```bash
# Profile an entire script
python -m cProfile -o output.prof my_script.py
```

### Visualize with SnakeViz

```bash
pip install snakeviz
snakeviz output.prof  # Opens interactive flame chart in browser
```

### Memory Profiling

```python
from memory_profiler import profile

@profile
def process_large_data():
    data = load_all_from_db()        # check memory increment
    transformed = transform(data)     # intermediate copies?
    result = summarize(transformed)   # final output
    return result
```

```bash
python -m memory_profiler my_script.py
# Line-by-line memory usage printed to stdout
```

**Detecting memory leaks:**
```python
import tracemalloc
tracemalloc.start()
# run code
snapshot = tracemalloc.take_snapshot()
stats = snapshot.statistics("lineno")
for stat in stats[:10]:
    print(stat)
```

## Go: pprof

### CPU Profiling

```go
import (
    "os"
    "runtime/pprof"
)

func main() {
    f, _ := os.Create("cpu.prof")
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    runMyCode()
}
```

```bash
go tool pprof cpu.prof
(pprof) top     # top CPU consumers
(pprof) web     # open in browser
(pprof) list functionName  # line-by-line breakdown
```

### Memory (Heap) Profiling

```go
f, _ := os.Create("heap.prof")
pprof.WriteHeapProfile(f)
f.Close()
```

```bash
go tool pprof -inuse_space heap.prof    # find high memory usage
go tool pprof -alloc_objects heap.prof  # find high allocation count
go tool pprof -alloc_space heap.prof    # find high allocation size
```

### Goroutine and Mutex Profiling

```go
// Goroutine profile
pprof.Lookup("goroutine").WriteTo(f, 1)

// Mutex profile
pprof.Lookup("mutex").WriteTo(f, 1)
```

```bash
go tool pprof http://localhost:6060/debug/pprof/goroutine
go tool pprof http://localhost:6060/debug/pprof/mutex
```

### net/http/pprof (HTTP endpoint)

```go
import _ "net/http/pprof"
// Then start http.ListenAndServe(":6060", nil)
```

```bash
# Gather profiles live without stopping the app
go tool pprof -seconds 30 http://localhost:6060/debug/pprof/profile
go tool pprof http://localhost:6060/debug/pprof/heap
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

## Java: JFR, JMC, VisualVM

### Java Flight Recorder (JFR)

```bash
# Start recording (requires JDK 11+)
jcmd <pid> JFR.start duration=60s filename=recording.jfr

# Or from command line at startup
java -XX:StartFlightRecording=duration=60s,filename=recording.jfr -jar app.jar
```

### Java Mission Control (JMC)

GUI tool to open `.jfr` files and navigate events:
- **CPU** — hottest methods, thread CPU usage
- **Memory** — allocation pressure, GC pauses, heap usage
- **IO** — file and socket reads/writes
- **Threads** — lock contention, blocked threads
- **Exceptions** — where and how often errors occur

### VisualVM

All-in-one profiling tool for running JVMs:
- CPU sampler — see which methods consume CPU
- Memory sampler — track heap growth, detect leaks
- Threads — see thread states, detect deadlocks
- Heap dump analysis — find leak suspects, dominator tree

### On-Heap vs Off-Heap

```bash
# Heap dump for leak analysis
jmap -dump:live,format=b,file=heap.hprof <pid>
# Analyze with VisualVM or Eclipse MAT
```

## Common Performance Patterns

| Pattern | Symptom | Profiling Signal | Likely Cause |
|---------|---------|-----------------|--------------|
| **Memory leak** | Heap grows over time, OOM | Heap profile shows growth in one class | Objects retained longer than expected (e.g., unbounded cache, forgotten listener) |
| **CPU spike** | 100% CPU for extended period | CPU profile shows a single hot function | Inefficient algorithm, tight loop, infinite loop |
| **GC thrashing** | Frequent stop-the-world pauses, high GC % | Memory profile shows high allocation rate | Creating many short-lived objects; object pooling or allocation reduction needed |
| **Thread contention** | Low CPU but slow throughput | Thread dump shows many blocked threads | Lock contention, synchronized bottleneck |
| **IO bottleneck** | CPU low, response time high | Thread dump shows threads in network IO | Slow database query, chatty API calls |
| **Excessive allocations** | Frequent GC, high CPU | Allocation profile shows hot allocation sites | String concatenation in loops, boxing/unboxing |
| **Connection leak** | Connections exhausted, errors | Thread dump: many waiting for connection | Connection pool not released in `finally` |

## How to Read Flame Graphs

Flame graphs visualize stack traces over time.

```
_____________________________________________________________
| some_func                                    |            |
|   child_func                  |              |            |
|     grandchild  |  another    |  helper_func |            |
|  main           |  main     |  main       |  main        |
‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
```

**How to read:**
- **X-axis** — stack trace population (width = proportion of time spent there)
- **Y-axis** — stack depth (top = leaf function, bottom = entry point)
- **Color** — often randomized; some tools color by CPU vs IO vs lock
- **Wide blocks** = hot spots — that function + its children consume the most CPU

**Method:**
1. Find the widest block at the top (leaf) — that's the hottest function
2. Trace from there to understand the call path
3. If two hot paths share a parent, the problem may be upstream

**Generate flame graphs:**
```bash
# Go: convert pprof to flame graph
go tool pprof -http=:8080 cpu.prof

# Python: py-spy + flamegraph.pl
py-spy record -o profile.svg -- python my_script.py
```

## Performance Regression Testing

```python
import pytest

@pytest.mark.benchmark
def test_api_latency_regression(benchmark):
    client = TestClient(app)
    result = benchmark(client.get, "/api/users")
    assert result.status_code == 200
    # Benchmark framework asserts against thresholds
```

```go
func BenchmarkGetUser(b *testing.B) {
    for i := 0; i < b.N; i++ {
        app.GetUser("test-id")
    }
}
// go test -bench=. -benchtime=100x -count=5
```

**CI integration:**
1. Run performance tests on every PR merge to main
2. Compare against baseline (previous commit or predefined thresholds)
3. Fail the build if latency increased > 5% or memory > 10%
4. Store historical results for trend analysis

```yaml
# Example threshold config
thresholds:
  - name: "get_user_p99"
    max_ms: 200
  - name: "create_order_p95"
    max_ms: 500
  - name: "memory_after_test"
    max_mb: 256
```

## Profiling in CI

```yaml
# GitHub Actions + Go pprof
- name: Run performance tests
  run: |
    go test -bench=. -benchtime=100x -cpuprofile=cpu.prof -memprofile=mem.prof ./...
    go tool pprof -top -cum cpu.prof > cpu-top.txt
    go tool pprof -top -alloc_space mem.prof > mem-top.txt
  continue-on-error: true  # don't fail on perf regressions immediately

- name: Upload profiles
  uses: actions/upload-artifact@v4
  with:
    name: pprof-results
    path: |
      cpu.prof
      mem.prof
      cpu-top.txt
      mem-top.txt
```

**When to profile in CI:**
1. On every PR — lightweight quick benchmarks (not full CPU profiles)
2. Nightly — full profiling suite with flame graphs saved as artifacts
3. Pre-release — comprehensive performance validation

**Keep CI profiling fast:**
- Limit benchmark iterations (but enough for statistical significance)
- Use fixed resource limits (cpus, memory) for reproducible results
- Run in isolation — no other containers competing for CPU
- Compare against stored baselines, not absolute thresholds
