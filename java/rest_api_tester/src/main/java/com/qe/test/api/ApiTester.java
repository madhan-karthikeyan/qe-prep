package com.qe.test.api;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;
import java.time.Instant;
import java.util.concurrent.ThreadLocalRandom;

public class ApiTester {
    private final String baseUrl;
    private final Duration timeout;
    private final int maxRetries;
    private final Duration baseDelay;

    public ApiTester(String baseUrl, Duration timeout, int maxRetries, Duration baseDelay) {
        if (baseUrl == null || baseUrl.isBlank()) {
            throw new IllegalArgumentException("baseUrl must not be null or blank");
        }
        if (timeout == null || timeout.isNegative()) {
            throw new IllegalArgumentException("timeout must not be null or negative");
        }
        if (maxRetries < 0) {
            throw new IllegalArgumentException("maxRetries must be non-negative: " + maxRetries);
        }
        this.baseUrl = baseUrl.endsWith("/") ? baseUrl.substring(0, baseUrl.length() - 1) : baseUrl;
        this.timeout = timeout;
        this.maxRetries = maxRetries;
        this.baseDelay = baseDelay != null ? baseDelay : Duration.ofMillis(100);
    }

    public TestResult sendRequest(String method, String path, String body) {
        Instant start = Instant.now();
        try {
            var client = HttpClient.newBuilder()
                    .connectTimeout(timeout)
                    .build();

            HttpRequest request = buildRequest(method, path, body);
            HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());

            return new TestResult(
                    response.statusCode(),
                    response.body(),
                    Duration.between(start, Instant.now()),
                    true,
                    null
            );
        } catch (IOException | InterruptedException e) {
            return new TestResult(-1, null, Duration.between(start, Instant.now()), false, e.getMessage());
        }
    }

    public TestResult sendRequestWithRetry(String method, String path, String body) {
        Instant start = Instant.now();
        Exception lastException = null;

        for (int attempt = 0; attempt <= maxRetries; attempt++) {
            TestResult result = sendRequest(method, path, body);
            if (result.passed()) {
                return result;
            }
            lastException = new RuntimeException(result.errorMessage());
            if (attempt < maxRetries) {
                sleep(computeBackoff(attempt));
            }
        }

        return new TestResult(-1, null, Duration.between(start, Instant.now()), false,
                "Request failed after " + (maxRetries + 1) + " attempts: " + lastException.getMessage());
    }

    public void assertStatus(TestResult result, int expectedStatus) {
        if (result.statusCode() != expectedStatus) {
            throw new AssertionError("Expected status " + expectedStatus
                    + " but got " + result.statusCode()
                    + ". Body: " + result.body());
        }
    }

    public void assertJson(TestResult result, String expectedKey, String expectedValue) {
        if (result.body() == null || !result.body().contains("\"" + expectedKey + "\"")) {
            throw new AssertionError("Expected key '" + expectedKey + "' not found in body: " + result.body());
        }
        // Simple JSON value extraction
        String search = "\"" + expectedKey + "\":";
        int keyIdx = result.body().indexOf(search);
        if (keyIdx == -1) {
            throw new AssertionError("Key '" + expectedKey + "' not found in body: " + result.body());
        }
        int valueStart = keyIdx + search.length();
        if (valueStart >= result.body().length()) {
            throw new AssertionError("Unexpected body format: " + result.body());
        }
        char firstChar = result.body().charAt(valueStart);
        if (firstChar == '"') {
            int valueEnd = result.body().indexOf('"', valueStart + 1);
            if (valueEnd == -1) {
                throw new AssertionError("Unterminated string value in body: " + result.body());
            }
            String actualValue = result.body().substring(valueStart + 1, valueEnd);
            if (!actualValue.equals(expectedValue)) {
                throw new AssertionError("Expected value '" + expectedValue + "' but got '" + actualValue + "'");
            }
        }
    }

    private HttpRequest buildRequest(String method, String path, String body) {
        var builder = HttpRequest.newBuilder()
                .uri(URI.create(baseUrl + path))
                .timeout(timeout);

        return switch (method.toUpperCase()) {
            case "GET" -> builder.GET().build();
            case "DELETE" -> builder.DELETE().build();
            case "POST" -> {
                if (body != null) {
                    yield builder.POST(HttpRequest.BodyPublishers.ofString(body))
                            .header("Content-Type", "application/json")
                            .build();
                }
                yield builder.POST(HttpRequest.BodyPublishers.noBody()).build();
            }
            case "PUT" -> {
                if (body != null) {
                    yield builder.PUT(HttpRequest.BodyPublishers.ofString(body))
                            .header("Content-Type", "application/json")
                            .build();
                }
                yield builder.PUT(HttpRequest.BodyPublishers.noBody()).build();
            }
            default -> throw new IllegalArgumentException("Unsupported method: " + method);
        };
    }

    private long computeBackoff(int attempt) {
        long delay = baseDelay.toMillis() * (1L << attempt);
        double jitter = 1.0 + (ThreadLocalRandom.current().nextDouble() * 2 - 1) * 0.2;
        return (long) (delay * jitter);
    }

    private void sleep(long millis) {
        try {
            Thread.sleep(millis);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }

    public record TestResult(int statusCode, String body, Duration duration, boolean passed, String errorMessage) {
        public TestResult {
            if (duration == null) {
                throw new NullPointerException("duration must not be null");
            }
        }
    }
}
