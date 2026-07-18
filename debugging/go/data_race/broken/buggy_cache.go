package main

import (
	"fmt"
	"sync"
)

type Cache struct {
	data map[int]string
}

func (c *Cache) Set(key int, value string) {
	c.data[key] = value
}

func (c *Cache) Get(key int) string {
	return c.data[key]
}

func main() {
	cache := &Cache{data: make(map[int]string)}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cache.Set(i, fmt.Sprintf("value-%d", i))
		}(i)
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = cache.Get(i)
		}(i)
	}

	wg.Wait()
	fmt.Println("cache operations completed")
}
