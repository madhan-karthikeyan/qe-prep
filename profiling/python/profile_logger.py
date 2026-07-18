import cProfile
import pstats
import os
from logging_toolkit.implementation.logger import Logger, RotatingFileHandler


def profile_logger():
    logger = Logger()
    handler = RotatingFileHandler(os.devnull)
    logger.add_handler(handler)
    for i in range(10000):
        logger.info(f"line {i}")


if __name__ == "__main__":
    prof_file = "logger.prof"

    cProfile.run("profile_logger()", prof_file)

    p = pstats.Stats(prof_file)
    p.sort_stats("cumtime").print_stats(20)
    print("\n--- By ncalls ---")
    p.sort_stats("ncalls").print_stats(20)
    print("\n--- By time ---")
    p.sort_stats("time").print_stats(20)
