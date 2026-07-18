#ifndef PRODUCER_CONSUMER_HPP
#define PRODUCER_CONSUMER_HPP

#include <queue>
#include <mutex>
#include <condition_variable>

template<typename T>
class BlockingQueue {
public:
    BlockingQueue(int capacity) : capacity_(capacity) {}

    void put(const T &val) {
        std::unique_lock<std::mutex> lock(mtx_);
        not_full_.wait(lock, [this]{ return (int)queue_.size() < capacity_; });
        queue_.push(val);
        not_empty_.notify_one();
    }

    T get() {
        std::unique_lock<std::mutex> lock(mtx_);
        not_empty_.wait(lock, [this]{ return !queue_.empty(); });
        T val = queue_.front();
        queue_.pop();
        not_full_.notify_one();
        return val;
    }

private:
    int capacity_;
    std::queue<T> queue_;
    std::mutex mtx_;
    std::condition_variable not_full_;
    std::condition_variable not_empty_;
};

#endif
