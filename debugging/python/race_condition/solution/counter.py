import threading

class Counter:
    def __init__(self):
        self.value = 0
        self.lock = threading.Lock()

    def increment(self):
        with self.lock:
            self.value += 1

    def get_value(self):
        return self.value


counter = Counter()

def worker():
    for _ in range(10000):
        counter.increment()

threads = [threading.Thread(target=worker) for _ in range(10)]

for t in threads:
    t.start()
for t in threads:
    t.join()

print(f"Final counter value: {counter.get_value()}")
print(f"Expected: {10000 * 10}")
