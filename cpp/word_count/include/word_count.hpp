#ifndef WORD_COUNT_HPP
#define WORD_COUNT_HPP

#include <istream>

struct WCResult {
    long lines;
    long words;
    long chars;
};

WCResult wc_count(std::istream &is);

#endif
