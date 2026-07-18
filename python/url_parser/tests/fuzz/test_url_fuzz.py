import random
import string
import unittest

from url_parser.implementation.url_parser import URLParser


def random_string(max_len: int = 100) -> str:
    chars = string.ascii_letters + string.digits + "/:?#[]@!$&'()*+,;=-._~%"
    length = random.randint(0, max_len)
    return "".join(random.choice(chars) for _ in range(length))


class TestURLParserFuzz(unittest.TestCase):
    def test_random_strings_no_crash(self) -> None:
        for _ in range(10000):
            s = random_string()
            try:
                URLParser.parse(s)
            except (ValueError, TypeError):
                pass

    def test_edge_case_strings(self) -> None:
        edge_cases = [
            "",
            "http://",
            "://",
            "http://a",
            "http://a:b@c",
            "http://[::1]:99999",
            "http://[::1]",
            "http://a?",
            "http://a#",
            "http://a?#",
            "http://a?query#frag",
            "a",
            "/",
            "//",
            "://",
            "http://host:port",
            "http://host:0",
            "http://host:65536",
        ]
        for s in edge_cases:
            try:
                URLParser.parse(s)
            except (ValueError, TypeError):
                pass

    def test_valid_urls_roundtrip(self) -> None:
        valid = [
            "http://example.com",
            "https://example.com:8443/path",
            "ftp://user@host.com/file",
            "http://[::1]:8080/path",
            "http://example.com/path?query=value#frag",
            "http://a.b/c/d/e/f",
        ]
        for url in valid:
            try:
                parsed = URLParser.parse(url)
                reconstructed = parsed.reconstruct()
                reparsed = URLParser.parse(reconstructed)
                self.assertEqual(reparsed.scheme, parsed.scheme)
                self.assertEqual(reparsed.host, parsed.host)
                self.assertEqual(reparsed.path, parsed.path)
            except (ValueError, TypeError):
                self.fail(f"Unexpected exception for {url!r}")


if __name__ == "__main__":
    unittest.main()
