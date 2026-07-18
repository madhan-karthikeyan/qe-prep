## Symptoms
The program runs for a while but eventually crashes with `java.lang.OutOfMemoryError: Java heap space`, even with modest heap settings.

## Root Cause
A static `HashMap` caches every computed value without any eviction policy. Each entry holds a 1 MB `byte[]`, so memory fills up quickly and never gets freed.

## Fix
Use an LRU cache implemented with `LinkedHashMap` and an overridden `removeEldestEntry`, or use `WeakHashMap` for keys that can be garbage-collected, or a dedicated caching library like Caffeine.

## Prevention
- Always bound caches with a maximum size or TTL.
- Use `-Xmx` to limit heap and `-XX:+HeapDumpOnOutOfMemoryError` to capture heap dumps for analysis.
- Profile with JVisualVM or Eclipse MAT to find memory hogs.
