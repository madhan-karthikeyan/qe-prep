from __future__ import annotations

import threading

from logging_toolkit.implementation.logger import Logger, RotatingFileHandler


def test_10k_log_entries_multiple_threads(tmp_path) -> None:
    log_file = tmp_path / "stress.log"
    handler = RotatingFileHandler(str(log_file), max_bytes=50 * 1024 * 1024)
    logger = Logger()
    logger.add_handler(handler)

    n_threads = 8
    entries_per_thread = 1250
    total = n_threads * entries_per_thread
    barrier = threading.Barrier(n_threads)

    def worker(tid: int) -> None:
        barrier.wait()
        for i in range(entries_per_thread):
            logger.info(f"thread-{tid} entry-{i}")

    threads = [threading.Thread(target=worker, args=(tid,)) for tid in range(n_threads)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    lines = log_file.read_text().strip().split("\n")
    assert len(lines) == total

    seen = set()
    for line in lines:
        assert line not in seen, f"Duplicate line: {line}"
        seen.add(line)
