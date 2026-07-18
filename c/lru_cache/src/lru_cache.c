#include "lru_cache.h"
#include <stdlib.h>
#include <string.h>

#define HASHMAP_SIZE 1024

lru_cache_t *lru_create(int capacity) {
    lru_cache_t *cache = malloc(sizeof(lru_cache_t));
    cache->hashmap = calloc(HASHMAP_SIZE, sizeof(lru_node_t *));
    cache->head = NULL;
    cache->tail = NULL;
    cache->capacity = capacity;
    cache->size = 0;
    return cache;
}

static void detach_node(lru_cache_t *cache, lru_node_t *node) {
    if (node->prev) node->prev->next = node->next;
    if (node->next) node->next->prev = node->prev;
    if (cache->head == node) cache->head = node->next;
    if (cache->tail == node) cache->tail = node->prev;
    node->prev = NULL;
    node->next = NULL;
}

static void move_to_front(lru_cache_t *cache, lru_node_t *node) {
    detach_node(cache, node);
    node->next = cache->head;
    if (cache->head) cache->head->prev = node;
    cache->head = node;
    if (!cache->tail) cache->tail = node;
}

static lru_node_t *remove_tail(lru_cache_t *cache) {
    lru_node_t *node = cache->tail;
    if (!node) return NULL;
    cache->tail = node->prev;
    if (cache->tail) cache->tail->next = NULL;
    else cache->head = NULL;
    return node;
}

int lru_get(lru_cache_t *cache, int key) {
    int idx = (unsigned)key % HASHMAP_SIZE;
    lru_node_t *node = cache->hashmap[idx];
    while (node) {
        if (node->key == key) {
            move_to_front(cache, node);
            return node->value;
        }
        node = node->next;
    }
    return -1;
}

void lru_put(lru_cache_t *cache, int key, int value) {
    int idx = (unsigned)key % HASHMAP_SIZE;
    lru_node_t *node = cache->hashmap[idx];
    while (node) {
        if (node->key == key) {
            node->value = value;
            move_to_front(cache, node);
            return;
        }
        node = node->next;
    }

    if (cache->size >= cache->capacity) {
        lru_node_t *evict = remove_tail(cache);
        if (evict) {
            int evict_idx = (unsigned)evict->key % HASHMAP_SIZE;
            lru_node_t **p = &cache->hashmap[evict_idx];
            while (*p) {
                if (*p == evict) { *p = evict->next; break; }
                p = &(*p)->next;
            }
            free(evict);
            cache->size--;
        }
    }

    lru_node_t *new_node = malloc(sizeof(lru_node_t));
    new_node->key = key;
    new_node->value = value;
    new_node->prev = NULL;
    new_node->next = cache->head;
    if (cache->head) cache->head->prev = new_node;
    cache->head = new_node;
    if (!cache->tail) cache->tail = new_node;

    new_node->next = cache->hashmap[idx];
    cache->hashmap[idx] = new_node;
    cache->size++;
}

void lru_free(lru_cache_t *cache) {
    lru_node_t *node = cache->head;
    while (node) {
        lru_node_t *next = node->next;
        free(node);
        node = next;
    }
    free(cache->hashmap);
    free(cache);
}
