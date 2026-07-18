#include "circular_queue.h"
#include <assert.h>
#include <stdio.h>

void test_cq() {
    circular_queue_t q;
    cq_init(&q, 3);

    assert(cq_is_empty(&q));
    assert(!cq_is_full(&q));

    assert(cq_enqueue(&q, 1) == 0);
    assert(cq_enqueue(&q, 2) == 0);
    assert(cq_enqueue(&q, 3) == 0);
    assert(cq_is_full(&q));
    assert(cq_enqueue(&q, 4) == -1);

    int val;
    assert(cq_dequeue(&q, &val) == 0 && val == 1);
    assert(!cq_is_full(&q));
    assert(cq_enqueue(&q, 4) == 0);

    assert(cq_dequeue(&q, &val) == 0 && val == 2);
    assert(cq_dequeue(&q, &val) == 0 && val == 3);
    assert(cq_dequeue(&q, &val) == 0 && val == 4);
    assert(cq_is_empty(&q));

    cq_destroy(&q);
}

int main() {
    test_cq();
    printf("PASS: circular_queue\n");
    return 0;
}
