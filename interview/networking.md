# Networking Interview Guide — QE Engineer

## Overview

Networking knowledge is essential for QE engineers debugging distributed systems, configuring test environments, and understanding latency/throughput bottlenecks. Expect questions on TCP, HTTP, DNS, TLS, and load balancing.

## Top 20 Networking Questions

### TCP/IP Basics

1. **Describe the TCP 3-way handshake.**
   - **Expected**: 
     1. Client sends SYN (seq=x)
     2. Server responds SYN-ACK (seq=y, ack=x+1)
     3. Client sends ACK (seq=x+1, ack=y+1)
   - After this, connection is established.
   - **Mistake**: Forgetting ACK number is SYN number + 1.

2. **Why does TCP need a 3-way handshake? Why not 2-way?**
   - **Expected**: Prevents stale SYN packets from previous connections creating phantom connections (a delayed SYN gets ACKed by a listening server). The third ACK confirms the client received the server's SYN.

3. **Explain TCP connection termination (4-way handshake).**
   - **Expected**: FIN from client → ACK from server → FIN from server → ACK from client. Or half-close: one side sends FIN, still receives data.

4. **What is the TIME_WAIT state? Why does it exist?**
   - **Expected**: After sending final ACK, the closing side waits 2×MSL (typically 60s). Ensures delayed packets don't corrupt new connections. Testing issue: socket exhaustion under high connection rates.

5. **What is the difference between TCP and UDP? When would you use each?**
   - **TCP**: Reliable, ordered, connection-oriented. Use: HTTP, databases, file transfers.
   - **UDP**: Unreliable, unordered, connectionless. Use: streaming, DNS, real-time comms, metrics.

### HTTP (Difficulty: ★☆☆ Easy–★★☆ Medium)

6. **What happens when you type "https://example.com" in a browser?**
   - **Expected**: DNS lookup → TCP handshake → TLS handshake → HTTP request → response → render. Hit cache along the way.
   - **Testing angle**: Measure each phase (DNS resolution time, TCP connect time, TLS time, TTFB).

7. **Compare HTTP/1.1 and HTTP/2.**
   | Feature | HTTP/1.1 | HTTP/2 |
   |---------|----------|--------|
   | Multiplexing | No (pipeline limited) | Yes (multiple streams over 1 connection) |
   | Head-of-line blocking | Yes (request blocked until response) | Application-level only (not fixed for TCP HOL) |
   | Header compression | No | Yes (HPACK) |
   | Server push | No | Yes |
   | Binary protocol | No (text) | Yes (binary framing) |

8. **What is HTTP/3? How does it differ from HTTP/2?**
   - **Expected**: HTTP/3 uses QUIC (UDP-based, not TCP). Eliminates TCP head-of-line blocking. Faster connection establishment (0-RTT).

9. **Explain the difference between `Keep-Alive` and connection pooling.**
   - **Keep-Alive**: HTTP header reusing the same TCP connection for multiple requests.
   - **Connection pooling**: Client-side pool of pre-established connections (e.g., HikariCP for databases). Reduces connection overhead.

10. **What status codes do you test for?**
    - **2xx**: Success (200 OK, 201 Created, 204 No Content)
    - **3xx**: Redirect (301 Moved, 304 Not Modified)
    - **4xx**: Client error (400 Bad Request, 401 Unauthorized, 403 Forbidden, 404 Not Found, 429 Rate Limited)
    - **5xx**: Server error (500 Internal Server Error, 502 Bad Gateway, 503 Service Unavailable, 504 Gateway Timeout)
    - **Test:** Verify both expected success codes AND proper error codes.

### DNS (Difficulty: ★☆☆ Easy)

11. **How does DNS resolution work?**
    - **Expected**: Browser cache → OS cache → hosts file → recursive resolver → root NS → TLD NS → authoritative NS → IP address.
    - Time scales: caching TTL controls propagation delay.

12. **What is DNS round-robin? How does it affect testing?**
    - **Expected**: Multiple A records for one hostname; DNS returns them in rotation. Test con: caching means clients may use same IP. Solution: use DNS-based load testing with random resolution.

13. **What is DNS spoofing/cache poisoning?**
    - **Expected**: Attacker injects fake DNS records into a resolver's cache. Mitigations: DNSSEC (signatures), random source ports, query ID randomization.

### Load Balancing (Difficulty: ★★☆ Medium)

14. **Compare load balancing algorithms.**
    | Algorithm | How It Works | Best For |
    |-----------|-------------|----------|
    | Round-robin | Distributes sequentially | Uniform request times |
    | Least connections | Sends to least busy server | Variable-length requests |
    | IP hash | Hash of client IP → server | Session persistence |
    | Weighted | Proportional to server capacity | Heterogeneous servers |
    | Random | Picks randomly | Simple, uniform load |

15. **What is the difference between L4 and L7 load balancing?**
    - **L4 (transport layer)**: Routes TCP/UDP connections by IP + port. Faster, no content awareness.
    - **L7 (application layer)**: Inspects HTTP headers, cookies, paths. Slower, but supports intelligent routing, SSL termination, content-based routing.

16. **What is a health check in a load balancer?**
    - **Expected**: Periodic probes to backend servers (TCP connect, HTTP 200 check, or custom script). Unhealthy servers are removed from the pool. Test: verify health check timeout, degraded response detection, graceful drain.

### TLS (Difficulty: ★★☆ Medium)

17. **Describe the TLS 1.3 handshake (simplified).**
    - **Expected**:
      1. Client sends ClientHello (supported ciphers, key share)
      2. Server responds ServerHello (chosen cipher, key share, cert), finished
      3. Client sends finished
    - **Key**: TLS 1.3 handshake is 1-RTT (vs 2-RTT for TLS 1.2). 0-RTT mode allows sending data with first message (requires pre-shared key).

18. **What is SNI (Server Name Indication)?**
    - **Expected**: Extension in TLS that allows the client to specify the hostname it's connecting to. Enables multiple HTTPS sites on a single IP. Without SNI, only one cert per IP.

19. **How would you test TLS misconfigurations?**
    - **Expected**: Use `sslscan` or `testssl.sh`. Check: weak ciphers (RC4, DES), expired certs, self-signed certs, mismatched hostname, incomplete chain, TLS version support (TLS <1.2 should be disabled).

### Common Networking Interview Problems

20. **A service is returning intermittent 502 Bad Gateway. How do you debug?**
    - **Expected**:
      1. Check load balancer health checks — is backend marked unhealthy?
      2. Check backend logs — errors, OOM, connection pool exhaustion?
      3. Check backend response time — is it exceeding timeout?
      4. Check upstream — is the upstream service (DB, API) slow?
      5. Check network — packet loss? TCP retransmissions? (`tcpdump`, `netstat -s`)
      6. Reproduce — curl with specific headers, timeout flags

21. **What tools do you use for network troubleshooting?**
    | Tool | Use |
    |------|-----|
    | `tcpdump` / Wireshark | Packet capture, analyze handshake, retransmission |
    | `curl -v` / `wget` | HTTP debugging, headers, timing |
    | `ping` / `mtr` | Latency, packet loss, routing |
    | `nslookup` / `dig` | DNS resolution debugging |
    | `ss` / `netstat` / `lsof` | Connection state, open ports, listening services |
    | `iperf` / `netcat` | Bandwidth testing, raw data transfer |
    | `openssl s_client` | TLS debugging, cert inspection |
    | `traceroute` | Path discovery, latency per hop |

22. **What causes TCP retransmissions? How do you detect them?**
    - **Expected**: Packet loss, network congestion, timeouts, out-of-order delivery. Detect via `ss -i`, `netstat -s`, `tcpdump`, or Wireshark's "TCP Retransmission" filter.

23. **What is the Nagle algorithm? When would you disable it?**
    - **Expected**: Coalesces small packets to reduce overhead. Waits for ACK before sending next. Disable when low latency is critical (game servers, real-time apps, SSH). Enable with `TCP_NODELAY` socket option.

24. **Explain BDP (Bandwidth-Delay Product). Why does it matter for TCP?**
    - **Expected**: BDP = bandwidth × RTT. TCP window size should be ≥ BDP to fully utilize the link. On high-latency links, default window is too small — need window scaling (RFC 1323) or alternative congestion control (BBR).

25. **What is a reverse proxy vs forward proxy?**
    - **Forward proxy**: Client-side, hides client from server (corporate proxy, VPN).
    - **Reverse proxy**: Server-side, hides server from client (load balancer, CDN, API gateway).

---

## How to Approach Networking Questions

- **Use layers**: Physical → Data Link → Network → Transport → Application. Frame answers at the relevant layer.
- **Show tool knowledge**: Mention `tcpdump`, `curl -v`, `dig` — interviewers want hands-on experience.
- **Connect to testing**: "I would verify this by capturing traffic with tcpdump during the test and looking for..."

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Handshake details wrong | Memorize exact SYN/SYN-ACK/ACK flow |
| Ignoring TIME_WAIT impact | Mention connection reuse and socket exhaustion |
| HTTP/2 = SPDY | HTTP/2 is based on Google's SPDY, but they differ |
| Forgetting DNS caching | Tests can give misleading results due to cache — use `dig` with +trace |

## Difficulty Levels

| Topic | Difficulty |
|-------|-----------|
| TCP handshake | ★☆☆ |
| DNS basics | ★☆☆ |
| HTTP/1.1 vs HTTP/2 | ★★☆ |
| Load balancing | ★★☆ |
| TLS handshake | ★★☆ |
| TCP retransmission analysis | ★★★ |
| BDP / congestion control | ★★★ |
