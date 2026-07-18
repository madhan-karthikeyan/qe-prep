#include "word_count.hpp"
#include <cctype>

WCResult wc_count(std::istream &is) {
    WCResult result = {0, 0, 0};
    char c;
    bool in_word = false;

    while (is.get(c)) {
        result.chars++;
        if (c == '\n') result.lines++;
        if (std::isalnum(c) || c == '_' || c == '\'') {
            if (!in_word) {
                in_word = true;
                result.words++;
            }
        } else {
            in_word = false;
        }
    }
    return result;
}
