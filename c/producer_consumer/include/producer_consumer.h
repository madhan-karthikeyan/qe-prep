#ifndef PRODUCER_CONSUMER_H
#define PRODUCER_CONSUMER_H

#include <pthread.h>

typedef struct {
    int *buffer;
    int head;
    int tail;
    int capacity;
    int size;
    pthread_mutex_t mutex;
    pthread_cond_t not_full;
    pthread_cond_t not_empty;
} blocking_queue_t;

void bq_init(blocking_queue_t *q, int capacity);
void bq_put(blocking_queue_t *q, int val);
int bq_get(blocking_queue_t *q);
void bq_destroy(blocking_queue_t *q);

#endif
