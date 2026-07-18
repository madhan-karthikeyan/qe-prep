package com.qe.test.cache;

import java.util.HashMap;
import java.util.Map;
import java.util.Objects;
import java.util.concurrent.locks.ReentrantLock;
import java.util.function.BiConsumer;

public class LRUCache<K, V> {
    private final int capacity;
    private final boolean threadSafe;
    private final Map<K, Node<K, V>> map;
    private final Node<K, V> head;
    private final Node<K, V> tail;
    private final ReentrantLock lock;

    public LRUCache(int capacity) {
        this(capacity, false);
    }

    public LRUCache(int capacity, boolean threadSafe) {
        if (capacity <= 0) throw new IllegalArgumentException("capacity must be positive");
        this.capacity = capacity;
        this.threadSafe = threadSafe;
        this.map = new HashMap<>(capacity);
        this.lock = threadSafe ? new ReentrantLock() : null;
        this.head = new Node<>(null, null);
        this.tail = new Node<>(null, null);
        head.next = tail;
        tail.prev = head;
    }

    public V get(K key) {
        if (threadSafe) {
            lock.lock();
            try {
                return getInternal(key);
            } finally {
                lock.unlock();
            }
        }
        return getInternal(key);
    }

    public void put(K key, V value) {
        Objects.requireNonNull(key, "key must not be null");
        if (threadSafe) {
            lock.lock();
            try {
                putInternal(key, value);
            } finally {
                lock.unlock();
            }
        } else {
            putInternal(key, value);
        }
    }

    public boolean contains(K key) {
        if (threadSafe) {
            lock.lock();
            try {
                return map.containsKey(key);
            } finally {
                lock.unlock();
            }
        }
        return map.containsKey(key);
    }

    public int size() {
        if (threadSafe) {
            lock.lock();
            try {
                return map.size();
            } finally {
                lock.unlock();
            }
        }
        return map.size();
    }

    public void forEach(BiConsumer<? super K, ? super V> action) {
        Objects.requireNonNull(action);
        if (threadSafe) {
            lock.lock();
            try {
                forEachInternal(action);
            } finally {
                lock.unlock();
            }
        } else {
            forEachInternal(action);
        }
    }

    private V getInternal(K key) {
        var node = map.get(key);
        if (node == null) return null;
        moveToHead(node);
        return node.value;
    }

    private void putInternal(K key, V value) {
        var existing = map.get(key);
        if (existing != null) {
            existing.value = value;
            moveToHead(existing);
            return;
        }
        if (map.size() >= capacity) {
            evict();
        }
        var node = new Node<>(key, value);
        map.put(key, node);
        addToHead(node);
    }

    private void evict() {
        var lru = tail.prev;
        if (lru == head) return;
        removeNode(lru);
        map.remove(lru.key);
    }

    private void moveToHead(Node<K, V> node) {
        removeNode(node);
        addToHead(node);
    }

    private void addToHead(Node<K, V> node) {
        node.next = head.next;
        node.prev = head;
        head.next.prev = node;
        head.next = node;
    }

    private void removeNode(Node<K, V> node) {
        node.prev.next = node.next;
        node.next.prev = node.prev;
    }

    private void forEachInternal(BiConsumer<? super K, ? super V> action) {
        var cur = head.next;
        while (cur != tail) {
            action.accept(cur.key, cur.value);
            cur = cur.next;
        }
    }

    private static class Node<K, V> {
        final K key;
        V value;
        Node<K, V> prev;
        Node<K, V> next;

        Node(K key, V value) {
            this.key = key;
            this.value = value;
        }
    }
}
