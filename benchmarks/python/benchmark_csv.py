import io
import tempfile

import pytest
from file_processing.implementation.csv_parser import parse_csv


def _make_csv(num_rows: int, num_cols: int = 10) -> io.StringIO:
    header = ",".join(f"col_{i}" for i in range(num_cols))
    rows = [",".join(str(i * j) for j in range(num_cols)) for i in range(num_rows)]
    return io.StringIO(header + "\n" + "\n".join(rows))


@pytest.mark.parametrize("num_rows", [100, 1000, 10000])
def test_csv_parse_small(benchmark, num_rows):
    data = _make_csv(num_rows)
    benchmark(lambda: list(parse_csv(data)))


@pytest.mark.parametrize("num_rows", [100, 1000, 10000])
def test_csv_parse_no_header(benchmark, num_rows):
    data = _make_csv(num_rows)
    benchmark(lambda: list(parse_csv(data, has_header=False)))


@pytest.mark.parametrize("delimiter", [",", ";", "\t"])
def test_csv_delimiters(benchmark, delimiter):
    header = delimiter.join(f"col_{i}" for i in range(5))
    rows = [delimiter.join(str(i * j) for j in range(5)) for i in range(500)]
    data = io.StringIO(header + "\n" + "\n".join(rows))
    benchmark(lambda: list(parse_csv(data, delimiter=delimiter)))


def test_csv_large_file(benchmark):
    num_rows = 50000
    num_cols = 5
    header = ",".join(f"col_{i}" for i in range(num_cols))
    rows = [",".join(str(i * j) for j in range(num_cols)) for i in range(num_rows)]
    data = io.StringIO(header + "\n" + "\n".join(rows))
    benchmark(lambda: list(parse_csv(data)))
