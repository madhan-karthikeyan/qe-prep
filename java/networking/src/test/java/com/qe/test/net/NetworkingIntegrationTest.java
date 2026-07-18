package com.qe.test.net;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;

import java.io.IOException;
import java.time.Duration;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

@DisplayName("Networking Integration")
class NetworkingIntegrationTest {

    @Test
    @DisplayName("TCP echo server echoes message back")
    @Timeout(10)
    void tcpEcho() throws IOException, InterruptedException {
        var server = new TcpEchoServer(0, 2);
        server.start();
        Thread.sleep(100);

        var client = new TcpEchoClient("localhost", server.getPort(), Duration.ofSeconds(5));
        client.connect();

        String response = client.sendAndReceive("Hello, Server!");
        assertEquals("Hello, Server!", response);

        client.close();
        server.stop();
    }

    @Test
    @DisplayName("echo server handles multiple clients")
    @Timeout(15)
    void multipleClients() throws IOException, InterruptedException {
        var server = new TcpEchoServer(0, 5);
        server.start();
        Thread.sleep(100);

        var client1 = new TcpEchoClient("localhost", server.getPort(), Duration.ofSeconds(5));
        var client2 = new TcpEchoClient("localhost", server.getPort(), Duration.ofSeconds(5));
        client1.connect();
        client2.connect();

        assertEquals("msg1", client1.sendAndReceive("msg1"));
        assertEquals("msg2", client2.sendAndReceive("msg2"));

        client1.close();
        client2.close();
        server.stop();
    }
}
