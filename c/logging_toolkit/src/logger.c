#include "logger.h"
#include <time.h>
#include <string.h>

static FILE *log_file = NULL;

static const char *level_label(log_level_t level) {
    switch (level) {
        case LOG_DEBUG: return "DEBUG";
        case LOG_INFO:  return "INFO";
        case LOG_WARN:  return "WARN";
        case LOG_ERROR: return "ERROR";
    }
    return "UNKNOWN";
}

int log_init(const char *filename) {
    log_file = fopen(filename, "a");
    return log_file ? 0 : -1;
}

void log_msg(log_level_t level, const char *msg) {
    if (!log_file) return;

    time_t now = time(NULL);
    struct tm *tm_info = localtime(&now);
    char timestamp[20];
    strftime(timestamp, sizeof(timestamp), "%Y-%m-%d %H:%M:%S", tm_info);

    fprintf(log_file, "[%s] [%s] %s\n", timestamp, level_label(level), msg);
    fflush(log_file);

    fprintf(stdout, "[%s] [%s] %s\n", timestamp, level_label(level), msg);
    fflush(stdout);
}

void log_close(void) {
    if (log_file) {
        fclose(log_file);
        log_file = NULL;
    }
}
