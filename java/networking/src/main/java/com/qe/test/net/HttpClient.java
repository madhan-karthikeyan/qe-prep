package com.qe.test.net;

import java.io.IOException;
import java.io.InputStream;
import java.net.HttpURLConnection;
import java.net.URI;
import java.net.URISyntaxException;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.concurrent.ThreadLocalRandom;

public class HttpClient {
    private final Duration baseDelay;
    private final Duration maxDelay;
    private final double jitterFactor;
    private final int maxRetries;

    public HttpClient(Duration baseDelay, Duration maxDelay, double jitterFactor, int maxRetries) {
        if (baseDelay == null || baseDelay.isNegative()) {
            throw new IllegalArgumentException("baseDelay must not be null or negative");
        }
        if (maxDelay == null || maxDelay.isNegative()) {
            throw new IllegalArgumentException("maxDelay must not be null or negative");
        }
        if (jitterFactor < 0 || jitterFactor > 1) {
            throw new IllegalArgumentException("jitterFactor must be in [0, 1]: " + jitterFactor);
        }
        if (maxRetries < 0) {
            throw new IllegalArgumentException("maxRetries must be non-negative: " + maxRetries);
        }
        this.baseDelay = baseDelay;
        this.maxDelay = maxDelay;
        this.jitterFactor = jitterFactor;
        this.maxRetries = maxRetries;
    }

    public HttpResponse get(String url) throws IOException, InterruptedException {
        return executeWithRetry(url, "GET", null);
    }

    public HttpResponse post(String url, String body) throws IOException, InterruptedException {
        return executeWithRetry(url, "POST", body);
    }

    public HttpResponse delete(String url) throws IOException, InterruptedException {
        return executeWithRetry(url, "DELETE", null);
    }

    private HttpResponse executeWithRetry(String url, String method, String body)
            throws IOException, InterruptedException {
        IOException lastException = null;
        for (int attempt = 0; attempt <= maxRetries; attempt++) {
            try {
                return execute(url, method, body);
            } catch (IOException e) {
                lastException = e;
                if (attempt < maxRetries) {
                    long delay = computeDelay(attempt);
                    Thread.sleep(delay);
                }
            }
        }
        throw new IOException("Request failed after " + (maxRetries + 1) + " attempts", lastException);
    }

    private HttpResponse execute(String url, String method, String body) throws IOException {
        try {
            HttpURLConnection conn = (HttpURLConnection) new URI(url).toURL().openConnection();
            conn.setRequestMethod(method);
            conn.setConnectTimeout(5000);
            conn.setReadTimeout(5000);

            if (body != null && !body.isEmpty()) {
                conn.setDoOutput(true);
                conn.getOutputStream().write(body.getBytes(StandardCharsets.UTF_8));
            }

            int statusCode = conn.getResponseCode();
            String responseBody;
            try (InputStream in = (statusCode >= 400 ? conn.getErrorStream() : conn.getInputStream())) {
                if (in == null) {
                    responseBody = "";
                } else {
                    responseBody = new String(in.readAllBytes(), StandardCharsets.UTF_8);
                }
            }
            conn.disconnect();
            return new HttpResponse(statusCode, responseBody);
        } catch (URISyntaxException e) {
            throw new IOException("Invalid URL: " + url, e);
        }
    }

    private long computeDelay(int attempt) {
        long exponentialDelay = Math.min(baseDelay.toMillis() * (1L << attempt), maxDelay.toMillis());
        if (jitterFactor > 0) {
            double jitter = 1.0 + (ThreadLocalRandom.current().nextDouble() * 2 - 1) * jitterFactor;
            exponentialDelay = (long) (exponentialDelay * jitter);
        }
        return Math.max(0, Math.min(exponentialDelay, maxDelay.toMillis()));
    }

    public record HttpResponse(int statusCode, String body) { }
}
