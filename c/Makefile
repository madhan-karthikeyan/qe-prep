CC = gcc
CFLAGS = -Wall -Wextra -std=c11

all: test_logger test_lru test_stack_queue test_tcp_echo test_word_count test_circular_queue test_producer_consumer test_config_parser

test_logger: logging_toolkit/src/logger.c logging_toolkit/tests/test_logger.c
	$(CC) $(CFLAGS) -I logging_toolkit/include -o logging_toolkit/tests/test_logger $^

test_lru: lru_cache/src/lru_cache.c lru_cache/tests/test_lru_cache.c
	$(CC) $(CFLAGS) -I lru_cache/include -o lru_cache/tests/test_lru $^

test_stack_queue: linear_data_structures/src/stack_queue.c linear_data_structures/tests/test_stack_queue.c
	$(CC) $(CFLAGS) -I linear_data_structures/include -o linear_data_structures/tests/test_stack_queue $^

test_tcp_echo: tcp_echo/src/tcp_echo.c tcp_echo/tests/test_tcp_echo.c
	$(CC) $(CFLAGS) -I tcp_echo/include -o tcp_echo/tests/test_tcp_echo $^

test_word_count: word_count/src/word_count.c word_count/tests/test_word_count.c
	$(CC) $(CFLAGS) -I word_count/include -o word_count/tests/test_word_count $^

test_circular_queue: circular_queue/src/circular_queue.c circular_queue/tests/test_circular_queue.c
	$(CC) $(CFLAGS) -I circular_queue/include -o circular_queue/tests/test_circular_queue $^

test_producer_consumer: producer_consumer/src/producer_consumer.c producer_consumer/tests/test_producer_consumer.c
	$(CC) $(CFLAGS) -I producer_consumer/include -o producer_consumer/tests/test_producer_consumer $^ -lpthread

test_config_parser: config_parser/src/config_parser.c config_parser/tests/test_config_parser.c
	$(CC) $(CFLAGS) -I config_parser/include -o config_parser/tests/test_config_parser $^

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
