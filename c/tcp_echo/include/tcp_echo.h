#ifndef TCP_ECHO_H
#define TCP_ECHO_H

#include <sys/types.h>

int tcp_server_start(int port);
int tcp_client_connect(const char *host, int port);
ssize_t tcp_send(int fd, const void *buf, size_t len);
ssize_t tcp_recv(int fd, void *buf, size_t len);
void tcp_close(int fd);

#endif
