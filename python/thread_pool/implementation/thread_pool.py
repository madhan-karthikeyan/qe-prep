from __future__ import annotations

import threading
from collections.abc import Callable
from concurrent.futures import Future
from queue import Queue as ThreadQueue
from typing import Any


class ThreadPool:
    def __init__(self, num_workers: int, max_queue_size: int = 0) -> None:
        if num_workers < 1:
            raise ValueError("num_workers must be >= 1")
        self._num_workers = num_workers
        _task_t = tuple[Callable[..., Any], tuple[Any, ...], dict[str, Any], Future[Any]]
        self._task_queue: ThreadQueue[_task_t | None] = ThreadQueue(maxsize=max_queue_size)
        self._shutdown = threading.Event()
        self._workers: list[threading.Thread] = []
        self._start_workers()

    def _start_workers(self) -> None:
        for _ in range(self._num_workers):
            t = threading.Thread(target=self._worker_loop, daemon=True)
            t.start()
            self._workers.append(t)

    def _worker_loop(self) -> None:
        while not self._shutdown.is_set():
            try:
                task = self._task_queue.get(timeout=0.1)
            except Exception:
                continue
            if task is None:
                break
            fn, args, kwargs, future = task
            try:
                result = fn(*args, **kwargs)
                future.set_result(result)
            except BaseException as exc:
                future.set_exception(exc)

    def submit(
        self,
        fn: Callable[..., Any],
        *args: Any,
        **kwargs: Any,
    ) -> Future[Any]:
        if self._shutdown.is_set():
            raise RuntimeError("ThreadPool is shut down")
        future: Future[Any] = Future()
        self._task_queue.put((fn, args, kwargs, future))
        return future

    def shutdown(self, wait: bool = True) -> None:
        self._shutdown.set()
        if wait:
            for t in self._workers:
                t.join()

    @property
    def num_workers(self) -> int:
        return self._num_workers
