#ifndef CONFIG_PARSER_H
#define CONFIG_PARSER_H

#define CFG_MAX_ENTRIES 64
#define CFG_MAX_KEY 64
#define CFG_MAX_VALUE 256

typedef struct {
    char keys[CFG_MAX_ENTRIES][CFG_MAX_KEY];
    char values[CFG_MAX_ENTRIES][CFG_MAX_VALUE];
    int count;
} config_t;

int config_load(config_t *cfg, const char *filename);
const char *config_get(config_t *cfg, const char *key);

#endif
