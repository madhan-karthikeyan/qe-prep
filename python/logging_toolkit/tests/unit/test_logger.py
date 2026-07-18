from __future__ import annotations

import io

from logging_toolkit.implementation.logger import LogFilter, Logger, LogLevel, LogRecord


class TestLogRecord:
    def test_format_default(self) -> None:
        record = LogRecord(timestamp=0.0, level=LogLevel.INFO, message="hello")
        assert "1970-01-01 00:00:00 INFO hello" == record.format()

    def test_format_custom(self) -> None:
        record = LogRecord(timestamp=0.0, level=LogLevel.ERROR, message="fail")
        result = record.format("{level}: {message}")
        assert "ERROR: fail" == result


class TestLogFilter:
    def test_min_level_filters_lower(self) -> None:
        filt = LogFilter(min_level=LogLevel.WARN)
        assert not filt.accepts(LogRecord(0.0, LogLevel.DEBUG, "msg"))
        assert not filt.accepts(LogRecord(0.0, LogLevel.INFO, "msg"))
        assert filt.accepts(LogRecord(0.0, LogLevel.WARN, "msg"))
        assert filt.accepts(LogRecord(0.0, LogLevel.ERROR, "msg"))

    def test_pattern_filter(self) -> None:
        filt = LogFilter(pattern=r"error")
        assert filt.accepts(LogRecord(0.0, LogLevel.INFO, "something error happened"))
        assert not filt.accepts(LogRecord(0.0, LogLevel.INFO, "all good"))

    def test_combined_filter(self) -> None:
        filt = LogFilter(min_level=LogLevel.WARN, pattern=r"fail")
        assert filt.accepts(LogRecord(0.0, LogLevel.ERROR, "task failed"))
        assert not filt.accepts(LogRecord(0.0, LogLevel.DEBUG, "task failed"))
        assert not filt.accepts(LogRecord(0.0, LogLevel.ERROR, "task succeeded"))


class TestLogger:
    def test_log_to_stdout(self) -> None:
        buf = io.StringIO()
        logger = Logger()
        logger.add_handler(buf)
        logger.info("hello world")
        output = buf.getvalue()
        assert "INFO hello world" in output

    def test_log_levels_respected(self) -> None:
        buf = io.StringIO()
        logger = Logger(filter=LogFilter(min_level=LogLevel.WARN))
        logger.add_handler(buf)
        logger.info("should not appear")
        logger.warn("should appear")
        assert "should appear" in buf.getvalue()
        assert "should not appear" not in buf.getvalue()

    def test_format_template(self) -> None:
        buf = io.StringIO()
        logger = Logger(format_template="{level}: {message}")
        logger.add_handler(buf)
        logger.error("boom")
        assert "ERROR: boom\n" == buf.getvalue()

    def test_multiple_handlers(self) -> None:
        buf1 = io.StringIO()
        buf2 = io.StringIO()
        logger = Logger()
        logger.add_handler(buf1)
        logger.add_handler(buf2)
        logger.info("broadcast")
        assert "broadcast" in buf1.getvalue()
        assert "broadcast" in buf2.getvalue()

    def test_remove_handler(self) -> None:
        buf = io.StringIO()
        logger = Logger()
        logger.add_handler(buf)
        logger.remove_handler(buf)
        logger.info("silent")
        assert "" == buf.getvalue()
