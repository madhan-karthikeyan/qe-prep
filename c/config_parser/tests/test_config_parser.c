#include "config_parser.h"
#include <assert.h>
#include <stdio.h>
#include <string.h>

void test_config_parser() {
    config_t cfg;
    FILE *fp = fopen("/tmp/test_config.env", "w");
    assert(fp);
    fprintf(fp, "DB_HOST=localhost\nDB_PORT=5432\n# comment\nEMPTY=\n");
    fclose(fp);

    assert(config_load(&cfg, "/tmp/test_config.env") == 0);
    assert(cfg.count == 3);
    assert(strcmp(config_get(&cfg, "DB_HOST"), "localhost") == 0);
    assert(strcmp(config_get(&cfg, "DB_PORT"), "5432") == 0);
    assert(strcmp(config_get(&cfg, "EMPTY"), "") == 0);
    assert(config_get(&cfg, "NONEXIST") == NULL);
}

int main() {
    test_config_parser();
    printf("PASS: config_parser\n");
    return 0;
}
