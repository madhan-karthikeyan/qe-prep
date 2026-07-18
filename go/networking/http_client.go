package networking

import (
	"io"
	"math"
	"math/rand"
	"net/http"
	"time"
)

// RetryClient wraps net/http.Client with configurable retries, exponential
// backoff, and jitter.
type RetryClient struct {
	client     *http.Client
	maxRetries int
	baseDelay  time.Duration
	maxDelay   time.Duration
}

// NewRetryClient creates a new RetryClient.
func NewRetryClient(timeout time.Duration, maxRetries int, baseDelay, maxDelay time.Duration) *RetryClient {
	return &RetryClient{
		client: &http.Client{
			Timeout: timeout,
		},
		maxRetries: maxRetries,
		baseDelay:  baseDelay,
		maxDelay:   maxDelay,
	}
}

// Get issues a GET request with retry logic.
func (c *RetryClient) Get(url string) (*http.Response, error) {
	return c.doWithRetry(func() (*http.Response, error) {
		return c.client.Get(url)
	})
}

// Post issues a POST request with retry logic. The body is consumed on each
// attempt; for retries to work the caller must provide a fresh body or use a
// rewindable reader.
func (c *RetryClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return c.doWithRetry(func() (*http.Response, error) {
		return c.client.Post(url, contentType, body)
	})
}

func (c *RetryClient) doWithRetry(do func() (*http.Response, error)) (*http.Response, error) {
	var resp *http.Response
	var err error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		resp, err = do()
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		if attempt < c.maxRetries {
			time.Sleep(c.backoff(attempt))
		}
	}
	return resp, err
}

// backoff computes exponential backoff with jitter.
func (c *RetryClient) backoff(attempt int) time.Duration {
	delay := float64(c.baseDelay) * math.Pow(2, float64(attempt))
	if delay > float64(c.maxDelay) {
		delay = float64(c.maxDelay)
	}
	jitter := rand.Float64() * delay * 0.3
	return time.Duration(delay + jitter)
}
