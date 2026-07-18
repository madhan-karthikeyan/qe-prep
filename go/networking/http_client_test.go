package networking

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRetryClientGet(t *testing.T) {
	var attempts int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer ts.Close()

	client := NewRetryClient(time.Second, 2, 10*time.Millisecond, 100*time.Millisecond)
	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", attempts)
	}
}

func TestRetryClientRetriesOn5xx(t *testing.T) {
	var attempts int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := NewRetryClient(time.Second, 2, 10*time.Millisecond, 100*time.Millisecond)
	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	defer resp.Body.Close()
	if attempts != 3 {
		t.Errorf("expected 3 attempts (initial + 2 retries), got %d", attempts)
	}
}

func TestRetryClientSucceedsOnRetry(t *testing.T) {
	var attempts int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer ts.Close()

	client := NewRetryClient(time.Second, 3, 10*time.Millisecond, 100*time.Millisecond)
	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRetryClientZeroRetries(t *testing.T) {
	var attempts int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := NewRetryClient(time.Second, 0, 10*time.Millisecond, 100*time.Millisecond)
	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	defer resp.Body.Close()
	if attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", attempts)
	}
}
