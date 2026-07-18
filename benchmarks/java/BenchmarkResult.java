package benchmarks.java;

public record BenchmarkResult(String name, int iterations, double opsPerSec, double avgLatencyNanos, long totalNanos) {
}
