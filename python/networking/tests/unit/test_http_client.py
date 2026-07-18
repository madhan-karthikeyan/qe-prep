import unittest
from unittest.mock import MagicMock, patch

from networking.implementation.http_client import HTTPClient, HTTPError


class TestHTTPClient(unittest.TestCase):
    def test_invalid_construction(self) -> None:
        with self.assertRaises(ValueError):
            HTTPClient(max_retries=-1)
        with self.assertRaises(ValueError):
            HTTPClient(base_delay=0)
        with self.assertRaises(ValueError):
            HTTPClient(max_delay=0.1, base_delay=0.5)
        with self.assertRaises(ValueError):
            HTTPClient(timeout=0)

    @patch("networking.implementation.http_client.http.client.HTTPConnection")
    def test_get_success(self, mock_conn_class: MagicMock) -> None:
        mock_response = MagicMock()
        mock_response.status = 200
        mock_response.read.return_value = b'{"ok": true}'
        mock_conn_instance = MagicMock()
        mock_conn_instance.getresponse.return_value = mock_response
        mock_conn_class.return_value = mock_conn_instance

        client = HTTPClient(max_retries=0)
        status, body = client.get("http://example.com/test")
        self.assertEqual(status, 200)
        self.assertEqual(body, '{"ok": true}')

    def test_unsupported_scheme(self) -> None:
        client = HTTPClient(max_retries=0)
        with self.assertRaises(ValueError):
            client.get("ftp://example.com/")


class TestHTTPError(unittest.TestCase):
    def test_http_error_message(self) -> None:
        err = HTTPError(404, "Not Found")
        self.assertEqual(err.status, 404)
        self.assertIn("404", str(err))
        self.assertIn("Not Found", str(err))


if __name__ == "__main__":
    unittest.main()
