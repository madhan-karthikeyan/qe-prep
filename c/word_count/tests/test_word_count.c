#include "word_count.h"
#include <assert.h>
#include <stdio.h>

void test_word_count() {
    FILE *fp = fopen("/tmp/test_wc.txt", "w+");
    assert(fp);
    fprintf(fp, "hello world\nfoo bar baz\n");
    fflush(fp);
    rewind(fp);

    wc_result_t res = wc_count(fp);
    assert(res.lines == 2);
    assert(res.words == 5);
    assert(res.chars > 0);

    fclose(fp);
}

int main() {
    test_word_count();
    printf("PASS: word_count\n");
    return 0;
}
