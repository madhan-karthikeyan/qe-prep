package com.qe.test.config;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

@DisplayName("IniParser")
class IniParserTest {

    private final IniParser parser = new IniParser();

    @Test
    @DisplayName("parses section with key=value")
    void parsesSection() {
        var result = parser.parse("[database]\nhost=localhost");
        assertEquals("localhost", result.get("database").get("host"));
    }

    @Test
    @DisplayName("parses multiple sections")
    void parsesMultipleSections() {
        var result = parser.parse("""
                [database]
                host=localhost
                port=5432
                                
                [app]
                name=myapp
                """);
        assertEquals("localhost", result.get("database").get("host"));
        assertEquals("5432", result.get("database").get("port"));
        assertEquals("myapp", result.get("app").get("name"));
    }

    @Test
    @DisplayName("skips comments")
    void skipsComments() {
        var result = parser.parse("""
                ; comment
                # also comment
                [section]
                key=value
                """);
        assertEquals("value", result.get("section").get("key"));
    }

    @Test
    @DisplayName("handles quoted values")
    void quotedValues() {
        var result = parser.parse("[sec]\nname=\"John Doe\"");
        assertEquals("John Doe", result.get("sec").get("name"));
    }

    @Test
    @DisplayName("handles keys without section")
    void keysWithoutSection() {
        var result = parser.parse("key=value");
        assertEquals("value", result.get("").get("key"));
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
    @DisplayName("handles multi-line values")
    void multiLineValues() {
        var result = parser.parse("[sec]\nkey=line1\\\nline2");
        assertEquals("line1\nline2", result.get("sec").get("key"));
    }

    @Test
    @DisplayName("handles whitespace around keys and values")
    void whitespaceHandling() {
        var result = parser.parse("[sec]\n  key =  value  ");
        assertEquals("value", result.get("sec").get("key"));
    }
}
