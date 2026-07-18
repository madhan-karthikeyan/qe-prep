#include "logger.h"
#include <assert.h>
#include <stdio.h>

void test_logger_basic() {
    assert(log_init("/tmp/test_log.log") == 0);
    log_debug("debug message");
    log_info("info message");
    log_warn("warn message");
    log_error("error message");
    log_close();
}

int main() {
    test_logger_basic();
    printf("PASS: logger\n");
    return 0;
}
