import com.qe.test.logger.LogLevel;
import com.qe.test.logger.Logger;
import com.qe.test.cache.LRUCache;
import com.qe.test.ratelimit.TokenBucket;
import com.qe.test.ratelimit.SlidingWindowLog;

/**
 * Java profiling example using manual timers.
 *
 * Run with JFR (Java Flight Recorder) for detailed profiling:
 *   javac -cp <classpath> ProfilingExample.java
 *   java -XX:StartFlightRecording=filename=record.jfr,duration=30s \
 *        -cp <classpath> ProfilingExample
 *
 * Then open record.jfr with JDK Mission Control (jmc).
 */
public class ProfilingExample {

    static final int ITERATIONS = 100_000;

    public static void main(String[] args) throws Exception {
        System.out.println("=== Profiling Example ===");
        System.out.println("Iterations: " + ITERATIONS + "\n");

        exerciseLogger();
        exerciseLRU();
        exerciseRateLimiter();

        System.out.println("\nDone. Profiles written to record.jfr (if JFR was enabled).");
        System.out.println("Open with: jmc record.jfr");
    }

    static void exerciseLogger() {
        Logger logger = new Logger("profiling-example");
        for (int i = 0; i < ITERATIONS; i++) {
            logger.log(LogLevel.INFO, "profiling message " + i);
        }
        System.out.println("Logger: " + ITERATIONS + " messages logged");
    }

    static void exerciseLRU() {
        LRUCache<Integer, Integer> cache = new LRUCache<>(1000);
        for (int i = 0; i < ITERATIONS; i++) {
            cache.put(i, i);
            if (i % 2 == 0) {
                cache.get(i);
            } else {
                cache.get(i - 1);
            }
        }
        System.out.println("LRU: " + ITERATIONS + " put+get operations");
    }

    static void exerciseRateLimiter() {
        TokenBucket tb = new TokenBucket(ITERATIONS, ITERATIONS);
        for (int i = 0; i < ITERATIONS; i++) {
            tb.allowRequest();
        }

        SlidingWindowLog sw = new SlidingWindowLog(1, SlidingWindowLog.TimeUnit.SECONDS, ITERATIONS);
        for (int i = 0; i < ITERATIONS; i++) {
            sw.allowRequest();
        }

        System.out.println("RateLimiter: " + ITERATIONS + " TokenBucket + " + ITERATIONS + " SlidingWindow ops");
    }
}
