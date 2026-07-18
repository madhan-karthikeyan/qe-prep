#ifndef TCP_ECHO_HPP
#define TCP_ECHO_HPP

#include <string>
#include <thread>
#include <vector>

class TcpServer {
public:
    TcpServer(int port);
    ~TcpServer();
    void start();

private:
    int port_;
    int server_fd_;
    std::vector<std::thread> clients_;
    void handle_client(int client_fd);
};

class TcpClient {
public:
    TcpClient(const std::string &host, int port);
    ~TcpClient();
    ssize_t send(const std::string &data);
    std::string recv(size_t len);

private:
    int fd_;
};

#endif
