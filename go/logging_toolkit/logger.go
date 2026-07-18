package logging_toolkit

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Level represents a log level.
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var levelNames = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
}

// FormatConfig controls the output format of log entries.
type FormatConfig struct {
	IncludeTime  bool
	IncludeLevel bool
	TimeFormat   string
}

// DefaultFormatConfig returns a FormatConfig with sensible defaults.
func DefaultFormatConfig() FormatConfig {
	return FormatConfig{
		IncludeTime:  true,
		IncludeLevel: true,
		TimeFormat:   time.RFC3339,
	}
}

// Logger provides thread-safe logging to multiple io.Writer targets.
type Logger struct {
	mu      sync.Mutex
	writers []io.Writer
	level   Level
	format  FormatConfig
}

// NewLogger creates a new Logger with the given writers and minimum level.
func NewLogger(writers []io.Writer, level Level, format FormatConfig) *Logger {
	if len(writers) == 0 {
		writers = []io.Writer{io.Discard}
	}
	return &Logger{
		writers: writers,
		level:   level,
		format:  format,
	}
}

// log writes a formatted message at the given level if it meets the minimum threshold.
func (l *Logger) log(level Level, format string, args ...any) {
	if level < l.level {
		return
	}
	msg := fmt.Sprintf(format, args...)
	var prefix string
	if l.format.IncludeTime {
		prefix += time.Now().Format(l.format.TimeFormat) + " "
	}
	if l.format.IncludeLevel {
		prefix += levelNames[level] + " "
	}
	line := prefix + msg + "\n"
	l.mu.Lock()
	for _, w := range l.writers {
		_, _ = w.Write([]byte(line))
	}
	l.mu.Unlock()
}

// Debug logs a message at DEBUG level.
func (l *Logger) Debug(format string, args ...any) {
	l.log(DEBUG, format, args...)
}

// Info logs a message at INFO level.
func (l *Logger) Info(format string, args ...any) {
	l.log(INFO, format, args...)
}

// Warn logs a message at WARN level.
func (l *Logger) Warn(format string, args ...any) {
	l.log(WARN, format, args...)
}

// Error logs a message at ERROR level.
func (l *Logger) Error(format string, args ...any) {
	l.log(ERROR, format, args...)
}

// SetLevel changes the minimum log level atomically.
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}
