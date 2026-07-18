package benchmarks.java;

import com.qe.test.logger.LogLevel;
import com.qe.test.logger.Logger;

public class BenchmarkLogger {

    public static BenchmarkResult run(int iterations) {
        Logger logger = new Logger("benchmark");

        long start = System.nanoTime();
        for (int i = 0; i < iterations; i++) {
            logger.log(LogLevel.INFO, "benchmark message " + i);
        }
        long end = System.nanoTime();

        long totalNanos = end - start;
        double opsPerSec = (double) iterations / (totalNanos / 1_000_000_000.0);
        double avgLatencyNanos = (double) totalNanos / iterations;

        return new BenchmarkResult("Logger.info", iterations, opsPerSec, avgLatencyNanos, totalNanos);
    }

    public static BenchmarkResult runWithFilter(int iterations) {
        Logger logger = new Logger("benchmark");
        logger.setThreshold(LogLevel.ERROR);

        long start = System.nanoTime();
        for (int i = 0; i < iterations; i++) {
            logger.log(LogLevel.INFO, "filtered message " + i);
        }
        long end = System.nanoTime();

        long totalNanos = end - start;
        double opsPerSec = (double) iterations / (totalNanos / 1_000_000_000.0);
        double avgLatencyNanos = (double) totalNanos / iterations;

        return new BenchmarkResult("Logger.info (filtered)", iterations, opsPerSec, avgLatencyNanos, totalNanos);
    }
}
