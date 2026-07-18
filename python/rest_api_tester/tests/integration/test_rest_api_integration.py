import unittest

from rest_api_tester.implementation.api_tester import APITester, TestAssertionError
from rest_api_tester.implementation.http_server import TestHTTPServer


class TestRESTAPIIntegration(unittest.TestCase):
    def setUp(self) -> None:
        self.server = TestHTTPServer()
        self.server.start()
        self.tester = APITester(f"http://127.0.0.1:{self.server.port}")

    def tearDown(self) -> None:
        self.server.stop()

    def test_full_api_workflow(self) -> None:
        self.tester.send_request("GET", "/items")
        self.tester.assert_status(404)

        self.tester.send_request(
            "POST", "/",
            body={"key": "item1", "value": {"name": "Widget", "price": 9.99}},
        )
        self.tester.assert_status(201)
        self.tester.assert_json("key", "item1")

        self.tester.send_request("GET", "/item1")
        self.tester.assert_status(200)
        self.tester.assert_json("value.name", "Widget")
        self.tester.assert_json("value.price", 9.99)

        self.tester.send_request(
            "POST", "/",
            body={"key": "item2", "value": {"name": "Gadget", "price": 24.99}},
        )
        self.tester.assert_status(201)

        self.tester.send_request("DELETE", "/item1")
        self.tester.assert_status(200)
        self.tester.assert_json("deleted", "item1")

        self.tester.send_request("GET", "/item1")
        self.tester.assert_status(404)

    def test_report_summary(self) -> None:
        self.tester.send_request("GET", "/missing")
        self.tester.assert_status(404)

        self.tester.send_request("POST", "/", body={"key": "k", "value": "v"})
        self.tester.assert_status(201)

        self.tester.send_request("GET", "/k")
        self.tester.assert_status(200)
        self.tester.assert_json("value", "v")

        report = self.tester.report
        self.assertGreaterEqual(report.total, 3)
        self.assertGreater(report.timing, 0)

    def test_multiple_data_types(self) -> None:
        test_data = {
            "key": "types",
            "value": {
                "string": "hello",
                "int": 42,
                "float": 3.14,
                "bool": True,
                "null_value": None,
                "list": [1, 2, 3],
            },
        }

        self.tester.send_request("POST", "/", body=test_data)
        self.tester.assert_status(201)

        self.tester.send_request("GET", "/types")
        self.tester.assert_status(200)
        self.tester.assert_json("value.string", "hello")
        self.tester.assert_json("value.int", 42)
        self.tester.assert_json("value.float", 3.14)
        self.tester.assert_json("value.bool", True)
        self.tester.assert_json("value.null_value", None)
        self.tester.assert_json("value.list", [1, 2, 3])

    def test_error_status_assertions(self) -> None:
        self.tester.send_request("DELETE", "/nonexistent")
        self.tester.assert_status(404)

        self.tester.send_request("POST", "/", body={"bad": "request"})
        self.tester.assert_status(400)

    def test_report_passes_and_failures(self) -> None:
        self.tester.send_request("GET", "/pass")
        try:
            self.tester.assert_status(200)
        except TestAssertionError:
            pass

        self.tester.send_request("POST", "/", body={"key": "x", "value": "y"})
        self.tester.assert_status(201)

        self.tester.send_request("GET", "/x")
        self.tester.assert_status(200)
        self.tester.assert_json("value", "y")

        report = self.tester.report
        self.assertGreaterEqual(report.passed, 2)
        self.assertGreaterEqual(report.total, 3)


if __name__ == "__main__":
    unittest.main()
