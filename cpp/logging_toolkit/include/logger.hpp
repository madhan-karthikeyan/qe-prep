#ifndef LOGGER_HPP
#define LOGGER_HPP

#include <string>
#include <fstream>
#include <mutex>

enum class LogLevel { DEBUG, INFO, WARN, ERROR };

class Logger {
public:
    Logger(const std::string &filename);
    ~Logger();
    void log(LogLevel level, const std::string &msg);

private:
    std::ofstream file_;
    std::mutex mtx_;
    std::string level_str(LogLevel level);
};

#endif
