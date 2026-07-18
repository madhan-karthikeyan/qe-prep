#include "stack.h"
#include "queue.h"
#include <stdlib.h>

void stack_init(stack_t *s, int capacity) {
    s->data = malloc(capacity * sizeof(int));
    s->top = -1;
    s->capacity = capacity;
}

void stack_push(stack_t *s, int val) {
    if (s->top >= s->capacity - 1) return;
    s->data[++s->top] = val;
}

int stack_pop(stack_t *s) {
    if (s->top < 0) return -1;
    return s->data[s->top--];
}

int stack_peek(stack_t *s) {
    if (s->top < 0) return -1;
    return s->data[s->top];
}

int stack_is_empty(stack_t *s) {
    return s->top < 0;
}

void stack_destroy(stack_t *s) {
    free(s->data);
    s->data = NULL;
    s->top = -1;
}

void queue_init(queue_t *q, int capacity) {
    q->data = malloc(capacity * sizeof(int));
    q->head = 0;
    q->tail = 0;
    q->capacity = capacity;
    q->size = 0;
}

void queue_enqueue(queue_t *q, int val) {
    if (q->size >= q->capacity) return;
    q->data[q->tail] = val;
    q->tail = (q->tail + 1) % q->capacity;
    q->size++;
}

int queue_dequeue(queue_t *q) {
    if (q->size <= 0) return -1;
    int val = q->data[q->head];
    q->head = (q->head + 1) % q->capacity;
    q->size--;
    return val;
}

int queue_peek(queue_t *q) {
    if (q->size <= 0) return -1;
    return q->data[q->head];
}

int queue_is_empty(queue_t *q) {
    return q->size <= 0;
}

void queue_destroy(queue_t *q) {
    free(q->data);
    q->data = NULL;
    q->size = 0;
}
