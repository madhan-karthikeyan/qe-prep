package com.qe.test.rotation;

import java.io.Closeable;
import java.io.Flushable;
import java.io.IOException;
import java.io.Writer;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardOpenOption;
import java.time.Instant;
import java.time.ZoneId;
import java.time.format.DateTimeFormatter;

public class RotatingFileWriter implements Closeable, Flushable {
    private static final DateTimeFormatter ARCHIVE_TIMESTAMP =
            DateTimeFormatter.ofPattern("yyyyMMdd-HHmmss").withZone(ZoneId.of("UTC"));

    private final Path basePath;
    private final long maxBytes;
    private Writer writer;
    private long bytesWritten;
    private boolean closed;

    public RotatingFileWriter(Path basePath, long maxBytes) {
        if (basePath == null) throw new IllegalArgumentException("basePath must not be null");
        if (maxBytes <= 0) throw new IllegalArgumentException("maxBytes must be positive");
        this.basePath = basePath.toAbsolutePath();
        this.maxBytes = maxBytes;
    }

    public void write(String data) throws IOException {
        if (closed) throw new IllegalStateException("writer is closed");
        byte[] bytes = data.getBytes(StandardCharsets.UTF_8);
        if (bytesWritten + bytes.length > maxBytes) {
            rotate();
        }
        ensureOpen();
        writer.write(data);
        writer.flush();
        bytesWritten += bytes.length;
    }

    private void ensureOpen() throws IOException {
        if (writer == null) {
            Files.createDirectories(basePath.getParent());
            writer = Files.newBufferedWriter(basePath, StandardCharsets.UTF_8,
                    StandardOpenOption.CREATE, StandardOpenOption.APPEND);
        }
    }

    private void rotate() throws IOException {
        closeWriter();
        Path parent = basePath.getParent();
        String name = basePath.getFileName().toString();
        String ts = ARCHIVE_TIMESTAMP.format(Instant.now());
        Path archive = parent.resolve(name + "." + ts);
        Files.move(basePath, archive);
        bytesWritten = 0;
        writer = null;
    }

    private void closeWriter() throws IOException {
        if (writer != null) {
            try {
                writer.close();
            } finally {
                writer = null;
            }
        }
    }

    @Override
    public void flush() throws IOException {
        if (writer != null) {
            writer.flush();
        }
    }

    @Override
    public void close() throws IOException {
        if (!closed) {
            closed = true;
            closeWriter();
        }
    }

    public Path basePath() { return basePath; }
    public long maxBytes() { return maxBytes; }
}
