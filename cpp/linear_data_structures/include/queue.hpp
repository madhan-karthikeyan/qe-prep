#ifndef QUEUE_HPP
#define QUEUE_HPP

#include <deque>
#include <stdexcept>

template<typename T>
class Queue {
public:
    void enqueue(const T &val) { data_.push_back(val); }
    T dequeue() {
        if (data_.empty()) throw std::out_of_range("queue empty");
        T val = data_.front();
        data_.pop_front();
        return val;
    }
    T &peek() {
        if (data_.empty()) throw std::out_of_range("queue empty");
        return data_.front();
    }
    bool is_empty() const { return data_.empty(); }

private:
    std::deque<T> data_;
};

#endif
