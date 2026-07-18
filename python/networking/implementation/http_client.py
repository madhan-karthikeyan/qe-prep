import http.client
import random
import time
import urllib.parse


class HTTPError(Exception):
    def __init__(self, status: int, body: str) -> None:
        self.status = status
        self.body = body
        super().__init__(f"HTTP {status}: {body[:200]}")


class HTTPClient:
    def __init__(
        self,
        max_retries: int = 3,
        base_delay: float = 0.5,
        max_delay: float = 10.0,
        jitter: bool = True,
        timeout: float = 10.0,
    ) -> None:
        if max_retries < 0:
            raise ValueError("max_retries must be non-negative")
        if base_delay <= 0:
            raise ValueError("base_delay must be positive")
        if max_delay < base_delay:
            raise ValueError("max_delay must be >= base_delay")
        if timeout <= 0:
            raise ValueError("timeout must be positive")

        self._max_retries = max_retries
        self._base_delay = base_delay
        self._max_delay = max_delay
        self._jitter = jitter
        self._timeout = timeout

    def get(
        self,
        url: str,
        headers: dict[str, str] | None = None,
    ) -> tuple[int, str]:
        return self._request("GET", url, headers=headers)

    def post(
        self,
        url: str,
        body: str | None = None,
        headers: dict[str, str] | None = None,
    ) -> tuple[int, str]:
        return self._request("POST", url, body=body, headers=headers)

    def _request(
        self,
        method: str,
        url: str,
        body: str | None = None,
        headers: dict[str, str] | None = None,
    ) -> tuple[int, str]:
        parsed = urllib.parse.urlparse(url)
        if parsed.scheme not in ("http", "https"):
            raise ValueError(f"Unsupported scheme: {parsed.scheme}")

        host = parsed.hostname or "localhost"
        port = parsed.port
        path = parsed.path or "/"
        if parsed.query:
            path = f"{path}?{parsed.query}"

        last_exc: Exception | None = None
        for attempt in range(self._max_retries + 1):
            try:
                if parsed.scheme == "https":
                    conn: http.client.HTTPConnection | http.client.HTTPSConnection = (
                        http.client.HTTPSConnection(host, port, timeout=self._timeout)
                    )
                else:
                    conn = http.client.HTTPConnection(host, port, timeout=self._timeout)

                try:
                    conn.request(method, path, body=body, headers=headers or {})
                    response = conn.getresponse()
                    status = response.status
                    data = response.read().decode("utf-8", errors="replace")
                finally:
                    conn.close()

                if 200 <= status < 300:
                    return status, data

                if status in (429,) or 500 <= status < 600:
                    last_exc = HTTPError(status, data)
                    if attempt < self._max_retries:
                        self._backoff(attempt)
                        continue

                return status, data

            except (http.client.HTTPException, OSError, TimeoutError) as e:
                last_exc = e
                if attempt < self._max_retries:
                    self._backoff(attempt)
                    continue
                raise

        raise HTTPError(0, f"Request failed after {self._max_retries} retries") from last_exc

    def _backoff(self, attempt: int) -> None:
        delay = min(self._base_delay * (2 ** attempt), self._max_delay)
        if self._jitter:
            delay = delay * (0.5 + random.random() * 0.5)
        time.sleep(delay)
