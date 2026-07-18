package com.qe.test.net;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.time.Duration;

import static org.junit.jupiter.api.Assertions.*;

@DisplayName("HttpClient")
class HttpClientTest {

    @Test
    @DisplayName("rejects null base delay")
    void rejectsNullBaseDelay() {
        assertThrows(IllegalArgumentException.class,
                () -> new HttpClient(null, Duration.ofSeconds(1), 0.2, 3));
    }

    @Test
    @DisplayName("rejects null max delay")
    void rejectsNullMaxDelay() {
        assertThrows(IllegalArgumentException.class,
                () -> new HttpClient(Duration.ofMillis(100), null, 0.2, 3));
    }

    @Test
    @DisplayName("rejects negative base delay")
    void rejectsNegativeBaseDelay() {
        assertThrows(IllegalArgumentException.class,
                () -> new HttpClient(Duration.ofMillis(-1), Duration.ofSeconds(1), 0.2, 3));
    }

    @Test
    @DisplayName("rejects jitter factor out of range")
    void rejectsInvalidJitterFactor() {
        assertThrows(IllegalArgumentException.class,
                () -> new HttpClient(Duration.ofMillis(100), Duration.ofSeconds(1), -0.1, 3));
        assertThrows(IllegalArgumentException.class,
                () -> new HttpClient(Duration.ofMillis(100), Duration.ofSeconds(1), 1.5, 3));
    }

    @Test
    @DisplayName("rejects negative max retries")
    void rejectsNegativeMaxRetries() {
        assertThrows(IllegalArgumentException.class,
                () -> new HttpClient(Duration.ofMillis(100), Duration.ofSeconds(1), 0.2, -1));
    }
}
