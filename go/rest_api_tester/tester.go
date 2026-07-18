package rest_api_tester

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TestResult holds the outcome of a single test case.
type TestResult struct {
	Name   string
	Passed bool
	Error  string
}

// TestCase defines an API test case to run.
type TestCase struct {
	Name    string
	Method  string
	Path    string
	Headers map[string]string
	Body    interface{}
	Assert  func(*http.Response) error
}

// APITester provides utilities for testing REST APIs with retries and
// assertions.
type APITester struct {
	baseURL    string
	client     *http.Client
	maxRetries int
	baseDelay  time.Duration
}

// NewAPITester creates a new APITester targeting the given base URL.
func NewAPITester(baseURL string) *APITester {
	return &APITester{
		baseURL:    baseURL,
		client:     &http.Client{Timeout: 5 * time.Second},
		maxRetries: 3,
		baseDelay:  50 * time.Millisecond,
	}
}

// SendRequest sends an HTTP request and returns the response.
func (t *APITester) SendRequest(method, path string, headers map[string]string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, t.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return t.client.Do(req)
}

// AssertStatus checks that the response has the expected status code.
func (t *APITester) AssertStatus(expected int, resp *http.Response) error {
	if resp.StatusCode != expected {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return fmt.Errorf("expected status %d, got %d: body=%s", expected, resp.StatusCode, string(body))
	}
	return nil
}

// AssertJSON checks that a JSON path in the response body matches the expected
// value. Path uses dot notation for nested keys (e.g., "data.name").
func (t *APITester) AssertJSON(path string, expected interface{}, resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}
	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	actual := resolveJSONPath(data, path)
	gotStr := fmt.Sprintf("%v", actual)
	wantStr := fmt.Sprintf("%v", expected)
	if gotStr != wantStr {
		return fmt.Errorf("path %q: expected %v, got %v", path, expected, actual)
	}
	return nil
}

// RunTests executes a slice of test cases and returns results.
func (t *APITester) RunTests(cases []TestCase) []TestResult {
	results := make([]TestResult, 0, len(cases))
	for _, tc := range cases {
		result := TestResult{Name: tc.Name}
		resp, err := t.sendWithRetry(tc.Method, tc.Path, tc.Headers, tc.Body)
		if err != nil {
			result.Error = fmt.Sprintf("request failed: %v", err)
			results = append(results, result)
			continue
		}
		if tc.Assert != nil {
			if err := tc.Assert(resp); err != nil {
				result.Error = err.Error()
				results = append(results, result)
				continue
			}
		}
		result.Passed = true
		results = append(results, result)
	}
	return results
}

func (t *APITester) sendWithRetry(method, path string, headers map[string]string, body interface{}) (*http.Response, error) {
	var resp *http.Response
	var err error
	for attempt := 0; attempt <= t.maxRetries; attempt++ {
		resp, err = t.SendRequest(method, path, headers, body)
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		if attempt < t.maxRetries {
			time.Sleep(t.backoff(attempt))
		}
	}
	return resp, err
}

func (t *APITester) backoff(attempt int) time.Duration {
	delay := float64(t.baseDelay) * math.Pow(2, float64(attempt))
	jitter := rand.Float64() * delay * 0.3
	return time.Duration(delay + jitter)
}

func resolveJSONPath(data interface{}, path string) interface{} {
	if path == "" {
		return data
	}
	parts := strings.Split(path, ".")
	current := data
	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			current = v[part]
		case []interface{}:
			idx, err := strconv.Atoi(part)
			if err != nil || idx < 0 || idx >= len(v) {
				return nil
			}
			current = v[idx]
		default:
			return nil
		}
	}
	return current
}
