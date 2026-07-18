package com.qe.test.file;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.Reader;
import java.util.Objects;

public class WordCounter {
    private long lines;
    private long words;
    private long characters;

    public WordCounter() {
    }

    public WordCounter count(Reader reader) throws IOException {
        Objects.requireNonNull(reader, "reader must not be null");
        var br = (reader instanceof BufferedReader) ? (BufferedReader) reader : new BufferedReader(reader);
        String line;
        while ((line = br.readLine()) != null) {
            lines++;
            characters += line.length() + 1; // +1 for newline
            countWords(line);
        }
        return this;
    }

    private void countWords(String line) {
        boolean inWord = false;
        for (int i = 0; i < line.length(); i++) {
            char c = line.charAt(i);
            if (Character.isWhitespace(c)) {
                inWord = false;
            } else if (!inWord) {
                inWord = true;
                words++;
            }
        }
    }

    public long lines() { return lines; }
    public long words() { return words; }
    public long characters() { return characters; }
}
