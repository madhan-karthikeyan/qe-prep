package com.qe.test.file;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.Reader;
import java.io.UncheckedIOException;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.NoSuchElementException;
import java.util.Objects;

public class CsvParser implements Iterator<String[]>, AutoCloseable {
    private final BufferedReader reader;
    private final char delimiter;
    private final char quoteChar;
    private String[] next;
    private boolean closed;

    public CsvParser(Reader reader) {
        this(reader, ',', '"');
    }

    public CsvParser(Reader reader, char delimiter, char quoteChar) {
        Objects.requireNonNull(reader, "reader must not be null");
        this.reader = (reader instanceof BufferedReader) ? (BufferedReader) reader : new BufferedReader(reader);
        this.delimiter = delimiter;
        this.quoteChar = quoteChar;
        advance();
    }

    @Override
    public boolean hasNext() {
        return next != null;
    }

    @Override
    public String[] next() {
        if (next == null) throw new NoSuchElementException();
        var result = next;
        advance();
        return result;
    }

    private void advance() {
        try {
            String line = reader.readLine();
            while (line != null && line.isBlank()) {
                line = reader.readLine();
            }
            if (line == null) {
                next = null;
                return;
            }
            if (closed) {
                next = null;
                return;
            }
            next = parseLine(line);
        } catch (IOException e) {
            throw new UncheckedIOException(e);
        }
    }

    private String[] parseLine(String line) {
        var fields = new ArrayList<String>();
        var cur = new StringBuilder();
        boolean inQuotes = false;
        for (int i = 0; i < line.length(); i++) {
            char c = line.charAt(i);
            if (inQuotes) {
                if (c == quoteChar) {
                    if (i + 1 < line.length() && line.charAt(i + 1) == quoteChar) {
                        cur.append(quoteChar);
                        i++;
                    } else {
                        inQuotes = false;
                    }
                } else {
                    cur.append(c);
                }
            } else {
                if (c == quoteChar) {
                    inQuotes = true;
                } else if (c == delimiter) {
                    fields.add(cur.toString());
                    cur = new StringBuilder();
                } else {
                    cur.append(c);
                }
            }
        }
        fields.add(cur.toString());
        return fields.toArray(String[]::new);
    }

    @Override
    public void close() throws IOException {
        closed = true;
        reader.close();
    }
}
