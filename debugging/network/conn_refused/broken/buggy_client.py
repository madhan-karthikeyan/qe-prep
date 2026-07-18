import socket

def run_client():
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect(("127.0.0.1", 5000))
    sock.sendall(b"Hello, server!")
    data = sock.recv(1024)
    print(f"Received: {data.decode()}")
    sock.close()

if __name__ == "__main__":
    run_client()
