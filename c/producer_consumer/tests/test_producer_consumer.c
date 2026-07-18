#include "producer_consumer.h"
#include <assert.h>
#include <stdio.h>
#include <pthread.h>

typedef struct {
    blocking_queue_t *q;
    int id;
} thread_arg_t;

void *producer(void *arg) {
    thread_arg_t *a = (thread_arg_t *)arg;
    bq_put(a->q, a->id);
    return NULL;
}

void *consumer(void *arg) {
    blocking_queue_t *q = (blocking_queue_t *)arg;
    int val = bq_get(q);
    (void)val;
    return NULL;
}

void test_pc() {
    blocking_queue_t q;
    bq_init(&q, 5);

    pthread_t prod[3], cons[3];
    thread_arg_t args[3];

    for (int i = 0; i < 3; i++) {
        args[i].q = &q;
        args[i].id = i + 1;
        pthread_create(&prod[i], NULL, producer, &args[i]);
    }

    for (int i = 0; i < 3; i++) {
        pthread_create(&cons[i], NULL, consumer, &q);
    }

    for (int i = 0; i < 3; i++) {
        pthread_join(prod[i], NULL);
        pthread_join(cons[i], NULL);
    }

    assert(q.size == 0);
    bq_destroy(&q);
}

int main() {
    test_pc();
    printf("PASS: producer_consumer\n");
    return 0;
}
