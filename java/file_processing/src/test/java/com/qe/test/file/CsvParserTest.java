package com.qe.test.file;

import static org.junit.jupiter.api.Assertions.*;

import java.io.StringReader;
import java.util.ArrayList;
import java.util.NoSuchElementException;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

@DisplayName("CsvParser")
class CsvParserTest {

    @Test
    @DisplayName("parses simple CSV")
    void simpleCsv() throws Exception {
        var csv = "a,b,c\n1,2,3";
        try (var parser = new CsvParser(new StringReader(csv))) {
            assertArrayEquals(new String[]{"a", "b", "c"}, parser.next());
            assertArrayEquals(new String[]{"1", "2", "3"}, parser.next());
            assertFalse(parser.hasNext());
        }
    }

    @Test
    @DisplayName("parses quoted fields with commas")
    void quotedFields() throws Exception {
        var csv = "\"hello, world\",foo,\"bar\"\"baz\"";
        try (var parser = new CsvParser(new StringReader(csv))) {
            var row = parser.next();
            assertEquals("hello, world", row[0]);
            assertEquals("foo", row[1]);
            assertEquals("bar\"baz", row[2]);
        }
    }

    @Test
    @DisplayName("skips blank lines")
    void skipsBlankLines() throws Exception {
        var csv = "a,b\n\n\nc,d";
        try (var parser = new CsvParser(new StringReader(csv))) {
            assertEquals("a", parser.next()[0]);
            assertEquals("c", parser.next()[0]);
            assertFalse(parser.hasNext());
        }
    }

    @Test
    @DisplayName("custom delimiter")
    void customDelimiter() throws Exception {
        var csv = "a|b|c\n1|2|3";
        try (var parser = new CsvParser(new StringReader(csv), '|', '"')) {
            assertArrayEquals(new String[]{"a", "b", "c"}, parser.next());
            assertArrayEquals(new String[]{"1", "2", "3"}, parser.next());
        }
    }

    @Test
    @DisplayName("empty file returns nothing")
    void emptyFile() throws Exception {
        try (var parser = new CsvParser(new StringReader(""))) {
            assertFalse(parser.hasNext());
        }
    }

    @Test
    @DisplayName("throws NoSuchElementException on empty")
    void throwsOnEmpty() throws Exception {
        try (var parser = new CsvParser(new StringReader(""))) {
            assertThrows(NoSuchElementException.class, parser::next);
        }
    }

    @Test
    @DisplayName("handles unicode")
    void unicode() throws Exception {
        var csv = "姓名,年龄\n张三,28";
        try (var parser = new CsvParser(new StringReader(csv))) {
            assertArrayEquals(new String[]{"姓名", "年龄"}, parser.next());
            assertArrayEquals(new String[]{"张三", "28"}, parser.next());
        }
    }

    @Test
    @DisplayName("collect all rows")
    void collectAll() throws Exception {
        var csv = "x,y\n1,2\n3,4";
        var all = new ArrayList<String[]>();
        try (var parser = new CsvParser(new StringReader(csv))) {
            parser.forEachRemaining(all::add);
        }
        assertEquals(3, all.size());
        assertArrayEquals(new String[]{"x", "y"}, all.get(0));
        assertArrayEquals(new String[]{"1", "2"}, all.get(1));
    }
}
