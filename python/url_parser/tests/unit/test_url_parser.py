import unittest

from url_parser.implementation.url_parser import URLParser


class TestURLParser(unittest.TestCase):
    def test_simple_http(self) -> None:
        p = URLParser.parse("http://example.com/path")
        self.assertEqual(p.scheme, "http")
        self.assertEqual(p.host, "example.com")
        self.assertEqual(p.path, "/path")
        self.assertEqual(p.port, 80)
        self.assertIsNone(p.query)
        self.assertIsNone(p.fragment)
        self.assertIsNone(p.userinfo)

    def test_https_with_port(self) -> None:
        p = URLParser.parse("https://example.com:8443/api/v1")
        self.assertEqual(p.scheme, "https")
        self.assertEqual(p.host, "example.com")
        self.assertEqual(p.port, 8443)
        self.assertEqual(p.path, "/api/v1")

    def test_with_query_and_fragment(self) -> None:
        p = URLParser.parse("http://example.com/path?q=hello&n=1#section")
        self.assertEqual(p.query, "q=hello&n=1")
        self.assertEqual(p.fragment, "section")
        self.assertEqual(p.path, "/path")

    def test_with_userinfo(self) -> None:
        p = URLParser.parse("ftp://user:pass@ftp.example.com/file")
        self.assertEqual(p.userinfo, "user:pass")
        self.assertEqual(p.host, "ftp.example.com")
        self.assertEqual(p.path, "/file")

    def test_ipv6(self) -> None:
        p = URLParser.parse("http://[::1]:8080/path")
        self.assertEqual(p.host, "::1")
        self.assertEqual(p.port, 8080)
        self.assertEqual(p.path, "/path")

    def test_ipv6_no_port(self) -> None:
        p = URLParser.parse("http://[::1]/path")
        self.assertEqual(p.host, "::1")
        self.assertEqual(p.port, 80)

    def test_empty_path(self) -> None:
        p = URLParser.parse("http://example.com")
        self.assertEqual(p.path, "/")

    def test_no_scheme_relative(self) -> None:
        p = URLParser.parse("/path/to/resource")
        self.assertEqual(p.path, "/path/to/resource")
        self.assertEqual(p.scheme, "")
        self.assertEqual(p.host, "")

    def test_reconstruct_basic(self) -> None:
        url = "http://example.com/path"
        p = URLParser.parse(url)
        self.assertEqual(p.reconstruct(), url)

    def test_reconstruct_with_query_fragment(self) -> None:
        url = "http://example.com/path?a=b#frag"
        p = URLParser.parse(url)
        self.assertEqual(p.reconstruct(), url)

    def test_reconstruct_ipv6(self) -> None:
        url = "http://[::1]:8080/path"
        p = URLParser.parse(url)
        self.assertEqual(p.reconstruct(), url)

    def test_reconstruct_with_userinfo(self) -> None:
        url = "ftp://user@host.com/file"
        p = URLParser.parse(url)
        rec = p.reconstruct()
        self.assertIn("user@", rec)
        self.assertIn("host.com", rec)

    def test_validate_valid(self) -> None:
        p = URLParser.parse("http://example.com")
        self.assertTrue(p.validate())

    def test_validate_missing_scheme(self) -> None:
        p = URLParser.parse("/path")
        self.assertFalse(p.validate())

    def test_validate_bad_port(self) -> None:
        with self.assertRaises(ValueError):
            URLParser.parse("http://example.com:99999")

    def test_validate_default_port(self) -> None:
        p = URLParser.parse("http://example.com")
        self.assertTrue(p.validate())

    def test_empty_url(self) -> None:
        with self.assertRaises(ValueError):
            URLParser.parse("")

    def test_reconstruct_with_default_port(self) -> None:
        p = URLParser.parse("http://example.com:80/path")
        self.assertEqual(p.reconstruct(), "http://example.com/path")

    def test_reconstruct_non_default_port(self) -> None:
        p = URLParser.parse("http://example.com:8080/path")
        self.assertEqual(p.reconstruct(), "http://example.com:8080/path")

    def test_no_host_after_scheme(self) -> None:
        with self.assertRaises(ValueError):
            URLParser.parse("http://")

    def test_missing_scheme_with_colon(self) -> None:
        with self.assertRaises(ValueError):
            URLParser.parse("://host.com")

    def test_unclosed_ipv6(self) -> None:
        with self.assertRaises(ValueError):
            URLParser.parse("http://[::1/path")

    def test_invalid_port_char(self) -> None:
        with self.assertRaises(ValueError):
            URLParser.parse("http://example.com:abc")

    def test_port_zero(self) -> None:
        with self.assertRaises(ValueError):
            URLParser.parse("http://example.com:0")

    def test_encoded_chars(self) -> None:
        url = "http://example.com/path%20with%20spaces?q=%E2%9C%93"
        p = URLParser.parse(url)
        self.assertEqual(p.path, "/path%20with%20spaces")
        self.assertEqual(p.query, "q=%E2%9C%93")

    def test_empty_query(self) -> None:
        p = URLParser.parse("http://example.com/path?")
        self.assertIsNone(p.query)

    def test_empty_fragment(self) -> None:
        p = URLParser.parse("http://example.com/path#")
        self.assertIsNone(p.fragment)

    def test_default_port_known_schemes(self) -> None:
        for scheme, default_port in {"http": 80, "https": 443, "ftp": 21}.items():
            p = URLParser.parse(f"{scheme}://example.com/")
            self.assertEqual(p.port, default_port)


class TestParsedURL(unittest.TestCase):
    def test_default_port_none_when_unknown_scheme(self) -> None:
        p = URLParser.parse("unknown://host.com/path")
        self.assertIsNone(p.port)


if __name__ == "__main__":
    unittest.main()
