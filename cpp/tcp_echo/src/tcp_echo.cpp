#include "tcp_echo.hpp"
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <netdb.h>
#include <unistd.h>
#include <cstring>
#include <iostream>
#include <algorithm>

TcpServer::TcpServer(int port) : port_(port), server_fd_(-1) {}

TcpServer::~TcpServer() {
    if (server_fd_ >= 0) ::close(server_fd_);
    for (auto &t : clients_)
        if (t.joinable()) t.join();
}

void TcpServer::handle_client(int client_fd) {
    char buf[1024];
    ssize_t n;
    while ((n = recv(client_fd, buf, sizeof(buf), 0)) > 0) {
        send(client_fd, buf, n, 0);
    }
    ::close(client_fd);
}

void TcpServer::start() {
    server_fd_ = socket(AF_INET, SOCK_STREAM, 0);
    int opt = 1;
    setsockopt(server_fd_, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt));

    struct sockaddr_in addr;
    std::memset(&addr, 0, sizeof(addr));
    addr.sin_family = AF_INET;
    addr.sin_addr.s_addr = INADDR_ANY;
    addr.sin_port = htons(port_);

    bind(server_fd_, (struct sockaddr *)&addr, sizeof(addr));
    listen(server_fd_, 5);

    while (true) {
        int client_fd = accept(server_fd_, nullptr, nullptr);
        if (client_fd < 0) break;
        clients_.emplace_back(&TcpServer::handle_client, this, client_fd);
        clients_.back().detach();
    }
}

TcpClient::TcpClient(const std::string &host, int port) {
    fd_ = socket(AF_INET, SOCK_STREAM, 0);
    struct hostent *he = gethostbyname(host.c_str());
    struct sockaddr_in addr;
    std::memset(&addr, 0, sizeof(addr));
    addr.sin_family = AF_INET;
    addr.sin_port = htons(port);
    std::memcpy(&addr.sin_addr, he->h_addr_list[0], he->h_length);
    connect(fd_, (struct sockaddr *)&addr, sizeof(addr));
}

TcpClient::~TcpClient() {
    if (fd_ >= 0) ::close(fd_);
}

ssize_t TcpClient::send(const std::string &data) {
    return ::send(fd_, data.data(), data.size(), 0);
}

std::string TcpClient::recv(size_t len) {
    char buf[1024];
    ssize_t n = ::recv(fd_, buf, std::min(len, sizeof(buf)), 0);
    if (n > 0) return std::string(buf, n);
    return "";
}
