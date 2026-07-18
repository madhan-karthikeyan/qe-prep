#ifndef CIRCULAR_QUEUE_H
#define CIRCULAR_QUEUE_H

typedef struct {
    int *buffer;
    int head;
    int tail;
    int capacity;
    int full;
} circular_queue_t;

void cq_init(circular_queue_t *q, int capacity);
int cq_enqueue(circular_queue_t *q, int val);
int cq_dequeue(circular_queue_t *q, int *val);
int cq_is_full(circular_queue_t *q);
int cq_is_empty(circular_queue_t *q);
void cq_destroy(circular_queue_t *q);

#endif
