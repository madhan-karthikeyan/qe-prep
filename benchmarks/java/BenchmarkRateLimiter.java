package benchmarks.java;

import com.qe.test.ratelimit.SlidingWindowLog;
import com.qe.test.ratelimit.TokenBucket;

public class BenchmarkRateLimiter {

    public static BenchmarkResult runTokenBucket(int iterations) {
        TokenBucket tb = new TokenBucket(iterations, iterations);

        long start = System.nanoTime();
        for (int i = 0; i < iterations; i++) {
            tb.allowRequest();
        }
        long end = System.nanoTime();

        long totalNanos = end - start;
        double opsPerSec = (double) iterations / (totalNanos / 1_000_000_000.0);
        double avgLatencyNanos = (double) totalNanos / iterations;

        return new BenchmarkResult("TokenBucket.allowRequest", iterations, opsPerSec, avgLatencyNanos, totalNanos);
    }

    public static BenchmarkResult runSlidingWindow(int iterations) {
        SlidingWindowLog sw = new SlidingWindowLog(1, SlidingWindowLog.TimeUnit.SECONDS, iterations);

        long start = System.nanoTime();
        for (int i = 0; i < iterations; i++) {
            sw.allowRequest();
        }
        long end = System.nanoTime();

        long totalNanos = end - start;
        double opsPerSec = (double) iterations / (totalNanos / 1_000_000_000.0);
        double avgLatencyNanos = (double) totalNanos / iterations;

        return new BenchmarkResult("SlidingWindowLog.allowRequest", iterations, opsPerSec, avgLatencyNanos, totalNanos);
    }
}
