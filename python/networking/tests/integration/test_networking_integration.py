import threading
import unittest

from networking.implementation.connection_pool import ConnectionPool
from networking.implementation.tcp_echo import TCPEchoClient, TCPEchoServer


class TestNetworkingIntegration(unittest.TestCase):
    def setUp(self) -> None:
        self.server = TCPEchoServer()
        self.server.start()

    def tearDown(self) -> None:
        self.server.stop()

    def test_echo_round_trip(self) -> None:
        client = TCPEchoClient(port=self.server.port)
        payload = b"Hello, World!"
        response = client.echo(payload)
        self.assertEqual(response, payload)

    def test_echo_multiple_round_trips(self) -> None:
        client = TCPEchoClient(port=self.server.port)
        for msg in (b"alpha", b"beta", b"gamma"):
            self.assertEqual(client.echo(msg), msg)

    def test_pool_echo_round_trip(self) -> None:
        pool = ConnectionPool("127.0.0.1", self.server.port, max_connections=5)
        sock = pool.acquire()
        sock.sendall(b"pool test")
        data = sock.recv(4096)
        self.assertEqual(data, b"pool test")
        pool.release(sock)
        pool.close_all()

    def test_concurrent_echo_via_pool(self) -> None:
        pool = ConnectionPool("127.0.0.1", self.server.port, max_connections=10)
        results: list[bool] = []
        lock = threading.Lock()

        def worker(msg: bytes) -> None:
            sock = pool.acquire()
            sock.sendall(msg)
            data = sock.recv(4096)
            with lock:
                results.append(data == msg)
            pool.release(sock)

        messages = [f"msg{i}".encode() for i in range(10)]
        threads = [threading.Thread(target=worker, args=(m,)) for m in messages]
        for t in threads:
            t.start()
        for t in threads:
            t.join()
        pool.close_all()
        self.assertTrue(all(results))


if __name__ == "__main__":
    unittest.main()
