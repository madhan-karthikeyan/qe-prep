import os
import tempfile

import pytest
from logging_toolkit.implementation.logger import Logger, RotatingFileHandler


def test_logger_speed(benchmark):
    logger = Logger()
    handler = RotatingFileHandler(os.devnull)
    logger.add_handler(handler)
    benchmark(logger.info, "test message 42")


def test_logger_with_format(benchmark):
    logger = Logger(format_template="{timestamp} {level} {message}")
    handler = RotatingFileHandler(os.devnull)
    logger.add_handler(handler)
    benchmark(logger.info, "formatted message")


def test_logger_multiple_handlers(benchmark):
    logger = Logger()
    for _ in range(3):
        logger.add_handler(RotatingFileHandler(os.devnull))
    benchmark(logger.info, "multi handler message")


def test_logger_filtered(benchmark):
    from logging_toolkit.implementation.logger import LogFilter, LogLevel

    logger = Logger(filter=LogFilter(min_level=LogLevel.ERROR))
    handler = RotatingFileHandler(os.devnull)
    logger.add_handler(handler)
    benchmark(logger.info, "this should be filtered out")
