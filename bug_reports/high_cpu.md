# High CPU usage during log rotation on large files

**Severity:** Critical
**Priority:** P1
**Environment:** All platforms, Python 3.8+, log files >500MB
**Component:** Logging subsystem — RotatingFileHandler

## Summary

When `RotatingFileHandler` triggers a rotation on a log file >= 2GB, the process consumes 100% CPU for 30-120 seconds, blocking the main thread and causing request timeouts for all services sharing the process.

## Steps to Reproduce

1. Configure `RotatingFileHandler` with `maxBytes=500MB`, `backupCount=5`
2. Generate log data until the active file exceeds 2GB
3. Observe rotation trigger
4. Monitor CPU usage during rotation

## Expected Behavior

Rotation should complete within <500ms with minimal CPU impact. The application should remain responsive during rotation.

## Actual Behavior

- CPU spikes to 100% for 30-120 seconds
- All application threads are blocked (rotation runs in the logging handler's critical section)
- Request latency spikes to timeout levels
- In extreme cases, health check probes fail and the container is killed

## Logs / Screenshots

```
2025-03-15 10:15:01 [INFO] Log rotation triggered for app.log
2025-03-15 10:16:45 [INFO] Log rotation complete: app.log -> app.log.1
```
(90 second gap with no logs — rotation blocked all I/O)

CPU profile during rotation:
```
  55.2s  read()  (reading entire 2GB file into memory)
  12.1s  gzip.compress()
  22.7s  write()  (writing compressed archive)
```

## Root Cause Analysis

The `RotatingFileHandler.doRollover()` method reads the **entire** log file into memory before performing the archive operation:

```python
# Simplified problematic code
def doRollover(self):
    self.stream.close()
    if os.path.exists(self.baseFilename):
        with open(self.baseFilename, 'rb') as f:
            content = f.read()            # <-- loads 2GB into memory
        with open(self.archive_name, 'wb') as f:
            f.write(compress(content))    # <-- also in memory
    os.remove(self.baseFilename)
```

This design has two issues:
1. **Memory**: Requires 2× file size RAM (read buffer + compressed buffer)
2. **CPU**: Blocking read of a cold-cache 2GB file causes page cache thrashing

## Fix

Replace the in-memory copy with a streaming copy using a fixed-size buffer:

```python
def doRollover(self):
    self.stream.close()
    if os.path.exists(self.baseFilename):
        BUFFER_SIZE = 64 * 1024  # 64KB buffer
        with open(self.baseFilename, 'rb') as src, \
             open(self.archive_name, 'wb') as dst:
            while True:
                chunk = src.read(BUFFER_SIZE)
                if not chunk:
                    break
                dst.write(chunk)    # or apply compression per-chunk
        os.remove(self.baseFilename)
    self.stream = open(self.baseFilename, 'a')
```

For compression scenarios, use a compressor wrapper that accepts chunked input:

```python
import gzip
CHUNK = 64 * 1024
with open(self.baseFilename, 'rb') as src, \
     gzip.open(self.archive_name, 'wb') as dst:
    while chunk := src.read(CHUNK):
        dst.write(chunk)
```

## Regression Tests

### 1. Rotation Performance Benchmark

```python
def test_rotation_performance():
    handler = RotatingFileHandler('test.log', maxBytes=1024, backupCount=3)
    # Fill file to 10MB (simulating large rotation)
    with open('test.log', 'wb') as f:
        f.write(b'x' * 10 * 1024 * 1024)
    
    start = time.perf_counter()
    handler.doRollover()
    elapsed = time.perf_counter() - start
    
    # Assert: rotation completes in <500ms for 10MB file
    assert elapsed < 0.5, f"Rotation took {elapsed:.2f}s"
    
    # Assert: peak memory stays reasonable
    import tracemalloc
    tracemalloc.start()
    handler.doRollover()
    _, peak = tracemalloc.get_traced_memory()
    assert peak < 100 * 1024 * 1024  # <100MB peak
```

### 2. Non-Blocking Assertion

```python
def test_rotation_does_not_block_other_threads():
    import threading
    results = []
    
    def background_work():
        for _ in range(50):
            time.sleep(0.01)
        results.append("done")
    
    t = threading.Thread(target=background_work)
    t.start()
    
    handler.doRollover()  # large file rotation
    t.join()
    
    assert results[0] == "done", "Rotation blocked background thread"
```

### 3. CI Integration

Add a benchmark step in CI that creates a 100MB log file, rotates it, and asserts:
- Time < 2s
- Memory delta < 200MB
- No other goroutines/threads are starved
