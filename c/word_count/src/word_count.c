#include "word_count.h"
#include <ctype.h>

static int is_word_char(int c) {
    return isalnum(c) || c == '_' || c == '\'';
}

wc_result_t wc_count(FILE *fp) {
    wc_result_t result = {0, 0, 0};
    int c;
    int in_word = 0;

    while ((c = fgetc(fp)) != EOF) {
        result.chars++;
        if (c == '\n') result.lines++;
        if (is_word_char(c)) {
            if (!in_word) {
                in_word = 1;
                result.words++;
            }
        } else {
            in_word = 0;
        }
    }
    return result;
}
