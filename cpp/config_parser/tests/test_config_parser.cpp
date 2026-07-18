#include "config_parser.hpp"
#include <cassert>
#include <iostream>
#include <fstream>

void test_config_parser() {
    std::ofstream out("/tmp/test_config.env");
    out << "DB_HOST=localhost\nDB_PORT=5432\n# comment\nEMPTY=\n";
    out.close();

    ConfigParser cfg;
    assert(cfg.load("/tmp/test_config.env"));
    assert(cfg.get("DB_HOST") == "localhost");
    assert(cfg.get("DB_PORT") == "5432");
    assert(cfg.get("EMPTY") == "");
    assert(cfg.get("NONEXIST") == "");
}

int main() {
    test_config_parser();
    std::cout << "PASS: config_parser" << std::endl;
    return 0;
}
