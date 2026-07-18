package com.qe.test.api;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;

import java.time.Duration;

import static org.junit.jupiter.api.Assertions.*;

@DisplayName("ApiTester")
class ApiTesterTest {

    @Test
    @DisplayName("rejects null baseUrl")
    void rejectsNullBaseUrl() {
        assertThrows(IllegalArgumentException.class,
                () -> new ApiTester(null, Duration.ofSeconds(5), 0, Duration.ofMillis(100)));
    }

    @Test
    @DisplayName("rejects blank baseUrl")
    void rejectsBlankBaseUrl() {
        assertThrows(IllegalArgumentException.class,
                () -> new ApiTester("  ", Duration.ofSeconds(5), 0, Duration.ofMillis(100)));
    }

    @Test
    @DisplayName("rejects null timeout")
    void rejectsNullTimeout() {
        assertThrows(IllegalArgumentException.class,
                () -> new ApiTester("http://localhost", null, 0, Duration.ofMillis(100)));
    }

    @Test
    @DisplayName("rejects negative maxRetries")
    void rejectsNegativeMaxRetries() {
        assertThrows(IllegalArgumentException.class,
                () -> new ApiTester("http://localhost", Duration.ofSeconds(5), -1, Duration.ofMillis(100)));
    }

    @Test
    @DisplayName("fails gracefully when server is not running")
    @Timeout(10)
    void failsWhenServerNotRunning() {
        var tester = new ApiTester("http://localhost:19999", Duration.ofSeconds(1), 0, Duration.ofMillis(50));
        var result = tester.sendRequest("GET", "/api/items", null);
        assertFalse(result.passed());
        assertEquals(-1, result.statusCode());
    }

    @Test
    @DisplayName("assertStatus passes on matching status")
    void assertStatusPasses() {
        var result = new ApiTester.TestResult(200, "ok", Duration.ZERO, true, null);
        assertDoesNotThrow(() -> new ApiTester("http://localhost", Duration.ofSeconds(5), 0, Duration.ofMillis(100))
                .assertStatus(result, 200));
    }

    @Test
    @DisplayName("assertStatus throws on non-matching status")
    void assertStatusThrows() {
        var result = new ApiTester.TestResult(404, "not found", Duration.ZERO, false, null);
        assertThrows(AssertionError.class,
                () -> new ApiTester("http://localhost", Duration.ofSeconds(5), 0, Duration.ofMillis(100))
                        .assertStatus(result, 200));
    }

    @Test
    @DisplayName("assertJson finds key in body")
    void assertJsonFindsKey() {
        var result = new ApiTester.TestResult(200, "{\"name\":\"test\"}", Duration.ZERO, true, null);
        assertDoesNotThrow(() -> new ApiTester("http://localhost", Duration.ofSeconds(5), 0, Duration.ofMillis(100))
                .assertJson(result, "name", "test"));
    }

    @Test
    @DisplayName("assertJson throws when key not found")
    void assertJsonThrowsOnMissingKey() {
        var result = new ApiTester.TestResult(200, "{\"other\":\"value\"}", Duration.ZERO, true, null);
        assertThrows(AssertionError.class,
                () -> new ApiTester("http://localhost", Duration.ofSeconds(5), 0, Duration.ofMillis(100))
                        .assertJson(result, "name", "test"));
    }

    @Test
    @DisplayName("retry on failure")
    @Timeout(10)
    void retryOnFailure() {
        var tester = new ApiTester("http://localhost:19998", Duration.ofSeconds(1), 2, Duration.ofMillis(10));
        var result = tester.sendRequestWithRetry("GET", "/api/items", null);
        assertFalse(result.passed());
        assertTrue(result.errorMessage().contains("3 attempts"));
    }

    @Test
    @DisplayName("TestResult rejects null duration")
    void testResultRejectsNullDuration() {
        assertThrows(NullPointerException.class,
                () -> new ApiTester.TestResult(200, "ok", null, true, null));
    }

    @Test
    @DisplayName("TestResult records duration")
    void testResultRecordsDuration() {
        var result = new ApiTester.TestResult(200, "ok", Duration.ofMillis(50), true, null);
        assertEquals(Duration.ofMillis(50), result.duration());
        assertTrue(result.passed());
    }
}
