#include "circular_queue.h"
#include <stdlib.h>

void cq_init(circular_queue_t *q, int capacity) {
    q->buffer = malloc(capacity * sizeof(int));
    q->head = 0;
    q->tail = 0;
    q->capacity = capacity;
    q->full = 0;
}

int cq_enqueue(circular_queue_t *q, int val) {
    if (q->full) return -1;
    q->buffer[q->tail] = val;
    q->tail = (q->tail + 1) % q->capacity;
    q->full = (q->tail == q->head);
    return 0;
}

int cq_dequeue(circular_queue_t *q, int *val) {
    if (cq_is_empty(q)) return -1;
    *val = q->buffer[q->head];
    q->head = (q->head + 1) % q->capacity;
    q->full = 0;
    return 0;
}

int cq_is_full(circular_queue_t *q) {
    return q->full;
}

int cq_is_empty(circular_queue_t *q) {
    return !q->full && q->head == q->tail;
}

void cq_destroy(circular_queue_t *q) {
    free(q->buffer);
    q->buffer = NULL;
}
