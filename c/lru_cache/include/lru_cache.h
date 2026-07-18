#ifndef LRU_CACHE_H
#define LRU_CACHE_H

typedef struct lru_node {
    int key;
    int value;
    struct lru_node *prev;
    struct lru_node *next;
} lru_node_t;

typedef struct {
    lru_node_t **hashmap;
    lru_node_t *head;
    lru_node_t *tail;
    int capacity;
    int size;
} lru_cache_t;

lru_cache_t *lru_create(int capacity);
int lru_get(lru_cache_t *cache, int key);
void lru_put(lru_cache_t *cache, int key, int value);
void lru_free(lru_cache_t *cache);

#endif
