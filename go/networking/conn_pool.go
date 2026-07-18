package networking

import (
	"errors"
	"net"
	"sync"
	"time"
)

// ErrPoolExhausted is returned when the connection pool has reached its
// maximum number of active connections.
var ErrPoolExhausted = errors.New("connection pool exhausted")

// ConnPool manages a pool of reusable net.Conn connections.
type ConnPool struct {
	mu       sync.Mutex
	maxConns int
	timeout  time.Duration
	factory  func() (net.Conn, error)
	idle     []net.Conn
	active   int
}

// NewConnPool creates a new connection pool. factory is called to create new
// connections when no idle ones are available.
func NewConnPool(maxConns int, timeout time.Duration, factory func() (net.Conn, error)) *ConnPool {
	return &ConnPool{
		maxConns: maxConns,
		timeout:  timeout,
		factory:  factory,
		idle:     make([]net.Conn, 0, maxConns),
	}
}

// Acquire returns a connection from the pool or creates a new one. Returns
// ErrPoolExhausted if the pool is at capacity.
func (p *ConnPool) Acquire() (net.Conn, error) {
	p.mu.Lock()
	if len(p.idle) > 0 {
		conn := p.idle[len(p.idle)-1]
		p.idle = p.idle[:len(p.idle)-1]
		p.active++
		p.mu.Unlock()
		if p.timeout > 0 {
			_ = conn.SetDeadline(time.Now().Add(p.timeout))
		}
		return conn, nil
	}
	if p.active >= p.maxConns {
		p.mu.Unlock()
		return nil, ErrPoolExhausted
	}
	p.active++
	p.mu.Unlock()
	conn, err := p.factory()
	if err != nil {
		p.mu.Lock()
		p.active--
		p.mu.Unlock()
		return nil, err
	}
	if p.timeout > 0 {
		_ = conn.SetDeadline(time.Now().Add(p.timeout))
	}
	return conn, nil
}

// Release returns a connection to the pool for reuse.
func (p *ConnPool) Release(conn net.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.idle = append(p.idle, conn)
	p.active--
}

// Close closes all idle connections in the pool.
func (p *ConnPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, conn := range p.idle {
		conn.Close()
	}
	p.idle = nil
	p.active = 0
}

// Len returns the total number of connections (idle + active).
func (p *ConnPool) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.idle) + p.active
}
