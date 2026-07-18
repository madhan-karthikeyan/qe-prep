import socket
import threading


class TCPEchoServer:
    def __init__(self, host: str = "127.0.0.1", port: int = 0) -> None:
        self._host = host
        self._port = port
        self._server_socket: socket.socket | None = None
        self._thread: threading.Thread | None = None
        self._running = False

    def start(self) -> None:
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        self._server_socket.bind((self._host, self._port))
        self._server_socket.listen(5)
        self._port = self._server_socket.getsockname()[1]
        self._running = True
        self._thread = threading.Thread(target=self._serve, daemon=True)
        self._thread.start()

    def _serve(self) -> None:
        assert self._server_socket is not None
        self._server_socket.settimeout(1.0)
        while self._running:
            try:
                client_sock, addr = self._server_socket.accept()
            except TimeoutError:
                continue
            except OSError:
                break
            t = threading.Thread(
                target=self._handle_client,
                args=(client_sock, addr),
                daemon=True,
            )
            t.start()

    def _handle_client(self, client_sock: socket.socket, addr: tuple[str, int]) -> None:
        with client_sock:
            try:
                while True:
                    data = client_sock.recv(4096)
                    if not data:
                        break
                    client_sock.sendall(data)
            except OSError:
                pass

    def stop(self) -> None:
        self._running = False
        if self._server_socket is not None:
            self._server_socket.close()
        if self._thread is not None:
            self._thread.join(timeout=2)

    @property
    def port(self) -> int:
        return self._port

    @property
    def host(self) -> str:
        return self._host


class TCPEchoClient:
    def __init__(self, host: str = "127.0.0.1", port: int = 0) -> None:
        self._host = host
        self._port = port

    def echo(self, data: bytes, timeout: float = 5.0) -> bytes:
        if not data:
            return b""
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(timeout)
        with sock:
            sock.connect((self._host, self._port))
            sock.sendall(data)
            chunks: list[bytes] = []
            while True:
                chunk = sock.recv(4096)
                if not chunk:
                    break
                chunks.append(chunk)
                if len(b"".join(chunks)) >= len(data):
                    break
            return b"".join(chunks)
