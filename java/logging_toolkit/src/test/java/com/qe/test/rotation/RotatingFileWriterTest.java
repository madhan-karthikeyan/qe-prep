package com.qe.test.rotation;

import static org.junit.jupiter.api.Assertions.*;

import java.nio.file.Files;
import java.nio.file.Path;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

@DisplayName("RotatingFileWriter")
class RotatingFileWriterTest {

    @Test
    @DisplayName("writes data to file")
    void writesData() throws Exception {
        Path tmp = Files.createTempFile("rfw-test", ".log");
        try (var w = new RotatingFileWriter(tmp, 1024)) {
            w.write("hello");
            w.flush();
            assertEquals("hello", Files.readString(tmp));
        } finally {
            Files.deleteIfExists(tmp);
        }
    }

    @Test
    @DisplayName("rotates when exceeding maxBytes")
    void rotatesOnSize() throws Exception {
        Path tmp = Files.createTempFile("rfw-rotate", ".log");
        try (var w = new RotatingFileWriter(tmp, 10)) {
            w.write("1234567890");
            assertEquals("1234567890", Files.readString(tmp));
            w.write("more data");
            // original should have been renamed
            assertTrue(Files.exists(tmp)); // fresh rotated file
            String content = Files.readString(tmp);
            assertEquals("more data", content);
        } finally {
            Path parent = tmp.getParent();
            String name = tmp.getFileName().toString();
            try (var list = Files.list(parent)) {
                list.filter(p -> p.getFileName().toString().startsWith(name))
                    .forEach(p -> {
                        try { Files.deleteIfExists(p); } catch (Exception ignored) {}
                    });
            }
        }
    }

    @Test
    @DisplayName("throws when writing after close")
    void throwsAfterClose() throws Exception {
        Path tmp = Files.createTempFile("rfw-closed", ".log");
        var w = new RotatingFileWriter(tmp, 1024);
        w.close();
        assertThrows(IllegalStateException.class, () -> w.write("data"));
        Files.deleteIfExists(tmp);
    }

    @Test
    @DisplayName("constructor rejects null path")
    void constructorValidation() {
        assertThrows(IllegalArgumentException.class, () -> new RotatingFileWriter(null, 1024));
        assertThrows(IllegalArgumentException.class, () -> new RotatingFileWriter(Path.of("test"), -1));
        assertThrows(IllegalArgumentException.class, () -> new RotatingFileWriter(Path.of("test"), 0));
    }
}
