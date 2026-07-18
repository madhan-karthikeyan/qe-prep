package com.qe.test.net;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.ServerSocket;
import java.net.Socket;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;

public class TcpEchoServer implements AutoCloseable {
    private final int port;
    private final ExecutorService threadPool;
    private ServerSocket serverSocket;
    private final AtomicBoolean running = new AtomicBoolean(false);
    private Thread acceptThread;

    public TcpEchoServer(int port, int maxThreads) {
        if (port < 0 || port > 65535) {
            throw new IllegalArgumentException("Invalid port: " + port);
        }
        this.port = port;
        this.threadPool = Executors.newFixedThreadPool(maxThreads);
    }

    public void start() throws IOException {
        if (running.compareAndSet(false, true)) {
            serverSocket = new ServerSocket(port);
            acceptThread = Thread.ofVirtual().start(() -> {
                while (running.get() && !serverSocket.isClosed()) {
                    try {
                        Socket clientSocket = serverSocket.accept();
                        threadPool.submit(() -> handleClient(clientSocket));
                    } catch (IOException e) {
                        if (running.get()) {
                            System.err.println("Accept error: " + e.getMessage());
                        }
                    }
                }
            });
        }
    }

    public void stop() {
        running.set(false);
        try {
            if (serverSocket != null && !serverSocket.isClosed()) {
                serverSocket.close();
            }
        } catch (IOException e) {
            System.err.println("Error closing server: " + e.getMessage());
        }
        threadPool.shutdown();
        try {
            if (!threadPool.awaitTermination(5, TimeUnit.SECONDS)) {
                threadPool.shutdownNow();
            }
        } catch (InterruptedException e) {
            threadPool.shutdownNow();
            Thread.currentThread().interrupt();
        }
    }

    public int getPort() {
        if (serverSocket != null && serverSocket.isBound()) {
            return serverSocket.getLocalPort();
        }
        return port;
    }

    public boolean isRunning() {
        return running.get();
    }

    private void handleClient(Socket clientSocket) {
        try (clientSocket;
             InputStream in = clientSocket.getInputStream();
             OutputStream out = clientSocket.getOutputStream()) {
            byte[] buffer = new byte[4096];
            int bytesRead;
            while ((bytesRead = in.read(buffer)) != -1) {
                out.write(buffer, 0, bytesRead);
                out.flush();
            }
        } catch (IOException e) {
            if (running.get()) {
                System.err.println("Client handler error: " + e.getMessage());
            }
        }
    }

    @Override
    public void close() {
        stop();
    }
}
