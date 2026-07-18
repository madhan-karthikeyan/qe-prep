# Debugging Techniques

## Print Debugging vs Debugger vs Logging

| Technique | When to Use | Pros | Cons |
|-----------|------------|------|------|
| **Print debugging** | Quick exploration, one-off scripts | Zero setup, works everywhere | Noisy, must remove after, no persistence |
| **Debugger** | Complex logic, stepping through code | Inspect state at any point, conditional breakpoints | Requires setup, slows down execution |
| **Logging** | Production issues, long-running processes | Structured, searchable, persistent | Adds overhead, must think about levels beforehand |

**Rule of thumb:** Use `print()` for prototyping, debugger for complex logic, logging for production.

### Debugger basics

```python
import pdb; pdb.set_trace()  # Python — set breakpoint
# Or in Python 3.7+: breakpoint()
```

```go
dlv debug main.go  # Go — Delve debugger
```

```java
// Java — IDE debugger (IntelliJ, Eclipse):
// Click gutter to set breakpoint, Run > Debug
```

## Binary Search (git bisect)

Find the commit that introduced a bug in logarithmic time.

```bash
git bisect start
git bisect bad           # current commit is broken
git bisect good v1.0     # known good release
# Git checks out middle commit
# Run your test:
python -m pytest tests/test_regression.py
git bisect good          # if test passes
git bisect bad           # if test fails
# Repeat until the first bad commit is identified
git bisect reset
```

**Automated bisect:**
```bash
git bisect run python -m pytest tests/test_regression.py
```

## Rubber Duck Debugging

Explain the problem out loud (or to an inanimate object). The act of verbalizing often reveals the missing assumption or flawed logic.

**Process:**
1. State what the code *should* do
2. State what it *actually* does
3. Walk through the code line by line
4. At some point you'll notice: "Wait, that variable isn't what I thought..."

Works best when you're stuck and can't find a mentor.

## Reading Stack Traces

```
Traceback (most recent call last):
  File "app.py", line 25, in process_order
    result = payment_service.charge(order.total)
             ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "payment.py", line 42, in charge
    response = http_client.post(url, data=payload)
               ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "http.py", line 88, in post
    return self.session.post(url, json=data, timeout=10)
           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
requests.exceptions.ConnectTimeout: Connection to payments.example.com timed out
```

**How to read it:**
1. **Last line** — the actual error. Always start here.
2. **Bottom of stack** — where the error occurred.
3. **Top of stack** — the entry point that triggered the call chain.
4. **Each frame** — file, line number, and the line of code.

**Key questions:**
- What was the input that went in at the top?
- Where did the first unexpected value appear?

## Profiling Basics

### CPU Profiling

Find what's consuming CPU cycles.

```bash
python -m cProfile -o output.prof my_script.py
python -m snakeviz output.prof  # Visualize in browser
```

```go
import _ "net/http/pprof"
// Then: go tool pprof http://localhost:6060/debug/pprof/profile
```

```bash
# Java: Start with -XX:+FlightRecorder
jcmd <pid> JFR.start duration=60s filename=recording.jfr
```

### Memory Profiling

Find memory leaks and excessive allocations.

```python
from memory_profiler import profile

@profile
def process_large_file():
    data = load_file()  # <-- see memory usage per line
```

```go
go tool pprof -alloc_space http://localhost:6060/debug/pprof/heap
```

```bash
# Java
jcmd <pid> GC.heap_info
jmap -heap <pid>
```

### I/O Profiling

```bash
# System level
iostat -x 1      # disk I/O
lsof -p <pid>    # open files by process
strace -p <pid>  # system calls
```

## Log Analysis Patterns

| Pattern | Signal |
|---------|--------|
| **Error burst** | Connection pool exhausted, rate limited |
| **Latency spike** | GC pause, network blip, resource contention |
| **Missing logs** | Code path not reached — early return or exception swallowed |
| **Repeated retry logs** | Downstream service unstable |
| **Timestamp gap** | Blocking operation (full GC, file lock, deadlock) |

**Approach:**
1. Find the first error, not the last (cascading failures amplify)
2. Look at the 5 seconds before the error
3. Correlate across services (request ID, trace ID)
4. Count occurrences — is it 1 in 1000 or 1 in 2?

## Reproducing Intermittent Failures

| Technique | Method |
|-----------|--------|
| **Run tests 100x** | `pytest --count=100` or `go test -count=100` |
| **Randomize order** | `pytest --random-order` |
| **Stress the system** | Run under load, race detection enabled |
| **Capture state on failure** | Save logs, DB snapshot, thread dump on assertion failure |
| **Add more logging** | Temporarily add debug logs to narrow down timing |
| **Simplify the test** | Strip it down to the minimum that still fails |
| **Try different environments** | CI vs local, different OS, different resource limits |

**Mental model:** An intermittent failure is *always* deterministic — you just haven't found the control variable yet.
