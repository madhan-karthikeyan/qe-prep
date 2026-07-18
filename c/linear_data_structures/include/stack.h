#ifndef STACK_H
#define STACK_H

typedef struct {
    int *data;
    int top;
    int capacity;
} stack_t;

void stack_init(stack_t *s, int capacity);
void stack_push(stack_t *s, int val);
int stack_pop(stack_t *s);
int stack_peek(stack_t *s);
int stack_is_empty(stack_t *s);
void stack_destroy(stack_t *s);

#endif
