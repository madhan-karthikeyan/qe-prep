package logging_toolkit

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
)

func TestLoggerLevels(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	log := NewLogger([]io.Writer{&buf}, INFO, DefaultFormatConfig())
	log.Debug("debug msg")
	if buf.Len() > 0 {
		t.Error("expected no DEBUG output with INFO minimum")
	}
	log.Info("info msg")
	if !strings.Contains(buf.String(), "INFO") {
		t.Error("expected INFO output")
	}
}

func TestLoggerMultiWriter(t *testing.T) {
	t.Parallel()
	var buf1, buf2 bytes.Buffer
	log := NewLogger([]io.Writer{&buf1, &buf2}, DEBUG, DefaultFormatConfig())
	log.Info("hello %s", "world")
	if !strings.Contains(buf1.String(), "hello world") {
		t.Error("buf1 missing message")
	}
	if !strings.Contains(buf2.String(), "hello world") {
		t.Error("buf2 missing message")
	}
}

func TestLoggerFormat(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		format FormatConfig
		check  func(t *testing.T, s string)
	}{
		{
			"both",
			FormatConfig{IncludeTime: true, IncludeLevel: true, TimeFormat: "15:04:05"},
			func(t *testing.T, s string) {
				if !strings.Contains(s, "INFO") || !strings.Contains(s, ":") {
					t.Errorf("expected level+time, got %q", s)
				}
			},
		},
		{
			"level only",
			FormatConfig{IncludeTime: false, IncludeLevel: true},
			func(t *testing.T, s string) {
				if !strings.HasPrefix(s, "INFO") {
					t.Errorf("expected prefix INFO, got %q", s)
				}
			},
		},
		{
			"message only",
			FormatConfig{IncludeTime: false, IncludeLevel: false},
			func(t *testing.T, s string) {
				if s != "hey\n" {
					t.Errorf("expected \"hey\\n\", got %q", s)
				}
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			log := NewLogger([]io.Writer{&buf}, DEBUG, tc.format)
			log.Info("hey")
			tc.check(t, buf.String())
		})
	}
}

func TestLoggerConcurrent(t *testing.T) {
	var buf bytes.Buffer
	log := NewLogger([]io.Writer{&buf}, DEBUG, DefaultFormatConfig())
	var wg sync.WaitGroup
	for i := range 20 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			log.Info("goroutine %d", n)
		}(i)
	}
	wg.Wait()
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 20 {
		t.Errorf("expected 20 lines, got %d", len(lines))
	}
}

func TestRotatingFileWriter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")
	maxSize := int64(100)
	w, err := NewRotatingFileWriter(path, maxSize)
	if err != nil {
		t.Fatal(err)
	}
	data := []byte(strings.Repeat("a", 80))
	for range 3 {
		_, err := w.Write(data)
		if err != nil {
			t.Fatal(err)
		}
	}
	w.Close()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) < 2 {
		t.Fatalf("expected at least 2 files (rotated), got %d", len(entries))
	}
}

func TestRotatingFileWriterClose(t *testing.T) {
	w, err := NewRotatingFileWriter(filepath.Join(t.TempDir(), "close.log"), 1024)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	// second close should not panic
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestRotatingFileWriterWriteAfterClose(t *testing.T) {
	w, err := NewRotatingFileWriter(filepath.Join(t.TempDir(), "closed.log"), 1024)
	if err != nil {
		t.Fatal(err)
	}
	w.Close()
	_, err = w.Write([]byte("test"))
	if err == nil {
		t.Error("expected error writing to closed writer")
	}
}

func TestFilter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		filter   Filter
		level    Level
		msg      string
		expected bool
	}{
		{
			"min level passes",
			Filter{MinLevel: INFO},
			INFO,
			"ok",
			true,
		},
		{
			"below min level",
			Filter{MinLevel: WARN},
			INFO,
			"skip",
			false,
		},
		{
			"regex match",
			Filter{MinLevel: DEBUG, Pattern: regexp.MustCompile("error")},
			ERROR,
			"something error happened",
			true,
		},
		{
			"regex no match",
			Filter{MinLevel: DEBUG, Pattern: regexp.MustCompile("error")},
			INFO,
			"all good",
			false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.filter.ShouldLog(tc.level, tc.msg)
			if got != tc.expected {
				t.Errorf("ShouldLog(%v, %q) = %v, want %v", tc.level, tc.msg, got, tc.expected)
			}
		})
	}
}

func TestIntegrationMultiGoroutineRotation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "integ.log")
	maxSize := int64(512)
	w, err := NewRotatingFileWriter(path, maxSize)
	if err != nil {
		t.Fatal(err)
	}
	log := NewLogger([]io.Writer{w}, DEBUG, DefaultFormatConfig())
	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for j := range 20 {
				log.Debug("goroutine %d message %d: %s", n, j, strings.Repeat("x", 30))
			}
		}(i)
	}
	wg.Wait()
	w.Close()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) < 2 {
		t.Errorf("expected rotation to produce multiple files, got %d", len(entries))
	}
}
