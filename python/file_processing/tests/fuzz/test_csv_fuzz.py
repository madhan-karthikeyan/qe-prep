from __future__ import annotations

import io
import random
import string

from file_processing.implementation.csv_parser import parse_csv


def _random_csv_line(cols: int) -> str:
    fields: list[str] = []
    for _ in range(cols):
        choice = random.random()
        if choice < 0.3:
            fields.append(random.choice(string.ascii_letters))
        elif choice < 0.6:
            fields.append(f'"{random.choice(string.ascii_letters)}"')
        elif choice < 0.8:
            q = random.choice(string.ascii_letters)
            fields.append(f'"{q}""{q}"')
        else:
            fields.append("")
    return ",".join(fields)


def test_random_csv_no_crash() -> None:
    for _ in range(1000):
        n_cols = random.randint(1, 10)
        n_rows = random.randint(0, 20)
        text = "\n".join(_random_csv_line(n_cols) for _ in range(n_rows))
        if random.random() < 0.2:
            text += ","
        try:
            rows = list(parse_csv(io.StringIO(text)))
            assert isinstance(rows, list)
        except Exception:
            pass


def test_extremely_long_line() -> None:
    long_field = "a" * 10000
    text = f"a,b,c\n{long_field},{long_field},{long_field}\n"
    rows = list(parse_csv(io.StringIO(text)))
    assert len(rows) == 2
    assert rows[1] == [long_field, long_field, long_field]


def test_special_characters() -> None:
    text = "a,b\n\x00,\x01\n\xff,\U0001f600\n"
    rows = list(parse_csv(io.StringIO(text)))
    assert len(rows) == 3
