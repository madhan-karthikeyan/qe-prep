import threading
import unittest

from networking.implementation.connection_pool import ConnectionPool
from networking.implementation.tcp_echo import TCPEchoServer


class TestConnectionPool(unittest.TestCase):
    def setUp(self) -> None:
        self.server = TCPEchoServer()
        self.server.start()

    def tearDown(self) -> None:
        self.server.stop()

    def test_acquire_release(self) -> None:
        pool = ConnectionPool("127.0.0.1", self.server.port, max_connections=5)
        sock = pool.acquire()
        self.assertEqual(pool.active_connections, 1)
        self.assertEqual(pool.idle_connections, 0)
        pool.release(sock)
        self.assertEqual(pool.active_connections, 0)
        self.assertEqual(pool.idle_connections, 1)
        pool.close_all()

    def test_reuse_connection(self) -> None:
        pool = ConnectionPool("127.0.0.1", self.server.port, max_connections=5)
        sock1 = pool.acquire()
        pool.release(sock1)
        sock2 = pool.acquire()
        self.assertEqual(sock1, sock2)
        pool.close_all()

    def test_max_connections(self) -> None:
        pool = ConnectionPool("127.0.0.1", self.server.port, max_connections=3)
        socks = [pool.acquire() for _ in range(3)]
        with self.assertRaises(RuntimeError):
            pool.acquire()
        for s in socks:
            pool.release(s)
        pool.close_all()

    def test_invalid_construction(self) -> None:
        with self.assertRaises(ValueError):
            ConnectionPool("127.0.0.1", 0, max_connections=0)
        with self.assertRaises(ValueError):
            ConnectionPool("127.0.0.1", 0, timeout=0)
        with self.assertRaises(ValueError):
            ConnectionPool("127.0.0.1", 0, max_idle_seconds=0)

    def test_concurrent_access(self) -> None:
        pool = ConnectionPool("127.0.0.1", self.server.port, max_connections=20)
        results: list[int] = []
        lock = threading.Lock()

        def worker(n: int) -> None:
            conn = pool.acquire()
            with lock:
                results.append(n)
            pool.release(conn)

        threads = [threading.Thread(target=worker, args=(i,)) for i in range(20)]
        for t in threads:
            t.start()
        for t in threads:
            t.join()
        self.assertEqual(len(results), 20)
        pool.close_all()

    def test_close_all(self) -> None:
        pool = ConnectionPool("127.0.0.1", self.server.port, max_connections=5)
        _ = pool.acquire()
        pool.close_all()
        self.assertEqual(pool.active_connections, 0)
        self.assertEqual(pool.idle_connections, 0)
        with self.assertRaises(RuntimeError):
            pool.acquire()


if __name__ == "__main__":
    unittest.main()
