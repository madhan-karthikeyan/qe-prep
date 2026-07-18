CXX = g++
CXXFLAGS = -Wall -Wextra -std=c++17

all: test_logger test_lru test_stack_queue test_tcp_echo test_word_count test_circular_queue test_producer_consumer test_config_parser

test_logger: logging_toolkit/src/logger.cpp logging_toolkit/tests/test_logger.cpp
	$(CXX) $(CXXFLAGS) -I logging_toolkit/include -o logging_toolkit/tests/test_logger $^

test_lru: lru_cache/src/lru_cache.cpp lru_cache/tests/test_lru_cache.cpp
	$(CXX) $(CXXFLAGS) -I lru_cache/include -o lru_cache/tests/test_lru $^

test_stack_queue: linear_data_structures/src/stack_queue.cpp linear_data_structures/tests/test_stack_queue.cpp
	$(CXX) $(CXXFLAGS) -I linear_data_structures/include -o linear_data_structures/tests/test_stack_queue $^

test_tcp_echo: tcp_echo/src/tcp_echo.cpp tcp_echo/tests/test_tcp_echo.cpp
	$(CXX) $(CXXFLAGS) -I tcp_echo/include -o tcp_echo/tests/test_tcp_echo $^

test_word_count: word_count/src/word_count.cpp word_count/tests/test_word_count.cpp
	$(CXX) $(CXXFLAGS) -I word_count/include -o word_count/tests/test_word_count $^

test_circular_queue: circular_queue/src/circular_queue.cpp circular_queue/tests/test_circular_queue.cpp
	$(CXX) $(CXXFLAGS) -I circular_queue/include -o circular_queue/tests/test_circular_queue $^

test_producer_consumer: producer_consumer/src/producer_consumer.cpp producer_consumer/tests/test_producer_consumer.cpp
	$(CXX) $(CXXFLAGS) -I producer_consumer/include -o producer_consumer/tests/test_producer_consumer $^ -lpthread

test_config_parser: config_parser/src/config_parser.cpp config_parser/tests/test_config_parser.cpp
	$(CXX) $(CXXFLAGS) -I config_parser/include -o config_parser/tests/test_config_parser $^

test: all
	./logging_toolkit/tests/test_logger && echo "PASS: logger" || echo "FAIL: logger"
	./lru_cache/tests/test_lru && echo "PASS: lru" || echo "FAIL: lru"
	./linear_data_structures/tests/test_stack_queue && echo "PASS: stack_queue" || echo "FAIL: stack_queue"
	./tcp_echo/tests/test_tcp_echo && echo "PASS: tcp_echo" || echo "FAIL: tcp_echo"
	./word_count/tests/test_word_count && echo "PASS: word_count" || echo "FAIL: word_count"
	./circular_queue/tests/test_circular_queue && echo "PASS: circular_queue" || echo "FAIL: circular_queue"
	./producer_consumer/tests/test_producer_consumer && echo "PASS: producer_consumer" || echo "FAIL: producer_consumer"
	./config_parser/tests/test_config_parser && echo "PASS: config_parser" || echo "FAIL: config_parser"

clean:
	rm -f */tests/test_*

.PHONY: all test clean
