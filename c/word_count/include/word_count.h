#ifndef WORD_COUNT_H
#define WORD_COUNT_H

#include <stdio.h>

typedef struct {
    long lines;
    long words;
    long chars;
} wc_result_t;

wc_result_t wc_count(FILE *fp);

#endif
