package benchmarks.java;

import com.qe.test.cache.LRUCache;

public class BenchmarkLRU {

    public static BenchmarkResult runPuts(int iterations) {
        LRUCache<Integer, Integer> cache = new LRUCache<>(iterations);

        long start = System.nanoTime();
        for (int i = 0; i < iterations; i++) {
            cache.put(i, i);
        }
        long end = System.nanoTime();

        long totalNanos = end - start;
        double opsPerSec = (double) iterations / (totalNanos / 1_000_000_000.0);
        double avgLatencyNanos = (double) totalNanos / iterations;

        return new BenchmarkResult("LRU.put", iterations, opsPerSec, avgLatencyNanos, totalNanos);
    }

    public static BenchmarkResult runGets(int iterations) {
        LRUCache<Integer, Integer> cache = new LRUCache<>(iterations);
        for (int i = 0; i < iterations; i++) {
            cache.put(i, i);
        }

        long start = System.nanoTime();
        for (int i = 0; i < iterations; i++) {
            cache.get(i);
        }
        long end = System.nanoTime();

        long totalNanos = end - start;
        double opsPerSec = (double) iterations / (totalNanos / 1_000_000_000.0);
        double avgLatencyNanos = (double) totalNanos / iterations;

        return new BenchmarkResult("LRU.get (hit)", iterations, opsPerSec, avgLatencyNanos, totalNanos);
    }

    public static BenchmarkResult runMixed(int iterations) {
        LRUCache<Integer, Integer> cache = new LRUCache<>(iterations / 2);

        long start = System.nanoTime();
        for (int i = 0; i < iterations; i++) {
            cache.put(i, i);
            cache.get(i / 2);
        }
        long end = System.nanoTime();

        long totalNanos = end - start;
        double opsPerSec = (double) iterations * 2 / (totalNanos / 1_000_000_000.0);
        double avgLatencyNanos = (double) totalNanos / (iterations * 2);

        return new BenchmarkResult("LRU mixed put+get", iterations * 2, opsPerSec, avgLatencyNanos, totalNanos);
    }
}
