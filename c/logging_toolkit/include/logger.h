#ifndef LOGGER_H
#define LOGGER_H

#include <stdio.h>

typedef enum {
    LOG_DEBUG,
    LOG_INFO,
    LOG_WARN,
    LOG_ERROR
} log_level_t;

int log_init(const char *filename);
void log_msg(log_level_t level, const char *msg);
void log_close(void);

#define log_debug(msg) log_msg(LOG_DEBUG, msg)
#define log_info(msg)  log_msg(LOG_INFO, msg)
#define log_warn(msg)  log_msg(LOG_WARN, msg)
#define log_error(msg) log_msg(LOG_ERROR, msg)

#endif
