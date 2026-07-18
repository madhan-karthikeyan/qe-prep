package com.qe.test.config;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

@DisplayName("EnvParser")
class EnvParserTest {

    private final EnvParser parser = new EnvParser();

    @Test
    @DisplayName("parses simple key=value")
    void parsesSimple() {
        var result = parser.parse("KEY=value");
        assertEquals("value", result.get("KEY"));
    }

    @Test
    @DisplayName("parses multiple lines")
    void parsesMultipleLines() {
        var result = parser.parse("""
                DB_HOST=localhost
                DB_PORT=5432
                """);
        assertEquals("localhost", result.get("DB_HOST"));
        assertEquals("5432", result.get("DB_PORT"));
    }

    @Test
    @DisplayName("skips comments and empty lines")
    void skipsComments() {
        var result = parser.parse("""
                # This is a comment
                KEY=value
                                
                ANOTHER=val
                """);
        assertEquals(2, result.size());
        assertEquals("value", result.get("KEY"));
        assertEquals("val", result.get("ANOTHER"));
    }

    @Test
    @DisplayName("handles quoted values")
    void quotedValues() {
        var result = parser.parse("KEY=\"quoted value\"");
        assertEquals("quoted value", result.get("KEY"));
    }

    @Test
    @DisplayName("handles values with equals sign")
    void valuesWithEquals() {
        var result = parser.parse("KEY=value=with=equals");
        assertEquals("value=with=equals", result.get("KEY"));
    }

    @Test
    @DisplayName("rejects invalid key format")
    void rejectsInvalidKey() {
        assertThrows(IllegalArgumentException.class, () -> parser.parse("123INVALID=value"));
    }

    @Test
    @DisplayName("rejects null content")
    void rejectsNullContent() {
        assertThrows(NullPointerException.class, () -> parser.parse((String) null));
    }

    @Test
    @DisplayName("rejects null path")
    void rejectsNullPath() {
        assertThrows(NullPointerException.class, () -> parser.parse((java.nio.file.Path) null));
    }

    @Test
    @DisplayName("handles empty value")
    void emptyValue() {
        var result = parser.parse("EMPTY=");
        assertEquals("", result.get("EMPTY"));
    }
}
