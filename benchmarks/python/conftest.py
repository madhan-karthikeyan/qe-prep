def pytest_benchmark_scale_unit(config, unit, benchmarks, best, worst, sort):
    if unit == "seconds":
        prefix = "m"
        scale = 0.001
    elif unit == "operations":
        prefix = "K"
        scale = 0.001
    else:
        prefix = ""
        scale = 1.0
    return prefix, scale
