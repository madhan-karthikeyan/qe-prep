#include "tcp_echo.h"
#include <assert.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <sys/wait.h>
#include <sys/socket.h>

void test_tcp_echo() {
    int port = 9999;
    int server_fd = tcp_server_start(port);
    assert(server_fd >= 0);

    pid_t pid = fork();
    if (pid == 0) {
        int client = tcp_client_connect("127.0.0.1", port);
        assert(client >= 0);

        const char *msg = "hello";
        ssize_t sent = tcp_send(client, msg, strlen(msg));
        assert(sent == (ssize_t)strlen(msg));

        char buf[64] = {0};
        ssize_t rcvd = tcp_recv(client, buf, sizeof(buf) - 1);
        assert(rcvd == (ssize_t)strlen(msg));
        assert(strcmp(buf, msg) == 0);

        tcp_close(client);
        _exit(0);
    } else {
        int client = accept(server_fd, NULL, NULL);
        assert(client >= 0);

        char buf[64] = {0};
        ssize_t n = tcp_recv(client, buf, sizeof(buf) - 1);
        assert(n > 0);

        tcp_send(client, buf, n);
        tcp_close(client);

        int status;
        waitpid(pid, &status, 0);
    }

    tcp_close(server_fd);
    printf("PASS: tcp_echo\n");
}

int main() {
    test_tcp_echo();
    return 0;
}
