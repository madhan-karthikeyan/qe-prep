#include "word_count.hpp"
#include <cassert>
#include <iostream>
#include <sstream>

void test_word_count() {
    std::istringstream ss("hello world\nfoo bar baz\n");
    WCResult res = wc_count(ss);
    assert(res.lines == 2);
    assert(res.words == 5);
    assert(res.chars > 0);
}

int main() {
    test_word_count();
    std::cout << "PASS: word_count" << std::endl;
    return 0;
}
