package logging_toolkit

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// RotatingFileWriter is an io.WriteCloser that rotates log files when they
// exceed a size threshold. The old file is archived with a timestamp suffix.
type RotatingFileWriter struct {
	path      string
	maxSize   int64
	mu        sync.Mutex
	file      *os.File
	written   int64
}

// NewRotatingFileWriter creates a RotatingFileWriter that writes to the given
// path and rotates when the file exceeds maxSize bytes.
func NewRotatingFileWriter(path string, maxSize int64) (*RotatingFileWriter, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create log directory: %w", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, fmt.Errorf("stat log file: %w", err)
	}
	return &RotatingFileWriter{
		path:    path,
		maxSize: maxSize,
		file:    f,
		written: info.Size(),
	}, nil
}

// Write implements io.Writer. It rotates the file if the write would exceed the
// size threshold.
func (w *RotatingFileWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil {
		return 0, fmt.Errorf("writer is closed")
	}
	if w.written+int64(len(p)) > w.maxSize {
		if err := w.rotate(); err != nil {
			return 0, fmt.Errorf("rotate: %w", err)
		}
	}
	n, err := w.file.Write(p)
	w.written += int64(n)
	return n, err
}

// rotate closes the current file, renames it with a timestamp, and opens a new one.
func (w *RotatingFileWriter) rotate() error {
	if err := w.file.Close(); err != nil {
		return err
	}
	ts := time.Now().Format("20060102T150405.000")
	archive := w.path + "." + ts
	if err := os.Rename(w.path, archive); err != nil {
		return err
	}
	f, err := os.OpenFile(w.path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	w.file = f
	w.written = 0
	return nil
}

// Close implements io.Closer.
func (w *RotatingFileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file == nil {
		return nil
	}
	err := w.file.Close()
	w.file = nil
	return err
}
