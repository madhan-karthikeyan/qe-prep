#include "producer_consumer.h"
#include <stdlib.h>

void bq_init(blocking_queue_t *q, int capacity) {
    q->buffer = malloc(capacity * sizeof(int));
    q->head = 0;
    q->tail = 0;
    q->capacity = capacity;
    q->size = 0;
    pthread_mutex_init(&q->mutex, NULL);
    pthread_cond_init(&q->not_full, NULL);
    pthread_cond_init(&q->not_empty, NULL);
}

void bq_put(blocking_queue_t *q, int val) {
    pthread_mutex_lock(&q->mutex);
    while (q->size >= q->capacity)
        pthread_cond_wait(&q->not_full, &q->mutex);
    q->buffer[q->tail] = val;
    q->tail = (q->tail + 1) % q->capacity;
    q->size++;
    pthread_cond_signal(&q->not_empty);
    pthread_mutex_unlock(&q->mutex);
}

int bq_get(blocking_queue_t *q) {
    pthread_mutex_lock(&q->mutex);
    while (q->size <= 0)
        pthread_cond_wait(&q->not_empty, &q->mutex);
    int val = q->buffer[q->head];
    q->head = (q->head + 1) % q->capacity;
    q->size--;
    pthread_cond_signal(&q->not_full);
    pthread_mutex_unlock(&q->mutex);
    return val;
}

void bq_destroy(blocking_queue_t *q) {
    free(q->buffer);
    pthread_mutex_destroy(&q->mutex);
    pthread_cond_destroy(&q->not_full);
    pthread_cond_destroy(&q->not_empty);
}
