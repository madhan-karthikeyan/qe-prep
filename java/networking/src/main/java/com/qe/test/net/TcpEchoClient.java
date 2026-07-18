package com.qe.test.net;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.Socket;
import java.nio.charset.StandardCharsets;
import java.time.Duration;

public class TcpEchoClient implements AutoCloseable {
    private final String host;
    private final int port;
    private final Duration timeout;
    private Socket socket;

    public TcpEchoClient(String host, int port, Duration timeout) {
        if (host == null || host.isBlank()) {
            throw new IllegalArgumentException("host must not be null or blank");
        }
        if (port < 0 || port > 65535) {
            throw new IllegalArgumentException("Invalid port: " + port);
        }
        if (timeout == null || timeout.isNegative()) {
            throw new IllegalArgumentException("timeout must not be null or negative");
        }
        this.host = host;
        this.port = port;
        this.timeout = timeout;
    }

    public void connect() throws IOException {
        socket = new Socket(host, port);
        socket.setSoTimeout((int) timeout.toMillis());
    }

    public String sendAndReceive(String message) throws IOException {
        if (socket == null || !socket.isConnected()) {
            throw new IllegalStateException("Not connected. Call connect() first.");
        }
        OutputStream out = socket.getOutputStream();
        out.write(message.getBytes(StandardCharsets.UTF_8));
        out.flush();

        socket.shutdownOutput();

        InputStream in = socket.getInputStream();
        byte[] buffer = new byte[4096];
        int bytesRead = in.read(buffer);
        if (bytesRead == -1) {
            return "";
        }
        return new String(buffer, 0, bytesRead, StandardCharsets.UTF_8);
    }

    public boolean isConnected() {
        return socket != null && socket.isConnected() && !socket.isClosed();
    }

    @Override
    public void close() throws IOException {
        if (socket != null && !socket.isClosed()) {
            socket.close();
        }
    }
}
