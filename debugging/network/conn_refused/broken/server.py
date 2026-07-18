import socket

def run_server():
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    sock.bind(("127.0.0.1", 9999))
    sock.listen(5)
    print("Echo server listening on port 9999...")
    while True:
        conn, addr = sock.accept()
        print(f"Client connected: {addr}")
        data = conn.recv(1024)
        if data:
            conn.sendall(data)
        conn.close()

if __name__ == "__main__":
    run_server()
