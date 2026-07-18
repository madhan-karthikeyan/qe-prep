package com.qe.test.filter;

import static org.junit.jupiter.api.Assertions.*;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import com.qe.test.logger.LogLevel;
import com.qe.test.logger.LogRecord;

@DisplayName("LogFilter")
class LogFilterTest {

    @Test
    @DisplayName("byLevel accepts records at or above min level")
    void byLevel() {
        var filter = LogFilter.byLevel(LogLevel.WARN);
        assertFalse(filter.accept(new LogRecord(LogLevel.DEBUG, "debug")));
        assertFalse(filter.accept(new LogRecord(LogLevel.INFO, "info")));
        assertTrue(filter.accept(new LogRecord(LogLevel.WARN, "warn")));
        assertTrue(filter.accept(new LogRecord(LogLevel.ERROR, "error")));
    }

    @Test
    @DisplayName("byMessagePattern accepts records matching regex")
    void byMessagePattern() {
        var filter = LogFilter.byMessagePattern("(?i)error");
        assertTrue(filter.accept(new LogRecord(LogLevel.INFO, "something ERROR here")));
        assertFalse(filter.accept(new LogRecord(LogLevel.INFO, "all good")));
    }

    @Test
    @DisplayName("byPredicate delegates to predicate")
    void byPredicate() {
        var filter = LogFilter.byPredicate(r -> r.message().startsWith("A"));
        assertTrue(filter.accept(new LogRecord(LogLevel.INFO, "Alpha")));
        assertFalse(filter.accept(new LogRecord(LogLevel.INFO, "Beta")));
    }

    @Test
    @DisplayName("factory methods reject null")
    void factoryValidation() {
        assertThrows(IllegalArgumentException.class, () -> LogFilter.byLevel(null));
        assertThrows(IllegalArgumentException.class, () -> LogFilter.byMessagePattern(null));
        assertThrows(IllegalArgumentException.class, () -> LogFilter.byMessagePattern(""));
        assertThrows(IllegalArgumentException.class, () -> LogFilter.byPredicate(null));
    }
}
