#include "logger.hpp"
#include <cassert>
#include <iostream>

void test_logger_basic() {
    Logger log("/tmp/test_log.txt");
    log.log(LogLevel::INFO, "test message");
}

int main() {
    test_logger_basic();
    std::cout << "PASS: logger" << std::endl;
    return 0;
}
