import unittest

from networking.implementation.tcp_echo import TCPEchoClient, TCPEchoServer


class TestTCPEcho(unittest.TestCase):
    def setUp(self) -> None:
        self.server = TCPEchoServer()
        self.server.start()

    def tearDown(self) -> None:
        self.server.stop()

    def test_echo_basic(self) -> None:
        client = TCPEchoClient(port=self.server.port)
        result = client.echo(b"hello")
        self.assertEqual(result, b"hello")

    def test_echo_multiple_messages(self) -> None:
        client = TCPEchoClient(port=self.server.port)
        self.assertEqual(client.echo(b"msg1"), b"msg1")
        self.assertEqual(client.echo(b"msg2"), b"msg2")
        self.assertEqual(client.echo(b"msg3"), b"msg3")

    def test_echo_binary_data(self) -> None:
        client = TCPEchoClient(port=self.server.port)
        data = bytes(range(256))
        result = client.echo(data)
        self.assertEqual(result, data)

    def test_echo_large_payload(self) -> None:
        client = TCPEchoClient(port=self.server.port)
        data = b"x" * 65536
        result = client.echo(data)
        self.assertEqual(result, data)

    def test_echo_empty(self) -> None:
        client = TCPEchoClient(port=self.server.port)
        result = client.echo(b"")
        self.assertEqual(result, b"")

    def test_multiple_clients(self) -> None:
        results: list[bytes] = []
        for i in range(10):
            client = TCPEchoClient(port=self.server.port)
            results.append(client.echo(f"msg{i}".encode()))
        for i, r in enumerate(results):
            self.assertEqual(r, f"msg{i}".encode())


if __name__ == "__main__":
    unittest.main()
