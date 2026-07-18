#include "tcp_echo.hpp"
#include <cassert>
#include <iostream>
#include <thread>
#include <chrono>

void test_tcp_echo() {
    int port = 9998;
    TcpServer server(port);

    std::thread server_thread([&server]() {
        server.start();
    });
    server_thread.detach();
    std::this_thread::sleep_for(std::chrono::milliseconds(100));

    TcpClient client("127.0.0.1", port);
    ssize_t sent = client.send("hello");
    assert(sent == 5);

    std::string reply = client.recv(64);
    assert(reply == "hello");

    std::cout << "PASS: tcp_echo" << std::endl;
}

int main() {
    test_tcp_echo();
    return 0;
}
