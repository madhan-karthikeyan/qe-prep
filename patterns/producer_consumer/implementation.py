import queue
import threading
from typing import Any, Callable, Optional

_POISON = object()


class ProducerConsumer:
    """Orchestrate producers and consumers communicating via a thread-safe queue.

    Uses the poison pill pattern for graceful shutdown.

    Args:
        num_producers: Number of producer threads.
        num_consumers: Number of consumer threads.
        maxsize: Maximum queue size (0 = unbounded).
    """

    def __init__(
        self,
        num_producers: int = 1,
        num_consumers: int = 1,
        maxsize: int = 0,
    ) -> None:
        self.num_producers = num_producers
        self.num_consumers = num_consumers
        self._queue: queue.Queue = queue.Queue(maxsize=maxsize)
        self._producers: list[threading.Thread] = []
        self._consumers: list[threading.Thread] = []

    def run(
        self,
        producer_fn: Callable[[queue.Queue], None],
        consumer_fn: Callable[[Any], None],
    ) -> None:
        """Start producers and consumers, then wait for completion.

        *producer_fn* receives the shared queue and should call queue.put(item).
        *consumer_fn* receives a single item.
        When *producer_fn* returns, a poison pill is sent per consumer.
        """
        self._producers = [
            threading.Thread(
                target=self._producer_wrapper,
                args=(producer_fn, i),
                daemon=True,
            )
            for i in range(self.num_producers)
        ]
        self._consumers = [
            threading.Thread(
                target=self._consumer_wrapper,
                args=(consumer_fn,),
                daemon=True,
            )
            for _ in range(self.num_consumers)
        ]

        for t in self._producers + self._consumers:
            t.start()

        for t in self._producers:
            t.join()

        for _ in range(self.num_consumers):
            self._queue.put(_POISON)

        for t in self._consumers:
            t.join()

    def _producer_wrapper(
        self, producer_fn: Callable[[queue.Queue], None], idx: int
    ) -> None:
        try:
            producer_fn(self._queue)
        except Exception:
            pass

    def _consumer_wrapper(
        self, consumer_fn: Callable[[Any], None]
    ) -> None:
        while True:
            item = self._queue.get()
            if item is _POISON:
                self._queue.task_done()
                break
            try:
                consumer_fn(item)
            except Exception:
                pass
            finally:
                self._queue.task_done()
