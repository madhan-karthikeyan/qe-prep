package benchmarks.java;

import java.util.ArrayList;
import java.util.List;

public class BenchmarkRunner {

    public static void main(String[] args) {
        int iterations = args.length > 0 ? Integer.parseInt(args[0]) : 100_000;
        System.out.println("BenchmarkRunner: " + iterations + " iterations per test\n");

        List<BenchmarkResult> results = new ArrayList<>();

        // Logger benchmarks
        System.out.println("--- Logger Benchmarks ---");
        results.add(BenchmarkLogger.run(iterations));
        results.add(BenchmarkLogger.runWithFilter(iterations));

        // LRU benchmarks
        System.out.println("--- LRU Cache Benchmarks ---");
        results.add(BenchmarkLRU.runPuts(iterations));
        results.add(BenchmarkLRU.runGets(iterations));
        results.add(BenchmarkLRU.runMixed(iterations));

        // Rate limiter benchmarks
        System.out.println("--- Rate Limiter Benchmarks ---");
        results.add(BenchmarkRateLimiter.runTokenBucket(iterations));
        results.add(BenchmarkRateLimiter.runSlidingWindow(iterations));

        // Summary table
        System.out.println("\n--- Summary ---");
        printSummaryTable(results);
    }

    private static void printSummaryTable(List<BenchmarkResult> results) {
        String header = "| Benchmark | Iterations | Ops/sec | Avg Latency (ns) | Total Time (ms) |";
        String sep = "|---|---|---:|---:|---:|";
        System.out.println(header);
        System.out.println(sep);
        for (BenchmarkResult r : results) {
            System.out.printf(
                "| %s | %,d | %,.0f | %,.0f | %,.0f |%n",
                r.name, r.iterations, r.opsPerSec, r.avgLatencyNanos, r.totalNanos / 1_000_000.0
            );
        }
    }
}
