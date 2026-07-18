package rest_api_tester

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestServerGetNotFound(t *testing.T) {
	s := NewTestServer()
	if err := s.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer s.Close()

	resp, err := http.Get(s.URL() + "/items/nonexistent")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestServerPostAndGet(t *testing.T) {
	s := NewTestServer()
	if err := s.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer s.Close()

	item := `{"id":"item1","name":"test-item","value":42}`
	resp, err := http.Post(s.URL()+"/items", "application/json", strings.NewReader(item))
	if err != nil {
		t.Fatalf("Post: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	resp, err = http.Get(s.URL() + "/items/item1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if result["name"] != "test-item" {
		t.Errorf("name = %v", result["name"])
	}
}

func TestServerDelete(t *testing.T) {
	s := NewTestServer()
	if err := s.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer s.Close()

	item := `{"id":"del-item","data":"to-delete"}`
	http.Post(s.URL()+"/items", "application/json", strings.NewReader(item))

	req, _ := http.NewRequest("DELETE", s.URL()+"/items/del-item", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", resp.StatusCode)
	}

	resp, err = http.Get(s.URL() + "/items/del-item")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", resp.StatusCode)
	}
}

func TestServerBadRequest(t *testing.T) {
	s := NewTestServer()
	if err := s.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer s.Close()

	resp, err := http.Post(s.URL()+"/items", "application/json", strings.NewReader(`invalid json`))
	if err != nil {
		t.Fatalf("Post: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}
