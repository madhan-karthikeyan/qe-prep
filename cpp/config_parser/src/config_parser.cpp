#include "config_parser.hpp"
#include <fstream>
#include <sstream>
#include <algorithm>

bool ConfigParser::load(const std::string &filename) {
    std::ifstream file(filename);
    if (!file.is_open()) return false;

    data_.clear();
    std::string line;
    while (std::getline(file, line)) {
        auto comment = line.find('#');
        if (comment != std::string::npos)
            line = line.substr(0, comment);

        auto eq = line.find('=');
        if (eq == std::string::npos) continue;

        std::string key = line.substr(0, eq);
        std::string value = line.substr(eq + 1);

        key.erase(0, key.find_first_not_of(" \t"));
        key.erase(key.find_last_not_of(" \t") + 1);
        value.erase(0, value.find_first_not_of(" \t"));
        value.erase(value.find_last_not_of(" \t\r\n") + 1);

        if (!key.empty())
            data_[key] = value;
    }
    return true;
}

std::string ConfigParser::get(const std::string &key) const {
    auto it = data_.find(key);
    if (it != data_.end()) return it->second;
    return "";
}
