import cProfile
import pstats
from lru_cache.implementation.lru_cache import LRUCache


def profile_lru():
    cache = LRUCache[int, int](1000)
    for i in range(100000):
        cache.put(i, i)
        if i % 2 == 0:
            cache.get(i // 2)
        if i > 500:
            cache.get(i - 500)


if __name__ == "__main__":
    prof_file = "lru.prof"

    cProfile.run("profile_lru()", prof_file)

    p = pstats.Stats(prof_file)
    p.sort_stats("cumtime").print_stats(20)
    print("\n--- By ncalls ---")
    p.sort_stats("ncalls").print_stats(20)
    print("\n--- By time ---")
    p.sort_stats("time").print_stats(20)
