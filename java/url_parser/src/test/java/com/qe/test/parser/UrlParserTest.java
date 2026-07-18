package com.qe.test.parser;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;

import java.net.URI;
import java.net.URISyntaxException;

import static org.junit.jupiter.api.Assertions.*;

@DisplayName("UrlParser")
class UrlParserTest {

    @Test
    @DisplayName("parses simple HTTP URL")
    void parsesSimpleHttpUrl() {
        var parser = UrlParser.parse("http://example.com/path");
        assertEquals("http", parser.getScheme());
        assertEquals("example.com", parser.getHost());
        assertEquals(80, parser.getPort());
        assertEquals("/path", parser.getPath());
    }

    @Test
    @DisplayName("parses HTTPS URL with port")
    void parsesHttpsWithPort() {
        var parser = UrlParser.parse("https://example.com:8443/api/v1");
        assertEquals("https", parser.getScheme());
        assertEquals("example.com", parser.getHost());
        assertEquals(8443, parser.getPort());
        assertEquals("/api/v1", parser.getPath());
    }

    @Test
    @DisplayName("parses URL with query and fragment")
    void parsesQueryAndFragment() {
        var parser = UrlParser.parse("http://example.com/page?key=value&foo=bar#section");
        assertEquals("key=value&foo=bar", parser.getQuery());
        assertEquals("section", parser.getFragment());
        assertEquals(2, parser.getQueryParams().size());
        assertEquals("value", parser.getQueryParams().get("key"));
        assertEquals("bar", parser.getQueryParams().get("foo"));
    }

    @Test
    @DisplayName("parses URL with user info")
    void parsesUserInfo() {
        var parser = UrlParser.parse("ftp://user:pass@ftp.example.com/file");
        assertEquals("user:pass", parser.getUserInfo());
        assertEquals("ftp.example.com", parser.getHost());
    }

    @Test
    @DisplayName("parses IPv6 address")
    void parsesIpv6() {
        var parser = UrlParser.parse("http://[::1]:8080/path");
        assertEquals("::1", parser.getHost());
        assertEquals(8080, parser.getPort());
        assertEquals("/path", parser.getPath());
    }

    @ParameterizedTest
    @CsvSource({
            "http://example.com, http://example.com",
            "https://example.com:443/path, https://example.com/path",
            "http://example.com:80, http://example.com",
            "http://user@host.com/x, http://user@host.com/x",
            "http://[::1]:8080/p, http://[::1]:8080/p"
    })
    @DisplayName("reconstructs URL correctly")
    void reconstructsUrl(String input, String expected) {
        var parser = UrlParser.parse(input);
        assertEquals(expected, parser.reconstruct());
    }

    @Test
    @DisplayName("matches java.net.URI parsing")
    void matchesJavaNetUri() throws URISyntaxException {
        String url = "https://user:pass@api.example.com:8443/v2/items?page=1#top";
        var parser = UrlParser.parse(url);
        var uri = new URI(url);
        assertEquals(uri.getScheme(), parser.getScheme());
        assertEquals(uri.getHost(), parser.getHost());
        assertEquals(uri.getPort(), parser.getPort());
        assertEquals(uri.getPath(), parser.getPath());
        assertEquals(uri.getQuery(), parser.getQuery());
        assertEquals(uri.getFragment(), parser.getFragment());
    }

    @Test
    @DisplayName("rejects null URL")
    void rejectsNull() {
        assertThrows(NullPointerException.class, () -> UrlParser.parse(null));
    }

    @Test
    @DisplayName("rejects malformed URL")
    void rejectsMalformed() {
        assertThrows(IllegalArgumentException.class, () -> UrlParser.parse("not a url"));
        assertThrows(IllegalArgumentException.class, () -> UrlParser.parse(""));
    }

    @Test
    @DisplayName("returns default port for known schemes")
    void defaultPorts() {
        assertEquals(80, UrlParser.parse("http://example.com").getPort());
        assertEquals(443, UrlParser.parse("https://example.com").getPort());
        assertEquals(21, UrlParser.parse("ftp://example.com").getPort());
        assertEquals(-1, UrlParser.parse("unknown://example.com").getPort());
    }

    @Test
    @DisplayName("handles empty path")
    void emptyPath() {
        var parser = UrlParser.parse("http://example.com");
        assertEquals("", parser.getPath());
    }

    @Test
    @DisplayName("parses URL with encoded characters")
    void encodedCharacters() {
        var parser = UrlParser.parse("http://example.com/path%20with%20spaces?q=hello%20world");
        assertEquals("/path%20with%20spaces", parser.getPath());
        assertEquals("q=hello%20world", parser.getQuery());
    }
}
