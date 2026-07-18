package com.qe.test.file;

import static org.junit.jupiter.api.Assertions.*;

import java.io.StringReader;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

@DisplayName("WordCounter")
class WordCounterTest {

    @Test
    @DisplayName("counts lines, words, chars for simple text")
    void simpleCount() throws Exception {
        var text = "hello world\nfoo bar baz\n";
        var wc = new WordCounter().count(new StringReader(text));
        assertEquals(2, wc.lines());
        assertEquals(5, wc.words());
        assertEquals(24, wc.characters()); // "hello world\nfoo bar baz\n" = 11 + 1 + 11 + 1 = 24
    }

    @Test
    @DisplayName("empty input")
    void emptyInput() throws Exception {
        var wc = new WordCounter().count(new StringReader(""));
        assertEquals(0, wc.lines());
        assertEquals(0, wc.words());
        assertEquals(0, wc.characters());
    }

    @Test
    @DisplayName("single line no words")
    void singleLineNoWords() throws Exception {
        var wc = new WordCounter().count(new StringReader("   \n"));
        assertEquals(1, wc.lines());
        assertEquals(0, wc.words());
        assertEquals(4, wc.characters());
    }

    @Test
    @DisplayName("handles unicode words")
    void unicodeWords() throws Exception {
        var text = "café résumé\n";
        var wc = new WordCounter().count(new StringReader(text));
        assertEquals(1, wc.lines());
        assertEquals(2, wc.words());
    }

    @Test
    @DisplayName("multiple calls accumulate")
    void multipleCalls() throws Exception {
        var wc = new WordCounter();
        wc.count(new StringReader("a b\n"));
        wc.count(new StringReader("c d e\n"));
        assertEquals(2, wc.lines());
        assertEquals(5, wc.words());
    }
}
