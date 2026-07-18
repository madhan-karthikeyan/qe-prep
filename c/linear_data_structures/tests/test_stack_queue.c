#include "stack.h"
#include "queue.h"
#include <assert.h>
#include <stdio.h>

void test_stack() {
    stack_t s;
    stack_init(&s, 5);
    assert(stack_is_empty(&s));

    stack_push(&s, 1);
    stack_push(&s, 2);
    stack_push(&s, 3);
    assert(!stack_is_empty(&s));
    assert(stack_peek(&s) == 3);
    assert(stack_pop(&s) == 3);
    assert(stack_pop(&s) == 2);
    assert(stack_pop(&s) == 1);
    assert(stack_is_empty(&s));

    stack_destroy(&s);
}

void test_queue() {
    queue_t q;
    queue_init(&q, 5);
    assert(queue_is_empty(&q));

    queue_enqueue(&q, 1);
    queue_enqueue(&q, 2);
    queue_enqueue(&q, 3);
    assert(!queue_is_empty(&q));
    assert(queue_peek(&q) == 1);
    assert(queue_dequeue(&q) == 1);
    assert(queue_dequeue(&q) == 2);
    assert(queue_dequeue(&q) == 3);
    assert(queue_is_empty(&q));

    queue_destroy(&q);
}

int main() {
    test_stack();
    test_queue();
    printf("PASS: stack_queue\n");
    return 0;
}
