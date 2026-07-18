import threading

import pytest

from pub_sub.implementation import PubSub


def test_message_delivery_to_all_subscribers():
    pubsub = PubSub(max_workers=2)
    received: list[str] = []
    lock = threading.Lock()

    def cb1(msg: str):
        with lock:
            received.append(f"1:{msg}")

    def cb2(msg: str):
        with lock:
            received.append(f"2:{msg}")

    pubsub.subscribe("events", cb1)
    pubsub.subscribe("events", cb2)
    pubsub.publish("events", "hello")
    pubsub.shutdown()

    assert sorted(received) == ["1:hello", "2:hello"]


def test_unsubscribe():
    pubsub = PubSub(max_workers=2)
    received: list[str] = []
    lock = threading.Lock()

    def cb(msg: str):
        with lock:
            received.append(msg)

    sid = pubsub.subscribe("events", cb)
    pubsub.publish("events", "first")
    pubsub.unsubscribe(sid)
    pubsub.publish("events", "second")
    pubsub.shutdown()

    assert received == ["first"]


def test_no_subscriber():
    pubsub = PubSub()
    pubsub.publish("nonexistent", "data")
    pubsub.shutdown()


def test_multiple_topics():
    pubsub = PubSub(max_workers=2)
    events: list[str] = []
    notifications: list[str] = []
    lock = threading.Lock()

    def on_event(msg: str):
        with lock:
            events.append(msg)

    def on_notification(msg: str):
        with lock:
            notifications.append(msg)

    pubsub.subscribe("events", on_event)
    pubsub.subscribe("notifications", on_notification)

    pubsub.publish("events", "e1")
    pubsub.publish("notifications", "n1")
    pubsub.publish("events", "e2")
    pubsub.shutdown()

    assert sorted(events) == ["e1", "e2"]
    assert notifications == ["n1"]


def test_subscribe_same_callback_multiple_times():
    pubsub = PubSub(max_workers=2)
    received: list[str] = []
    lock = threading.Lock()

    def cb(msg: str):
        with lock:
            received.append(msg)

    pubsub.subscribe("events", cb)
    pubsub.subscribe("events", cb)
    pubsub.publish("events", "hi")
    pubsub.shutdown()

    assert received.count("hi") == 2
