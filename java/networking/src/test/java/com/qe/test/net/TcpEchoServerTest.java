package com.qe.test.net;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;

import java.io.IOException;
import java.net.Socket;

import static org.junit.jupiter.api.Assertions.*;

@DisplayName("TcpEchoServer")
class TcpEchoServerTest {

    @Test
    @DisplayName("starts and stops without error")
    @Timeout(5)
    void startStop() throws IOException, InterruptedException {
        var server = new TcpEchoServer(0, 2);
        server.start();
        Thread.sleep(100);
        assertTrue(server.isRunning());
        assertTrue(server.getPort() > 0);
        server.stop();
        assertFalse(server.isRunning());
        assertTrue(server.getPort() > 0);
    }

    @Test
    @DisplayName("rejects invalid port")
    void rejectsInvalidPort() {
        assertThrows(IllegalArgumentException.class, () -> new TcpEchoServer(-1, 1));
        assertThrows(IllegalArgumentException.class, () -> new TcpEchoServer(70000, 1));
    }
}
