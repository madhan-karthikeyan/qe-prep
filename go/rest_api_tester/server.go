package rest_api_tester

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"
)

// TestServer is a minimal HTTP test server with in-memory storage and JSON
// request/response handling. Supports GET, POST, DELETE on /items.
type TestServer struct {
	server *http.Server
	data   map[string]interface{}
	mu     sync.RWMutex
	addr   string
}

// NewTestServer creates a new TestServer with empty storage.
func NewTestServer() *TestServer {
	return &TestServer{
		data: make(map[string]interface{}),
	}
}

// Start launches the HTTP server on a random available port.
func (s *TestServer) Start() error {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}
	s.addr = listener.Addr().String()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /items/{id}", s.getHandler)
	mux.HandleFunc("POST /items", s.postHandler)
	mux.HandleFunc("DELETE /items/{id}", s.deleteHandler)
	s.server = &http.Server{Handler: mux}
	go func() {
		s.server.Serve(listener)
	}()
	return nil
}

// Close shuts down the server.
func (s *TestServer) Close() error {
	return s.server.Close()
}

// Addr returns the server's listen address.
func (s *TestServer) Addr() string {
	return s.addr
}

// URL returns the full base URL for the server.
func (s *TestServer) URL() string {
	return "http://" + s.addr
}

func (s *TestServer) getHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s.mu.RLock()
	item, ok := s.data[id]
	s.mu.RUnlock()
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *TestServer) postHandler(w http.ResponseWriter, r *http.Request) {
	var item map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad request"})
		return
	}
	id, _ := item["id"].(string)
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id required"})
		return
	}
	s.mu.Lock()
	s.data[id] = item
	s.mu.Unlock()
	writeJSON(w, http.StatusCreated, item)
}

func (s *TestServer) deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s.mu.Lock()
	delete(s.data, id)
	s.mu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
