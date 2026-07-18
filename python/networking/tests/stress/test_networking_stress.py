import threading
import time
import unittest

from networking.implementation.connection_pool import ConnectionPool
from networking.implementation.tcp_echo import TCPEchoServer


class TestNetworkingStress(unittest.TestCase):
    def setUp(self) -> None:
        self.server = TCPEchoServer()
        self.server.start()

    def tearDown(self) -> None:
        self.server.stop()

    def test_50_concurrent_connections_through_pool(self) -> None:
        pool = ConnectionPool(
            "127.0.0.1",
            self.server.port,
            max_connections=50,
            timeout=10.0,
        )
        errors: list[Exception] = []
        lock = threading.Lock()

        def worker(n: int) -> None:
            try:
                for _ in range(10):
                    sock = pool.acquire()
                    msg = f"stress_{n}".encode()
                    sock.sendall(msg)
                    data = sock.recv(4096)
                    assert data == msg, f"Data mismatch: {data!r} != {msg!r}"
                    pool.release(sock)
                    time.sleep(0.001)
            except Exception as e:
                with lock:
                    errors.append(e)

        threads = [threading.Thread(target=worker, args=(i,)) for i in range(50)]
        start = time.monotonic()
        for t in threads:
            t.start()
        for t in threads:
            t.join()
        elapsed = time.monotonic() - start

        pool.close_all()
        self.assertEqual(len(errors), 0, f"Errors: {errors}")
        self.assertLess(elapsed, 60)


if __name__ == "__main__":
    unittest.main()
