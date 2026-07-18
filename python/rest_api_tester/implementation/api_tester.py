import json
import time
import urllib.error
import urllib.request
from dataclasses import dataclass, field
from typing import Any


class TestAssertionError(Exception):
    __test__ = False
    pass


@dataclass
class TestReport:
    total: int = 0
    passed: int = 0
    failed: int = 0
    errors: list[str] = field(default_factory=list)
    timing: float = 0.0


class APITester:
    def __init__(self, base_url: str, max_retries: int = 0, retry_delay: float = 0.5) -> None:
        if not base_url:
            raise ValueError("base_url must not be empty")
        self._base_url = base_url.rstrip("/")
        self._max_retries = max_retries
        self._retry_delay = retry_delay
        self._report = TestReport()
        self._last_status: int | None = None
        self._last_body: Any = None

    def send_request(
        self,
        method: str,
        path: str,
        headers: dict[str, str] | None = None,
        body: Any = None,
    ) -> "APITester":
        url = f"{self._base_url}{path}"
        data = json.dumps(body).encode("utf-8") if body is not None else None
        req = urllib.request.Request(url, data=data, method=method)
        if headers:
            for k, v in headers.items():
                req.add_header(k, v)
        if data is not None:
            req.add_header("Content-Type", "application/json")

        last_error: Exception | None = None
        for attempt in range(self._max_retries + 1):
            try:
                start = time.monotonic()
                with urllib.request.urlopen(req, timeout=10) as resp:
                    body_bytes = resp.read()
                elapsed = time.monotonic() - start
                self._report.timing += elapsed
                self._last_status = resp.status
                self._last_body = json.loads(body_bytes.decode("utf-8"))
                self._report.total += 1
                return self
            except (urllib.error.HTTPError, OSError, TimeoutError) as e:
                last_error = e
                if attempt < self._max_retries:
                    time.sleep(self._retry_delay)
                    continue
                if isinstance(e, urllib.error.HTTPError):
                    self._last_status = e.code
                    self._last_body = {"error": str(e)}
                    self._report.total += 1
                    return self
                raise

        raise RuntimeError(f"Request failed after {self._max_retries} retries") from last_error

    def assert_status(self, expected: int) -> "APITester":
        if self._last_status != expected:
            msg = (
                f"Expected status {expected}, got {self._last_status}. "
                f"Body: {self._last_body}"
            )
            self._report.failed += 1
            self._report.errors.append(msg)
            raise TestAssertionError(msg)
        self._report.passed += 1
        return self

    def assert_json(self, path: str, expected_value: Any) -> "APITester":
        parts = path.split(".")
        value: Any = self._last_body
        for part in parts:
            if isinstance(value, dict) and part in value:
                value = value[part]
            elif isinstance(value, list) and part.isdigit():
                idx = int(part)
                if 0 <= idx < len(value):
                    value = value[idx]
                else:
                    msg = (
                        f"Path {path!r}: index {idx} out of range. "
                        f"Body: {self._last_body}"
                    )
                    self._report.failed += 1
                    self._report.errors.append(msg)
                    raise TestAssertionError(msg)
            else:
                msg = (
                    f"Path {path!r}: key {part!r} not found. "
                    f"Body: {self._last_body}"
                )
                self._report.failed += 1
                self._report.errors.append(msg)
                raise TestAssertionError(msg)

        if value != expected_value:
            msg = (
                f"Path {path!r}: expected {expected_value!r}, got {value!r}. "
                f"Body: {self._last_body}"
            )
            self._report.failed += 1
            self._report.errors.append(msg)
            raise TestAssertionError(msg)

        self._report.passed += 1
        return self

    @property
    def report(self) -> TestReport:
        return self._report

    @property
    def last_status(self) -> int | None:
        return self._last_status

    @property
    def last_body(self) -> Any:
        return self._last_body
