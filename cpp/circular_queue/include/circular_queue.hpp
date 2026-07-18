#ifndef CIRCULAR_QUEUE_HPP
#define CIRCULAR_QUEUE_HPP

#include <vector>
#include <stdexcept>

template<typename T>
class CircularQueue {
public:
    CircularQueue(int capacity) : buffer_(capacity), head_(0), tail_(0), full_(false) {}

    void enqueue(const T &val) {
        if (full_) throw std::overflow_error("queue full");
        buffer_[tail_] = val;
        tail_ = (tail_ + 1) % buffer_.size();
        full_ = (tail_ == head_);
    }

    T dequeue() {
        if (is_empty()) throw std::underflow_error("queue empty");
        T val = buffer_[head_];
        head_ = (head_ + 1) % buffer_.size();
        full_ = false;
        return val;
    }

    bool is_full() const { return full_; }
    bool is_empty() const { return !full_ && head_ == tail_; }

private:
    std::vector<T> buffer_;
    int head_;
    int tail_;
    bool full_;
};

#endif
