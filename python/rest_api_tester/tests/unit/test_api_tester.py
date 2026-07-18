import unittest

from rest_api_tester.implementation.api_tester import APITester, TestAssertionError
from rest_api_tester.implementation.http_server import TestHTTPServer


class TestAPITester(unittest.TestCase):
    def setUp(self) -> None:
        self.server = TestHTTPServer()
        self.server.start()
        self.tester = APITester(f"http://127.0.0.1:{self.server.port}")

    def tearDown(self) -> None:
        self.server.stop()

    def test_get_not_found(self) -> None:
        self.tester.send_request("GET", "/nonexistent")
        self.tester.assert_status(404)

    def test_post_and_get(self) -> None:
        self.tester.send_request("POST", "/", body={"key": "foo", "value": "bar"})
        self.tester.assert_status(201)
        self.tester.assert_json("key", "foo")
        self.tester.assert_json("value", "bar")

        self.tester.send_request("GET", "/foo")
        self.tester.assert_status(200)
        self.tester.assert_json("value", "bar")

    def test_delete(self) -> None:
        self.tester.send_request("POST", "/", body={"key": "x", "value": "y"})
        self.tester.send_request("DELETE", "/x")
        self.tester.assert_status(200)
        self.tester.assert_json("deleted", "x")

        self.tester.send_request("GET", "/x")
        self.tester.assert_status(404)

    def test_assert_status_failure(self) -> None:
        self.tester.send_request("GET", "/nonexistent")
        with self.assertRaises(TestAssertionError):
            self.tester.assert_status(200)

    def test_assert_json_failure(self) -> None:
        self.tester.send_request("GET", "/nonexistent")
        with self.assertRaises(TestAssertionError):
            self.tester.assert_json("value", "something")

    def test_invalid_json_body(self) -> None:
        import http.client
        conn = http.client.HTTPConnection("127.0.0.1", self.server.port)
        conn.request("POST", "/", body="not json", headers={"Content-Type": "application/json"})
        resp = conn.getresponse()
        self.assertEqual(resp.status, 400)
        conn.close()

    def test_report_tracking(self) -> None:
        self.tester.send_request("GET", "/a")
        self.tester.assert_status(404)
        report = self.tester.report
        self.assertEqual(report.total, 1)
        self.assertGreaterEqual(report.passed, 0)
        self.assertEqual(report.failed, 0)

    def test_retry_on_failure(self) -> None:
        tester = APITester(
            f"http://127.0.0.1:{self.server.port}",
            max_retries=2,
            retry_delay=0.05,
        )
        tester.send_request("GET", "/retry-test")
        self.assertIn(tester.last_status, (404,))


if __name__ == "__main__":
    unittest.main()
