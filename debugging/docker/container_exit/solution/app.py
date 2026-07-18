from http.server import HTTPServer, BaseHTTPRequestHandler
import time

class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b"Hello, I'm alive!\n")

    def log_message(self, format, *args):
        print(f"{self.client_address[0]} - {format % args}")

if __name__ == "__main__":
    server = HTTPServer(("0.0.0.0", 8080), Handler)
    print("Server started on port 8080")
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        server.shutdown()
