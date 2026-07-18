package com.qe.test.filter;

import java.util.function.Predicate;
import java.util.regex.Pattern;

import com.qe.test.logger.LogLevel;
import com.qe.test.logger.LogRecord;

public class LogFilter {
    private final LogLevel minLevel;
    private final Pattern messagePattern;
    private final Predicate<LogRecord> predicate;

    private LogFilter(LogLevel minLevel, Pattern messagePattern, Predicate<LogRecord> predicate) {
        this.minLevel = minLevel;
        this.messagePattern = messagePattern;
        this.predicate = predicate;
    }

    public static LogFilter byLevel(LogLevel minLevel) {
        if (minLevel == null) throw new IllegalArgumentException("minLevel must not be null");
        return new LogFilter(minLevel, null, null);
    }

    public static LogFilter byMessagePattern(String regex) {
        if (regex == null || regex.isBlank()) throw new IllegalArgumentException("regex must not be null or blank");
        return new LogFilter(null, Pattern.compile(regex), null);
    }

    public static LogFilter byPredicate(Predicate<LogRecord> predicate) {
        if (predicate == null) throw new IllegalArgumentException("predicate must not be null");
        return new LogFilter(null, null, predicate);
    }

    public boolean accept(LogRecord record) {
        if (record == null) return false;
        if (minLevel != null && record.level().priority() < minLevel.priority()) return false;
        if (messagePattern != null && !messagePattern.matcher(record.message()).find()) return false;
        if (predicate != null && !predicate.test(record)) return false;
        return true;
    }
}
