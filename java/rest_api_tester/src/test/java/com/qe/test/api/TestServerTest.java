package com.qe.test.api;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;

import static org.junit.jupiter.api.Assertions.*;

@DisplayName("TestServer")
class TestServerTest {

    @Test
    @DisplayName("starts and stops")
    void startStop() throws IOException {
        var server = new TestServer(0);
        assertTrue(server.getPort() > 0);
        server.stop();
    }

    @Test
    @DisplayName("handles GET /api/items")
    @Timeout(10)
    void getAllItems() throws IOException, InterruptedException {
        var server = new TestServer(0);
        var client = HttpClient.newBuilder()
                .connectTimeout(Duration.ofSeconds(5))
                .build();
        var request = HttpRequest.newBuilder()
                .uri(URI.create("http://localhost:" + server.getPort() + "/api/items"))
                .GET()
                .build();
        var response = client.send(request, HttpResponse.BodyHandlers.ofString());
        assertEquals(200, response.statusCode());
        server.stop();
    }

    @Test
    @DisplayName("handles POST and GET by id")
    @Timeout(10)
    void postAndGetById() throws IOException, InterruptedException {
        var server = new TestServer(0);
        var client = HttpClient.newBuilder()
                .connectTimeout(Duration.ofSeconds(5))
                .build();

        var postReq = HttpRequest.newBuilder()
                .uri(URI.create("http://localhost:" + server.getPort() + "/api/items"))
                .POST(HttpRequest.BodyPublishers.ofString("{\"name\":\"test\"}"))
                .header("Content-Type", "application/json")
                .build();
        var postResp = client.send(postReq, HttpResponse.BodyHandlers.ofString());
        assertEquals(201, postResp.statusCode());

        var getReq = HttpRequest.newBuilder()
                .uri(URI.create("http://localhost:" + server.getPort() + "/api/items/1"))
                .GET()
                .build();
        var getResp = client.send(getReq, HttpResponse.BodyHandlers.ofString());
        assertEquals(200, getResp.statusCode());
        server.stop();
    }

    @Test
    @DisplayName("handles DELETE")
    @Timeout(10)
    void deleteItem() throws IOException, InterruptedException {
        var server = new TestServer(0);
        var client = HttpClient.newBuilder()
                .connectTimeout(Duration.ofSeconds(5))
                .build();

        client.send(HttpRequest.newBuilder()
                .uri(URI.create("http://localhost:" + server.getPort() + "/api/items"))
                .POST(HttpRequest.BodyPublishers.ofString("{\"name\":\"delete-me\"}"))
                .header("Content-Type", "application/json")
                .build(), HttpResponse.BodyHandlers.ofString());

        var delReq = HttpRequest.newBuilder()
                .uri(URI.create("http://localhost:" + server.getPort() + "/api/items/1"))
                .DELETE()
                .build();
        var delResp = client.send(delReq, HttpResponse.BodyHandlers.ofString());
        assertEquals(200, delResp.statusCode());

        var getReq = HttpRequest.newBuilder()
                .uri(URI.create("http://localhost:" + server.getPort() + "/api/items/1"))
                .GET()
                .build();
        var getResp = client.send(getReq, HttpResponse.BodyHandlers.ofString());
        assertEquals(404, getResp.statusCode());
        server.stop();
    }
}
