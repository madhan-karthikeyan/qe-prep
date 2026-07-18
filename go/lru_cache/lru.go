package lru_cache

import (
	"sync"
)

// Node is a node in the doubly-linked list.
type Node struct {
	Key   int
	Value int
	Prev  *Node
	Next  *Node
}

// Cache implements an LRU cache with O(1) Get and Put operations.
type Cache struct {
	capacity int
	items    map[int]*Node
	head     *Node
	tail     *Node
	mu       sync.Mutex
}

// NewCache creates a new LRU cache with the given capacity.
func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		items:    make(map[int]*Node, capacity),
	}
}

// remove removes a node from the linked list.
func (c *Cache) remove(node *Node) {
	if node.Prev != nil {
		node.Prev.Next = node.Next
	} else {
		c.head = node.Next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	} else {
		c.tail = node.Prev
	}
	node.Prev = nil
	node.Next = nil
}

// insertFront inserts a node at the front of the linked list.
func (c *Cache) insertFront(node *Node) {
	node.Next = c.head
	node.Prev = nil
	if c.head != nil {
		c.head.Prev = node
	}
	c.head = node
	if c.tail == nil {
		c.tail = node
	}
}

// moveToFront moves an existing node to the front.
func (c *Cache) moveToFront(node *Node) {
	c.remove(node)
	c.insertFront(node)
}

// removeOldest removes and returns the least recently used node.
func (c *Cache) removeOldest() *Node {
	if c.tail == nil {
		return nil
	}
	node := c.tail
	c.remove(node)
	return node
}

// Get retrieves a value by key. The second return value indicates whether the
// key was found.
func (c *Cache) Get(key int) (int, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, ok := c.items[key]
	if !ok {
		return 0, false
	}
	c.moveToFront(node)
	return node.Value, true
}

// Put inserts or updates a key-value pair. If the cache is full, the least
// recently used entry is evicted.
func (c *Cache) Put(key, value int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.capacity == 0 {
		return
	}

	if node, ok := c.items[key]; ok {
		node.Value = value
		c.moveToFront(node)
		return
	}

	node := &Node{Key: key, Value: value}
	if len(c.items) >= c.capacity {
		evicted := c.removeOldest()
		if evicted != nil {
			delete(c.items, evicted.Key)
		}
	}
	c.insertFront(node)
	c.items[key] = node
}

// Len returns the number of items in the cache.
func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}
