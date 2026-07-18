package com.qe.test.logger;

import java.io.IOException;
import java.io.PrintStream;
import java.time.ZoneId;
import java.time.ZonedDateTime;
import java.time.format.DateTimeFormatter;
import java.util.concurrent.locks.ReentrantLock;

import com.qe.test.filter.LogFilter;
import com.qe.test.rotation.RotatingFileWriter;

public class Logger implements AutoCloseable {
    private static final DateTimeFormatter FORMATTER = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss.SSS Z");

    private final String name;
    private LogLevel threshold;
    private PrintStream stdoutTarget;
    private RotatingFileWriter fileWriter;
    private LogFilter filter;
    private final ReentrantLock lock = new ReentrantLock();
    private boolean closed;

    public Logger(String name) {
        if (name == null || name.isBlank()) {
            throw new IllegalArgumentException("name must not be null or blank");
        }
        this.name = name;
        this.threshold = LogLevel.DEBUG;
        this.stdoutTarget = System.out;
    }

    public String name() { return name; }
    public LogLevel threshold() { return threshold; }

    public void setThreshold(LogLevel threshold) {
        if (threshold == null) throw new IllegalArgumentException("threshold must not be null");
        this.threshold = threshold;
    }

    public void setStdoutTarget(PrintStream out) {
        if (out == null) throw new IllegalArgumentException("out must not be null");
        this.stdoutTarget = out;
    }

    public void setFileWriter(RotatingFileWriter fileWriter) {
        this.fileWriter = fileWriter;
    }

    public void setFilter(LogFilter filter) {
        this.filter = filter;
    }

    private boolean shouldLog(LogLevel level) {
        return level.priority() >= threshold.priority();
    }

    public void log(LogLevel level, String message) {
        log(new LogRecord(level, message));
    }

    public void log(LogLevel level, String message, Throwable thrown) {
        log(new LogRecord(level, message, thrown));
    }

    private void log(LogRecord record) {
        if (closed) return;
        if (!shouldLog(record.level())) return;

        if (filter != null && !filter.accept(record)) return;

        String formatted = format(record);
        lock.lock();
        try {
            if (stdoutTarget != null) {
                stdoutTarget.println(formatted);
            }
            if (fileWriter != null) {
                try {
                    fileWriter.write(formatted + System.lineSeparator());
                } catch (IOException e) {
                    System.err.println("Failed to write to log file: " + e.getMessage());
                }
            }
        } finally {
            lock.unlock();
        }
    }

    private String format(LogRecord record) {
        var zdt = ZonedDateTime.ofInstant(record.timestamp(), ZoneId.systemDefault());
        var sb = new StringBuilder();
        sb.append('[').append(zdt.format(FORMATTER)).append(']')
          .append(' ').append(record.level().name())
          .append(' ').append(name)
          .append(" - ").append(record.message());
        if (record.thrown() != null) {
            sb.append(System.lineSeparator());
            var sw = new java.io.StringWriter();
            record.thrown().printStackTrace(new java.io.PrintWriter(sw));
            sb.append(sw);
        }
        return sb.toString();
    }

    @Override
    public void close() throws IOException {
        lock.lock();
        try {
            closed = true;
            if (fileWriter != null) {
                fileWriter.close();
            }
        } finally {
            lock.unlock();
        }
    }
}
