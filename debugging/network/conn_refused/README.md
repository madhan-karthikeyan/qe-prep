## Symptoms
Client crashes with `ConnectionRefusedError: [Errno 111] Connection refused` when trying to connect.

## Root Cause
The server listens on port 9999, but the client connects to port 5000 — a port where nothing is listening. TCP's SYN packet receives an RST because no socket is bound to that port.

## Fix
Change the client's connect port from 5000 to 9999 to match the server.

## Prevention
- Centralize port configuration in environment variables or a config file.
- Add startup logging to servers (e.g., "listening on port X").
- Use connection timeouts to distinguish "refused" from "unreachable" scenarios.
