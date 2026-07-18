#ifndef STACK_HPP
#define STACK_HPP

#include <vector>
#include <stdexcept>

template<typename T>
class Stack {
public:
    void push(const T &val) { data_.push_back(val); }
    T pop() {
        if (data_.empty()) throw std::out_of_range("stack empty");
        T val = data_.back();
        data_.pop_back();
        return val;
    }
    T &peek() {
        if (data_.empty()) throw std::out_of_range("stack empty");
        return data_.back();
    }
    bool is_empty() const { return data_.empty(); }

private:
    std::vector<T> data_;
};

#endif
