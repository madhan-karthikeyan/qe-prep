package com.qe.test.logger;

import static org.junit.jupiter.api.Assertions.*;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.atomic.AtomicInteger;

import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import com.qe.test.filter.LogFilter;
import com.qe.test.rotation.RotatingFileWriter;

@DisplayName("Logger")
class LoggerTest {

    private ByteArrayOutputStream baos;
    private PrintStream out;
    private Logger logger;

    @BeforeEach
    void setUp() {
        baos = new ByteArrayOutputStream();
        out = new PrintStream(baos, true, StandardCharsets.UTF_8);
        logger = new Logger("test");
        logger.setStdoutTarget(out);
    }

    @AfterEach
    void tearDown() throws IOException {
        logger.close();
        out.close();
    }

    @Test
    @DisplayName("throws on null or blank name")
    void constructorValidation() {
        assertThrows(IllegalArgumentException.class, () -> new Logger(null));
        assertThrows(IllegalArgumentException.class, () -> new Logger("   "));
    }

    @Test
    @DisplayName("log prints formatted message to stdout")
    void logToStdout() {
        logger.log(LogLevel.INFO, "hello world");
        String output = baos.toString(StandardCharsets.UTF_8);
        assertTrue(output.contains("INFO"));
        assertTrue(output.contains("test"));
        assertTrue(output.contains("hello world"));
    }

    @Test
    @DisplayName("log filtered by threshold does not print")
    void thresholdFiltering() {
        logger.setThreshold(LogLevel.ERROR);
        logger.log(LogLevel.INFO, "should not appear");
        String output = baos.toString(StandardCharsets.UTF_8);
        assertFalse(output.contains("should not appear"));
    }

    @Test
    @DisplayName("log includes stack trace when throwable provided")
    void logWithThrowable() {
        var cause = new RuntimeException("boom");
        logger.log(LogLevel.ERROR, "error occurred", cause);
        String output = baos.toString(StandardCharsets.UTF_8);
        assertTrue(output.contains("RuntimeException"));
        assertTrue(output.contains("boom"));
    }

    @Test
    @DisplayName("filter by level predicate")
    void filtersByLevelPredicate() {
        logger.setFilter(LogFilter.byLevel(LogLevel.WARN));
        logger.log(LogLevel.INFO, "info");
        logger.log(LogLevel.WARN, "warn");
        logger.log(LogLevel.ERROR, "error");
        String output = baos.toString(StandardCharsets.UTF_8);
        assertFalse(output.contains("info"));
        assertTrue(output.contains("warn"));
        assertTrue(output.contains("error"));
    }

    @Test
    @DisplayName("filter by message pattern")
    void filtersByMessagePattern() {
        logger.setFilter(LogFilter.byMessagePattern("ERROR"));
        logger.log(LogLevel.INFO, "info no match");
        logger.log(LogLevel.ERROR, "this is an ERROR message");
        String output = baos.toString(StandardCharsets.UTF_8);
        assertFalse(output.contains("info no match"));
        assertTrue(output.contains("ERROR message"));
    }

    @Test
    @DisplayName("writes to file when file writer configured")
    void writesToFile() throws Exception {
        Path tmp = Files.createTempFile("logger-test", ".log");
        try (var fw = new RotatingFileWriter(tmp, 1024 * 1024)) {
            logger.setFileWriter(fw);
            logger.log(LogLevel.INFO, "file content");
            fw.flush();
            String content = Files.readString(tmp);
            assertTrue(content.contains("file content"));
        } finally {
            Files.deleteIfExists(tmp);
        }
    }

    @Test
    @DisplayName("multi-threaded logging does not corrupt output")
    void multiThreadedLogging() throws Exception {
        int threadCount = 10;
        int linesPerThread = 50;
        var latch = new CountDownLatch(threadCount);
        var executor = Executors.newFixedThreadPool(threadCount);

        for (int i = 0; i < threadCount; i++) {
            int id = i;
            executor.submit(() -> {
                try {
                    for (int j = 0; j < linesPerThread; j++) {
                        logger.log(LogLevel.INFO, "thread-" + id + " line-" + j);
                    }
                } finally {
                    latch.countDown();
                }
            });
        }
        latch.await();
        executor.shutdown();
        out.flush();
        String output = baos.toString(StandardCharsets.UTF_8);
        long count = output.lines().count();
        assertEquals(threadCount * linesPerThread, count);
    }
}
