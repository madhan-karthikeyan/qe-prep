from __future__ import annotations

import csv as csv_module
from collections.abc import Iterator
from typing import TextIO


class CSVParseError(Exception):
    pass


def parse_csv(
    file: TextIO,
    delimiter: str = ",",
    has_header: bool = True,
) -> Iterator[list[str]]:
    reader = csv_module.reader(
        file,
        delimiter=delimiter,
        quoting=csv_module.QUOTE_MINIMAL,
        skipinitialspace=False,
    )
    if has_header:
        try:
            header = next(reader)
        except StopIteration:
            return
        yield header
    for row in reader:
        yield [str(field) if field is not None else "" for field in row]
