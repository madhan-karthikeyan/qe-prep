package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/madhan/qe-interview-playbook/go/lru_cache"
)

func main() {
	// ---- CPU Profiling ----
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	// Run benchmark-style operations
	cache := lru_cache.NewCache(5000)
	for i := 0; i < 200000; i++ {
		cache.Put(i, i)
		if i%2 == 0 {
			cache.Get(i)
		}
	}

	// Simulate string operations
	var b strings.Builder
	for i := 0; i < 100000; i++ {
		b.WriteString(fmt.Sprintf("iteration %d\n", i))
	}
	_ = b.String()

	fmt.Println("CPU profile written to cpu.prof")
	fmt.Println("View with: go tool pprof -http=:8080 cpu.prof")
}
