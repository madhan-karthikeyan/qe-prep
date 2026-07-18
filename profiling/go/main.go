package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/madhan/qe-interview-playbook/go/file_processing"
	"github.com/madhan/qe-interview-playbook/go/logging_toolkit"
	"github.com/madhan/qe-interview-playbook/go/lru_cache"
	"github.com/madhan/qe-interview-playbook/go/rate_limiter"
)

func main() {
	// ---- CPU profiling ----
	cpuF, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer cpuF.Close()

	if err := pprof.StartCPUProfile(cpuF); err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	// ---- Memory profiling (before) ----
	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	// ---- Benchmark operations ----
	runLoggerOps()
	runLRUOps()
	runRateLimiterOps()
	runCSVParse()

	// ---- Memory profiling (after) ----
	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	memF, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer memF.Close()

	if err := pprof.WriteHeapProfile(memF); err != nil {
		log.Fatal(err)
	}

	// Print memory stats
	fmt.Println("=== Memory Profile Summary ===")
	fmt.Printf("Alloc (before): %d KB\n", memBefore.Alloc/1024)
	fmt.Printf("Alloc (after):  %d KB\n", memAfter.Alloc/1024)
	fmt.Printf("Total Alloc:    %d KB\n", memAfter.TotalAlloc/1024)
	fmt.Printf("Heap Inuse:     %d KB\n", memAfter.HeapInuse/1024)
	fmt.Printf("GC Cycles:      %d\n", memAfter.NumGC)

	fmt.Println("\nProfile files written: cpu.prof, mem.prof")
	fmt.Println("View with:")
	fmt.Println("  go tool pprof -http=:8080 cpu.prof")
	fmt.Println("  go tool pprof -http=:8080 mem.prof")
}

func runLoggerOps() {
	var buf bytes.Buffer
	writers := []io.Writer{&buf}
	logger := logging_toolkit.NewLogger(writers, logging_toolkit.DEBUG, logging_toolkit.DefaultFormatConfig())
	for i := 0; i < 10000; i++ {
		logger.Info("log message number %d", i)
	}
}

func runLRUOps() {
	cache := lru_cache.NewCache(1000)
	for i := 0; i < 50000; i++ {
		cache.Put(i, i)
		cache.Get(i % 1000)
	}
}

func runRateLimiterOps() {
	tb := rate_limiter.NewTokenBucket(10000, 10000)
	for i := 0; i < 100000; i++ {
		tb.Allow()
	}

	sw := rate_limiter.NewSlidingWindow(100000, time.Second)
	for i := 0; i < 100000; i++ {
		sw.Allow()
	}
}

func runCSVParse() {
	var b strings.Builder
	for j := 0; j < 10; j++ {
		if j > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, "col_%d", j)
	}
	b.WriteByte('\n')
	for i := 0; i < 1000; i++ {
		for j := 0; j < 10; j++ {
			if j > 0 {
				b.WriteByte(',')
			}
			fmt.Fprintf(&b, "%d", i*j)
		}
		b.WriteByte('\n')
	}
	for k := 0; k < 100; k++ {
		r := strings.NewReader(b.String())
		_, _ = file_processing.ParseCSV(r, ',')
	}
}
