import itertools
import threading
from concurrent.futures import ThreadPoolExecutor
from typing import Any, Callable

SubscriptionId = int
Callback = Callable[[Any], None]


class PubSub:
    """Publish-subscribe with async delivery via thread pool.

    Args:
        max_workers: Max threads for concurrent message delivery.
    """

    def __init__(self, max_workers: int = 4) -> None:
        self._subscriptions: dict[str, dict[SubscriptionId, Callback]] = {}
        self._id_counter = itertools.count(1)
        self._lock = threading.Lock()
        self._executor = ThreadPoolExecutor(max_workers=max_workers)

    def subscribe(self, topic: str, callback: Callback) -> SubscriptionId:
        """Subscribe *callback* to *topic*. Returns a subscription ID."""
        with self._lock:
            sub_id = next(self._id_counter)
            if topic not in self._subscriptions:
                self._subscriptions[topic] = {}
            self._subscriptions[topic][sub_id] = callback
        return sub_id

    def unsubscribe(self, subscription_id: SubscriptionId) -> None:
        """Remove a subscription by ID."""
        with self._lock:
            for topic in list(self._subscriptions):
                self._subscriptions[topic].pop(subscription_id, None)

    def publish(self, topic: str, message: Any) -> None:
        """Deliver *message* to all subscribers of *topic* (async)."""
        with self._lock:
            callbacks = list(self._subscriptions.get(topic, {}).values())
        for cb in callbacks:
            self._executor.submit(cb, message)

    def shutdown(self) -> None:
        """Shut down the executor (waits for pending deliveries)."""
        self._executor.shutdown(wait=True)
