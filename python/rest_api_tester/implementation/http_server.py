import json
import threading
from http.server import BaseHTTPRequestHandler, HTTPServer
from typing import Any


class _RequestHandler(BaseHTTPRequestHandler):
    store: dict[str, Any] = {}
    configurable_status: int = 200

    def do_GET(self) -> None:
        self._handle_request("GET")

    def do_POST(self) -> None:
        self._handle_request("POST")

    def do_DELETE(self) -> None:
        self._handle_request("DELETE")

    def _handle_request(self, method: str) -> None:
        content_length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(content_length).decode("utf-8") if content_length else ""

        if method == "GET":
            key = self.path.lstrip("/")
            data = _RequestHandler.store.get(key)
            if data is not None:
                self._send_json(200, {"key": key, "value": data})
            else:
                self._send_json(404, {"error": "Not found"})
        elif method == "POST":
            try:
                payload = json.loads(body) if body else {}
            except json.JSONDecodeError:
                self._send_json(400, {"error": "Invalid JSON"})
                return
            if isinstance(payload, dict) and "key" in payload and "value" in payload:
                _RequestHandler.store[payload["key"]] = payload["value"]
                self._send_json(201, payload)
            else:
                self._send_json(400, {"error": "Missing key/value"})
        elif method == "DELETE":
            key = self.path.lstrip("/")
            existed = key in _RequestHandler.store
            _RequestHandler.store.pop(key, None)
            if existed:
                self._send_json(200, {"deleted": key})
            else:
                self._send_json(404, {"error": "Not found"})

    def _send_json(self, status: int, data: dict[str, Any]) -> None:
        status = _RequestHandler.configurable_status if status == 200 else status
        body = json.dumps(data).encode("utf-8")
        self.send_response(status)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def log_message(self, format: str, *args: Any) -> None:
        pass


class TestHTTPServer:
    __test__ = False
    def __init__(self, host: str = "127.0.0.1", port: int = 0) -> None:
        self._host = host
        self._port = port
        self._server: HTTPServer | None = None
        self._thread: threading.Thread | None = None

    def start(self) -> None:
        _RequestHandler.store.clear()
        _RequestHandler.configurable_status = 200
        self._server = HTTPServer((self._host, self._port), _RequestHandler)
        self._port = self._server.server_address[1]
        self._thread = threading.Thread(target=self._server.serve_forever, daemon=True)
        self._thread.start()

    def stop(self) -> None:
        if self._server is not None:
            self._server.shutdown()
            self._server.server_close()

    @property
    def port(self) -> int:
        return self._port

    @property
    def host(self) -> str:
        return self._host

    @staticmethod
    def set_default_status(status: int) -> None:
        _RequestHandler.configurable_status = status

    @staticmethod
    def clear_store() -> None:
        _RequestHandler.store.clear()
