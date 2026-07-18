package benchmarks

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"testing"

	"github.com/madhan/qe-interview-playbook/go/file_processing"
	"github.com/madhan/qe-interview-playbook/go/logging_toolkit"
	"github.com/madhan/qe-interview-playbook/go/lru_cache"
	"github.com/madhan/qe-interview-playbook/go/rate_limiter"
	"github.com/madhan/qe-interview-playbook/go/thread_pool"
)

// ---------------------------------------------------------------------------
// Logger benchmarks
// ---------------------------------------------------------------------------

func BenchmarkLogger(b *testing.B) {
	var buf bytes.Buffer
	writers := []io.Writer{&buf}
	logger := logging_toolkit.NewLogger(writers, logging_toolkit.DEBUG, logging_toolkit.DefaultFormatConfig())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("test message %d", i)
	}
}

func BenchmarkLoggerParallel(b *testing.B) {
	var buf bytes.Buffer
	writers := []io.Writer{&buf}
	logger := logging_toolkit.NewLogger(writers, logging_toolkit.DEBUG, logging_toolkit.DefaultFormatConfig())
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			logger.Info("parallel message %d", i)
			i++
		}
	})
}

func BenchmarkLoggerWithLevelFilter(b *testing.B) {
	var buf bytes.Buffer
	writers := []io.Writer{&buf}
	logger := logging_toolkit.NewLogger(writers, logging_toolkit.ERROR, logging_toolkit.DefaultFormatConfig())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("filtered message %d", i)
	}
}

// ---------------------------------------------------------------------------
// LRU Cache benchmarks
// ---------------------------------------------------------------------------

func BenchmarkLRUPut(b *testing.B) {
	cache := lru_cache.NewCache(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Put(i, i)
	}
}

func BenchmarkLRUGet(b *testing.B) {
	cache := lru_cache.NewCache(b.N)
	for i := 0; i < b.N; i++ {
		cache.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(i)
	}
}

func BenchmarkLRUGetMiss(b *testing.B) {
	cache := lru_cache.NewCache(1000)
	for i := 0; i < 1000; i++ {
		cache.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(i + 1000)
	}
}

func BenchmarkLRUMixed(b *testing.B) {
	cache := lru_cache.NewCache(b.N / 2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Put(i, i)
		cache.Get(i / 2)
	}
}

// ---------------------------------------------------------------------------
// CSV parsing benchmarks
// ---------------------------------------------------------------------------

func BenchmarkCSVSmall(b *testing.B) {
	data := makeCSVData(100, 10)
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(data)
		_, _ = file_processing.ParseCSV(r, ',')
	}
}

func BenchmarkCSVMedium(b *testing.B) {
	data := makeCSVData(1000, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(data)
		_, _ = file_processing.ParseCSV(r, ',')
	}
}

func BenchmarkCSVLarge(b *testing.B) {
	data := makeCSVData(10000, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(data)
		_, _ = file_processing.ParseCSV(r, ',')
	}
}

func makeCSVData(rows, cols int) string {
	var b strings.Builder
	for i := 0; i < cols; i++ {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, "col_%d", i)
	}
	b.WriteByte('\n')
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if j > 0 {
				b.WriteByte(',')
			}
			fmt.Fprintf(&b, "%d", i*j)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

// ---------------------------------------------------------------------------
// Rate limiter benchmarks
// ---------------------------------------------------------------------------

func BenchmarkTokenBucket(b *testing.B) {
	tb := rate_limiter.NewTokenBucket(b.N, float64(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tb.Allow()
	}
}

func BenchmarkSlidingWindow(b *testing.B) {
	sw := rate_limiter.NewSlidingWindow(b.N, 1e9) // 1 second window
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sw.Allow()
	}
}

func BenchmarkTokenBucketParallel(b *testing.B) {
	tb := rate_limiter.NewTokenBucket(b.N, float64(b.N))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tb.Allow()
		}
	})
}

// ---------------------------------------------------------------------------
// Thread pool benchmarks
// ---------------------------------------------------------------------------

func BenchmarkThreadPoolSubmit(b *testing.B) {
	pool := thread_pool.New(4, b.N)
	counter := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool.Submit(func() {
			counter++
		})
	}
	pool.Stop()
	pool.Wait()
}

func BenchmarkThreadPoolParallelSubmit(b *testing.B) {
	pool := thread_pool.New(4, b.N)
	counter := 0
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Submit(func() {
				counter++
			})
		}
	})
	pool.Stop()
	pool.Wait()
}

// ---------------------------------------------------------------------------
// Random utility benchmarks
// ---------------------------------------------------------------------------

func BenchmarkRandInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = rand.Intn(1000)
	}
}

func BenchmarkSprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("hello %d world %f", i, float64(i))
	}
}
