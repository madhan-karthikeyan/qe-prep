package networking

import (
	"net"
	"sync"
	"testing"
	"time"
)

func testFactory(addr string) func() (net.Conn, error) {
	return func() (net.Conn, error) {
		return net.Dial("tcp", addr)
	}
}

func TestConnPoolAcquireRelease(t *testing.T) {
	addr := "127.0.0.1:9999"
	pool := NewConnPool(2, time.Second, testFactory(addr))
	conn, err := pool.Acquire()
	if err == nil {
		pool.Release(conn)
	}
}

func TestConnPoolExhaustion(t *testing.T) {
	pool := NewConnPool(1, time.Second, func() (net.Conn, error) {
		return nil, net.UnknownNetworkError("test")
	})
	_, err := pool.Acquire()
	if err == nil {
		t.Fatal("expected error from factory")
	}
}

func TestConnPoolMaxConns(t *testing.T) {
	pool := NewConnPool(2, time.Second, func() (net.Conn, error) {
		return &mockConn{}, nil
	})
	c1, err := pool.Acquire()
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	c2, err := pool.Acquire()
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	_, err = pool.Acquire()
	if err != ErrPoolExhausted {
		t.Errorf("expected ErrPoolExhausted, got %v", err)
	}
	pool.Release(c1)
	pool.Release(c2)
}

func TestConnPoolReusesIdle(t *testing.T) {
	var factoryCalls int
	pool := NewConnPool(5, time.Second, func() (net.Conn, error) {
		factoryCalls++
		return &mockConn{}, nil
	})
	c1, _ := pool.Acquire()
	pool.Release(c1)
	c2, _ := pool.Acquire()
	if factoryCalls != 1 {
		t.Errorf("expected 1 factory call (reuse), got %d", factoryCalls)
	}
	_ = c2
}

func TestConnPoolStress(t *testing.T) {
	pool := NewConnPool(10, time.Second, func() (net.Conn, error) {
		return &mockConn{}, nil
	})
	var wg sync.WaitGroup
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := pool.Acquire()
			if err == nil {
				time.Sleep(time.Millisecond)
				pool.Release(conn)
			}
		}()
	}
	wg.Wait()
}

type mockConn struct {
	net.Conn
}

func (m *mockConn) Close() error                       { return nil }
func (m *mockConn) Read(b []byte) (int, error)         { return len(b), nil }
func (m *mockConn) Write(b []byte) (int, error)        { return len(b), nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }
