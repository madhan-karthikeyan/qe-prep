import queue
import threading
from typing import Any, Callable, Optional

_POISON = object()


class WorkerPool:
    """Fixed-size worker pool distributing tasks via a shared queue.

    Args:
        num_workers: Number of worker threads.
        maxsize: Maximum queued tasks (0 = unbounded).
    """

    def __init__(self, num_workers: int = 4, maxsize: int = 0) -> None:
        self._num_workers = num_workers
        self._results: queue.Queue = queue.Queue()
        self._task_queue: queue.Queue = queue.Queue(maxsize=maxsize)
        self._workers: list[threading.Thread] = []
        self._shutdown = False
        self._lock = threading.Lock()

    def _worker_loop(self) -> None:
        while True:
            try:
                task = self._task_queue.get(timeout=0.1)
            except queue.Empty:
                with self._lock:
                    if self._shutdown:
                        return
                continue

            if task is _POISON:
                self._task_queue.task_done()
                break

            func, args, kwargs, result_callback = task
            try:
                result = func(*args, **kwargs)
                if result_callback:
                    result_callback(result)
            except Exception as exc:
                if result_callback:
                    result_callback(exc)
            finally:
                self._task_queue.task_done()

    def start(self) -> None:
        """Start the worker threads."""
        self._workers = [
            threading.Thread(target=self._worker_loop, daemon=True)
            for _ in range(self._num_workers)
        ]
        for w in self._workers:
            w.start()

    def submit(
        self,
        func: Callable[..., Any],
        *args,
        result_callback: Optional[Callable[[Any], None]] = None,
        **kwargs,
    ) -> None:
        """Submit a task for async execution (non-blocking)."""
        if self._shutdown:
            raise RuntimeError("WorkerPool is shut down")
        self._task_queue.put((func, args, kwargs, result_callback))

    def map(
        self,
        func: Callable[[Any], Any],
        iterable: list[Any],
        result_callback: Optional[Callable[[Any], None]] = None,
    ) -> None:
        """Distribute items from *iterable* to workers via *func*."""
        for item in iterable:
            self.submit(func, item, result_callback=result_callback)

    def shutdown(self, wait: bool = True) -> None:
        """Shut down the pool. Optionally wait for pending tasks."""
        with self._lock:
            self._shutdown = True

        for _ in range(self._num_workers):
            self._task_queue.put(_POISON)

        if wait:
            self._task_queue.join()
            for w in self._workers:
                w.join()
