package com.qe.test.logger;

import java.time.Instant;

public record LogRecord(Instant timestamp, LogLevel level, String message, Throwable thrown) {
    public LogRecord {
        if (level == null) {
            throw new IllegalArgumentException("level must not be null");
        }
        if (message == null) {
            throw new IllegalArgumentException("message must not be null");
        }
        if (timestamp == null) {
            timestamp = Instant.now();
        }
    }

    public LogRecord(LogLevel level, String message) {
        this(Instant.now(), level, message, null);
    }

    public LogRecord(LogLevel level, String message, Throwable thrown) {
        this(Instant.now(), level, message, thrown);
    }
}
