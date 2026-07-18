package rest_api_tester

import (
	"net/http"
	"testing"
)

func TestTesterSendRequest(t *testing.T) {
	s := NewTestServer()
	if err := s.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer s.Close()

	tester := NewAPITester(s.URL())

	resp, err := tester.SendRequest("POST", "/items", nil, map[string]interface{}{
		"id":   "test1",
		"name": "item1",
	})
	if err != nil {
		t.Fatalf("SendRequest: %v", err)
	}
	if err := tester.AssertStatus(http.StatusCreated, resp); err != nil {
		t.Errorf("AssertStatus: %v", err)
	}
}

func TestTesterAssertJSON(t *testing.T) {
	s := NewTestServer()
	if err := s.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer s.Close()

	tester := NewAPITester(s.URL())
	resp, _ := tester.SendRequest("POST", "/items", nil, map[string]interface{}{
		"id":   "test2",
		"name": "test-item",
		"data": map[string]interface{}{
			"value": 42,
		},
	})

	if err := tester.AssertJSON("name", "test-item", resp); err != nil {
		t.Errorf("AssertJSON name: %v", err)
	}
}

func TestTesterRunTests(t *testing.T) {
	s := NewTestServer()
	if err := s.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer s.Close()

	tester := NewAPITester(s.URL())
	cases := []TestCase{
		{
			Name:   "create item",
			Method: "POST",
			Path:   "/items",
			Body:   map[string]interface{}{"id": "run1", "name": "run-item"},
			Assert: func(resp *http.Response) error {
				return tester.AssertStatus(http.StatusCreated, resp)
			},
		},
		{
			Name:   "get item",
			Method: "GET",
			Path:   "/items/run1",
			Assert: func(resp *http.Response) error {
				return tester.AssertStatus(http.StatusOK, resp)
			},
		},
		{
			Name:   "get nonexistent",
			Method: "GET",
			Path:   "/items/nope",
			Assert: func(resp *http.Response) error {
				return tester.AssertStatus(http.StatusNotFound, resp)
			},
		},
	}

	results := tester.RunTests(cases)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Passed {
			t.Errorf("test %q failed: %s", r.Name, r.Error)
		}
	}
}

func TestTesterRetry(t *testing.T) {
	s := NewTestServer()
	if err := s.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer s.Close()

	tester := NewAPITester(s.URL())
	tester.maxRetries = 1

	resp, err := tester.SendRequest("GET", "/items/missing", nil, nil)
	if err != nil {
		t.Fatalf("SendRequest: %v", err)
	}

	if err := tester.AssertStatus(http.StatusNotFound, resp); err != nil {
		t.Errorf("expected 404: %v", err)
	}
}
