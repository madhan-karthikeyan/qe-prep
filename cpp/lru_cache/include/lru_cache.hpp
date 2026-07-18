#ifndef LRU_CACHE_HPP
#define LRU_CACHE_HPP

#include <unordered_map>
#include <list>

template<typename K, typename V>
class LRUCache {
public:
    LRUCache(int capacity) : capacity_(capacity) {}

    V get(const K &key) {
        auto it = map_.find(key);
        if (it == map_.end()) return V{};
        order_.splice(order_.begin(), order_, it->second);
        return it->second->second;
    }

    void put(const K &key, const V &value) {
        auto it = map_.find(key);
        if (it != map_.end()) {
            it->second->second = value;
            order_.splice(order_.begin(), order_, it->second);
            return;
        }
        if ((int)order_.size() >= capacity_) {
            auto last = order_.back();
            map_.erase(last.first);
            order_.pop_back();
        }
        order_.emplace_front(key, value);
        map_[key] = order_.begin();
    }

private:
    int capacity_;
    std::list<std::pair<K, V>> order_;
    std::unordered_map<K, typename std::list<std::pair<K, V>>::iterator> map_;
};

#endif
