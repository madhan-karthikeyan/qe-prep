from __future__ import annotations

import threading

from logging_toolkit.implementation.logger import (
    Logger,
    RotatingFileHandler,
)


def test_file_rotation(tmp_path) -> None:
    log_file = tmp_path / "test.log"
    handler = RotatingFileHandler(str(log_file), max_bytes=100, backup_count=2)
    logger = Logger()
    logger.add_handler(handler)

    for i in range(200):
        logger.info(f"line {i}")

    assert log_file.exists() or (tmp_path / "test.log.1").exists()


def test_multi_threaded_writes(tmp_path) -> None:
    log_file = tmp_path / "mt.log"
    handler = RotatingFileHandler(str(log_file), max_bytes=10 * 1024 * 1024)
    logger = Logger()
    logger.add_handler(handler)

    n_threads = 10
    lines_per_thread = 100
    barrier = threading.Barrier(n_threads)

    def writer(thread_id: int) -> None:
        barrier.wait()
        for i in range(lines_per_thread):
            logger.info(f"thread-{thread_id} line-{i}")

    threads = [threading.Thread(target=writer, args=(tid,)) for tid in range(n_threads)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    content = log_file.read_text()
    expected_lines = n_threads * lines_per_thread
    assert len(content.strip().split("\n")) == expected_lines


def test_log_file_restored_after_rotation(tmp_path) -> None:
    log_file = tmp_path / "rotate.log"
    handler = RotatingFileHandler(str(log_file), max_bytes=50, backup_count=3)
    logger = Logger()
    logger.add_handler(handler)

    for i in range(100):
        logger.info(f"msg {i}")

    files = list(tmp_path.glob("rotate.log*"))
    assert len(files) > 1
