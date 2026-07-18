package networking

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
)

// StartServer starts a TCP echo server on the given address. It accepts
// concurrent connections, each handled in its own goroutine. The server
// shuts down gracefully when ctx is cancelled.
func StartServer(ctx context.Context, addr string) error {
	lc := net.ListenConfig{}
	listener, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}
	go func() {
		<-ctx.Done()
		listener.Close()
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return fmt.Errorf("accept: %w", err)
		}
		go handleEcho(conn)
	}
}

func handleEcho(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			return
		}
		if _, err := conn.Write(buf[:n]); err != nil {
			return
		}
	}
}

// EchoClient connects to a TCP echo server at addr, sends msg, and returns
// the echoed response.
func EchoClient(addr, msg string) (string, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", fmt.Errorf("dial %s: %w", addr, err)
	}
	defer conn.Close()
	if _, err := fmt.Fprint(conn, msg); err != nil {
		return "", fmt.Errorf("send: %w", err)
	}
	tcpConn, ok := conn.(*net.TCPConn)
	if ok {
		if err := tcpConn.CloseWrite(); err != nil {
			return "", fmt.Errorf("close write: %w", err)
		}
	}
	resp, err := io.ReadAll(conn)
	if err != nil {
		return "", fmt.Errorf("read: %w", err)
	}
	return string(resp), nil
}
