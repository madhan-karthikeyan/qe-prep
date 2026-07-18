from __future__ import annotations

import os
import re
import threading
import time
from dataclasses import dataclass
from enum import IntEnum
from typing import IO


class LogLevel(IntEnum):
    DEBUG = 10
    INFO = 20
    WARN = 30
    ERROR = 40


@dataclass
class LogRecord:
    timestamp: float
    level: LogLevel
    message: str

    def format(self, fmt: str = "{timestamp} {level} {message}") -> str:
        return fmt.format(
            timestamp=time.strftime("%Y-%m-%d %H:%M:%S", time.gmtime(self.timestamp)),
            level=self.level.name,
            message=self.message,
        )


class LogFilter:
    def __init__(self, min_level: LogLevel | None = None, pattern: str | None = None) -> None:
        self.min_level = min_level
        self._regex = re.compile(pattern) if pattern is not None else None

    def accepts(self, record: LogRecord) -> bool:
        if self.min_level is not None and record.level < self.min_level:
            return False
        if self._regex is not None and not self._regex.search(record.message):
            return False
        return True


class RotatingFileHandler:
    def __init__(self, path: str, max_bytes: int = 10 * 1024 * 1024, backup_count: int = 5) -> None:
        self.path = path
        self.max_bytes = max_bytes
        self.backup_count = backup_count

    def write(self, text: str) -> None:
        if os.path.exists(self.path) and os.path.getsize(self.path) >= self.max_bytes:
            self._rotate()
        with open(self.path, "a", encoding="utf-8") as f:
            f.write(text)

    def _rotate(self) -> None:
        for i in range(self.backup_count - 1, 0, -1):
            src = f"{self.path}.{i}"
            dst = f"{self.path}.{i + 1}"
            if os.path.exists(src):
                os.replace(src, dst)
        if os.path.exists(self.path):
            os.replace(self.path, f"{self.path}.1")


class Logger:
    def __init__(
        self,
        format_template: str = "{timestamp} {level} {message}",
        filter: LogFilter | None = None,
    ) -> None:
        self._format_template = format_template
        self._filter = filter or LogFilter()
        self._handlers: list[IO[str] | RotatingFileHandler] = []
        self._lock = threading.Lock()

    def add_handler(self, handler: IO[str] | RotatingFileHandler) -> None:
        self._handlers.append(handler)

    def remove_handler(self, handler: IO[str] | RotatingFileHandler) -> None:
        self._handlers.remove(handler)

    def log(self, level: LogLevel, message: str) -> None:
        record = LogRecord(timestamp=time.time(), level=level, message=message)
        if not self._filter.accepts(record):
            return
        line = record.format(self._format_template) + "\n"
        with self._lock:
            for handler in self._handlers:
                handler.write(line)

    def debug(self, message: str) -> None:
        self.log(LogLevel.DEBUG, message)

    def info(self, message: str) -> None:
        self.log(LogLevel.INFO, message)

    def warn(self, message: str) -> None:
        self.log(LogLevel.WARN, message)

    def error(self, message: str) -> None:
        self.log(LogLevel.ERROR, message)
