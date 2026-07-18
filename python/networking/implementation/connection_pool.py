import socket
import threading
import time
from collections import OrderedDict


class ConnectionPool:
    def __init__(
        self,
        host: str,
        port: int,
        max_connections: int = 10,
        timeout: float = 5.0,
        max_idle_seconds: float = 60.0,
    ) -> None:
        if max_connections <= 0:
            raise ValueError("max_connections must be positive")
        if timeout <= 0:
            raise ValueError("timeout must be positive")
        if max_idle_seconds <= 0:
            raise ValueError("max_idle_seconds must be positive")

        self._host = host
        self._port = port
        self._max_connections = max_connections
        self._timeout = timeout
        self._max_idle_seconds = max_idle_seconds

        self._available: OrderedDict[socket.socket, float] = OrderedDict()
        self._in_use: set[socket.socket] = set()
        self._lock = threading.Lock()
        self._closed = False

    def acquire(self) -> socket.socket:
        if self._closed:
            raise RuntimeError("Connection pool is closed")

        while True:
            sock = self._try_get_idle()
            if sock is not None:
                if self._is_healthy(sock):
                    with self._lock:
                        self._in_use.add(sock)
                    return sock
                self._close_socket(sock)
                continue

            if self._can_create():
                sock = self._create_connection()
                with self._lock:
                    self._in_use.add(sock)
                return sock

            raise RuntimeError("No available connections in pool")

    def _try_get_idle(self) -> socket.socket | None:
        with self._lock:
            self._evict_stale()
            if self._available:
                sock, _ = self._available.popitem(last=False)
                return sock
            return None

    def _can_create(self) -> bool:
        with self._lock:
            return len(self._in_use) + len(self._available) < self._max_connections

    def release(self, sock: socket.socket) -> None:
        with self._lock:
            self._in_use.discard(sock)
            if self._is_healthy(sock):
                self._available[sock] = time.monotonic()
            else:
                self._close_socket(sock)

    def _create_connection(self) -> socket.socket:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(self._timeout)
        try:
            sock.connect((self._host, self._port))
        except OSError:
            sock.close()
            raise
        return sock

    def _is_healthy(self, sock: socket.socket) -> bool:
        try:
            sock.settimeout(0.1)
            sock.recv(1, socket.MSG_PEEK)
            return True
        except (TimeoutError, BlockingIOError):
            return True
        except OSError:
            return False

    def _evict_stale(self) -> None:
        now = time.monotonic()
        stale = [
            sock
            for sock, last_used in self._available.items()
            if now - last_used > self._max_idle_seconds
        ]
        for sock in stale:
            self._available.pop(sock, None)
            self._close_socket(sock)

    def _close_socket(self, sock: socket.socket) -> None:
        try:
            sock.close()
        except OSError:
            pass

    @property
    def active_connections(self) -> int:
        with self._lock:
            return len(self._in_use)

    @property
    def idle_connections(self) -> int:
        with self._lock:
            return len(self._available)

    def close_all(self) -> None:
        with self._lock:
            self._closed = True
            for sock in list(self._available.keys()):
                self._close_socket(sock)
            self._available.clear()
            for sock in list(self._in_use):
                self._close_socket(sock)
            self._in_use.clear()
