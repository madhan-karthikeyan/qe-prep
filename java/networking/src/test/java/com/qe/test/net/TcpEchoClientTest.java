package com.qe.test.net;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;

import java.io.IOException;
import java.time.Duration;

import static org.junit.jupiter.api.Assertions.*;

@DisplayName("TcpEchoClient")
class TcpEchoClientTest {

    @Test
    @DisplayName("rejects null host")
    void rejectsNullHost() {
        assertThrows(IllegalArgumentException.class,
                () -> new TcpEchoClient(null, 8080, Duration.ofSeconds(5)));
    }

    @Test
    @DisplayName("rejects blank host")
    void rejectsBlankHost() {
        assertThrows(IllegalArgumentException.class,
                () -> new TcpEchoClient("  ", 8080, Duration.ofSeconds(5)));
    }

    @Test
    @DisplayName("rejects invalid port")
    void rejectsInvalidPort() {
        assertThrows(IllegalArgumentException.class,
                () -> new TcpEchoClient("localhost", -1, Duration.ofSeconds(5)));
        assertThrows(IllegalArgumentException.class,
                () -> new TcpEchoClient("localhost", 70000, Duration.ofSeconds(5)));
    }

    @Test
    @DisplayName("rejects null timeout")
    void rejectsNullTimeout() {
        assertThrows(IllegalArgumentException.class,
                () -> new TcpEchoClient("localhost", 8080, null));
    }

    @Test
    @DisplayName("throws when sending without connecting")
    void throwsWhenNotConnected() {
        var client = new TcpEchoClient("localhost", 9999, Duration.ofSeconds(1));
        assertThrows(IllegalStateException.class, () -> client.sendAndReceive("hello"));
    }
}
