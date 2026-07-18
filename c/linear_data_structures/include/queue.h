#ifndef QUEUE_H
#define QUEUE_H

typedef struct {
    int *data;
    int head;
    int tail;
    int capacity;
    int size;
} queue_t;

void queue_init(queue_t *q, int capacity);
void queue_enqueue(queue_t *q, int val);
int queue_dequeue(queue_t *q);
int queue_peek(queue_t *q);
int queue_is_empty(queue_t *q);
void queue_destroy(queue_t *q);

#endif
