#include "config_parser.h"
#include <stdio.h>
#include <string.h>

int config_load(config_t *cfg, const char *filename) {
    FILE *fp = fopen(filename, "r");
    if (!fp) return -1;

    cfg->count = 0;
    char line[CFG_MAX_VALUE];
    while (fgets(line, sizeof(line), fp) && cfg->count < CFG_MAX_ENTRIES) {
        char *p = line;
        while (*p == ' ' || *p == '\t') p++;
        if (*p == '#' || *p == '\n' || *p == '\0') continue;

        char *eq = strchr(p, '=');
        if (!eq) continue;

        int key_len = eq - p;
        if (key_len >= CFG_MAX_KEY) key_len = CFG_MAX_KEY - 1;
        strncpy(cfg->keys[cfg->count], p, key_len);
        cfg->keys[cfg->count][key_len] = '\0';

        char *vstart = eq + 1;
        while (*vstart == ' ' || *vstart == '\t') vstart++;
        int vlen = strlen(vstart);
        while (vlen > 0 && (vstart[vlen - 1] == '\n' || vstart[vlen - 1] == '\r' || vstart[vlen - 1] == ' ' || vstart[vlen - 1] == '\t'))
            vlen--;
        if (vlen >= CFG_MAX_VALUE) vlen = CFG_MAX_VALUE - 1;
        strncpy(cfg->values[cfg->count], vstart, vlen);
        cfg->values[cfg->count][vlen] = '\0';

        cfg->count++;
    }
    fclose(fp);
    return 0;
}

const char *config_get(config_t *cfg, const char *key) {
    for (int i = 0; i < cfg->count; i++) {
        if (strcmp(cfg->keys[i], key) == 0)
            return cfg->values[i];
    }
    return NULL;
}
