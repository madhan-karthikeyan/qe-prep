#include "logger.hpp"
#include <ctime>
#include <iomanip>
#include <sstream>
#include <iostream>

Logger::Logger(const std::string &filename) {
    file_.open(filename, std::ios::app);
}

Logger::~Logger() {
    if (file_.is_open()) file_.close();
}

std::string Logger::level_str(LogLevel level) {
    switch (level) {
        case LogLevel::DEBUG: return "DEBUG";
        case LogLevel::INFO:  return "INFO";
        case LogLevel::WARN:  return "WARN";
        case LogLevel::ERROR: return "ERROR";
    }
    return "UNKNOWN";
}

void Logger::log(LogLevel level, const std::string &msg) {
    std::lock_guard<std::mutex> lock(mtx_);

    auto now = std::time(nullptr);
    auto tm = *std::localtime(&now);
    std::ostringstream ts;
    ts << std::put_time(&tm, "%Y-%m-%d %H:%M:%S");

    std::string line = "[" + ts.str() + "] [" + level_str(level) + "] " + msg + "\n";

    if (file_.is_open()) {
        file_ << line;
        file_.flush();
    }
    std::cout << line;
    std::cout.flush();
}
