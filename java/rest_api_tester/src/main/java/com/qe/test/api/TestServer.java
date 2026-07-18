package com.qe.test.api;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.regex.Pattern;

public class TestServer implements AutoCloseable {
    private final HttpServer server;
    private final ExecutorService executor;
    private final Map<String, String> store = new ConcurrentHashMap<>();
    private final AtomicInteger idCounter = new AtomicInteger(1);

    public TestServer(int port) throws IOException {
        server = HttpServer.create(new InetSocketAddress(port), 0);
        executor = Executors.newFixedThreadPool(4);

        server.createContext("/api/items", new ItemsHandler());
        server.createContext("/api/items/", new ItemByIdHandler());

        server.setExecutor(executor);
        server.start();
    }

    public void stop() {
        server.stop(0);
        executor.shutdown();
        try {
            if (!executor.awaitTermination(5, TimeUnit.SECONDS)) {
                executor.shutdownNow();
            }
        } catch (InterruptedException e) {
            executor.shutdownNow();
            Thread.currentThread().interrupt();
        }
    }

    public int getPort() {
        return server.getAddress().getPort();
    }

    @Override
    public void close() {
        stop();
    }

    private class ItemsHandler implements HttpHandler {
        @Override
        public void handle(HttpExchange exchange) throws IOException {
            try (exchange) {
                String method = exchange.getRequestMethod();
                switch (method) {
                    case "GET" -> handleGetAll(exchange);
                    case "POST" -> handlePost(exchange);
                    default -> sendJson(exchange, 405, "{\"error\":\"Method not allowed\"}");
                }
            }
        }

        private void handleGetAll(HttpExchange exchange) throws IOException {
            String json = store.entrySet().stream()
                    .map(e -> "{\"id\":" + e.getKey() + ",\"data\":\"" + escapeJson(e.getValue()) + "\"}")
                    .toList()
                    .toString();
            sendJson(exchange, 200, "{\"items\":" + json + "}");
        }

        private void handlePost(HttpExchange exchange) throws IOException {
            String body = new String(exchange.getRequestBody().readAllBytes(), StandardCharsets.UTF_8);
            String id = String.valueOf(idCounter.getAndIncrement());
            store.put(id, body);
            sendJson(exchange, 201, "{\"id\":" + id + ",\"data\":\"" + escapeJson(body) + "\"}");
        }
    }

    private class ItemByIdHandler implements HttpHandler {
        private final Pattern idPattern = Pattern.compile("/api/items/(\\d+)");

        @Override
        public void handle(HttpExchange exchange) throws IOException {
            try (exchange) {
                var matcher = idPattern.matcher(exchange.getRequestURI().getPath());
                if (!matcher.matches()) {
                    sendJson(exchange, 400, "{\"error\":\"Invalid path\"}");
                    return;
                }
                String id = matcher.group(1);
                String method = exchange.getRequestMethod();
                switch (method) {
                    case "GET" -> handleGet(exchange, id);
                    case "DELETE" -> handleDelete(exchange, id);
                    default -> sendJson(exchange, 405, "{\"error\":\"Method not allowed\"}");
                }
            }
        }

        private void handleGet(HttpExchange exchange, String id) throws IOException {
            String data = store.get(id);
            if (data == null) {
                sendJson(exchange, 404, "{\"error\":\"Not found\"}");
                return;
            }
            sendJson(exchange, 200, "{\"id\":" + id + ",\"data\":\"" + escapeJson(data) + "\"}");
        }

        private void handleDelete(HttpExchange exchange, String id) throws IOException {
            String removed = store.remove(id);
            if (removed == null) {
                sendJson(exchange, 404, "{\"error\":\"Not found\"}");
                return;
            }
            sendJson(exchange, 200, "{\"deleted\":" + id + "}");
        }
    }

    private void sendJson(HttpExchange exchange, int status, String json) throws IOException {
        byte[] bytes = json.getBytes(StandardCharsets.UTF_8);
        exchange.getResponseHeaders().set("Content-Type", "application/json");
        exchange.sendResponseHeaders(status, bytes.length);
        try (OutputStream os = exchange.getResponseBody()) {
            os.write(bytes);
        }
    }

    static String escapeJson(String s) {
        return s.replace("\\", "\\\\")
                .replace("\"", "\\\"")
                .replace("\n", "\\n")
                .replace("\r", "\\r")
                .replace("\t", "\\t");
    }
}
