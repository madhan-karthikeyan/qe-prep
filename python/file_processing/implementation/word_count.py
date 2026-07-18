from __future__ import annotations

import sys
from typing import TextIO

WORD_COUNT_DISPLAY = "{lines:>7} {words:>7} {chars:>7} {name}"


def count_stream(stream: TextIO) -> tuple[int, int, int]:
    lines = 0
    words = 0
    chars = 0
    for line in stream:
        lines += 1
        words += len(line.split())
        chars += len(line)
    return lines, words, chars


def word_count(files: list[str]) -> None:
    total_lines = 0
    total_words = 0
    total_chars = 0

    if not files:
        lines, words, chars = count_stream(sys.stdin)
        print(WORD_COUNT_DISPLAY.format(lines=lines, words=words, chars=chars, name=""))
        return

    for path in files:
        try:
            with open(path, encoding="utf-8") as f:
                lines, words, chars = count_stream(f)
        except FileNotFoundError as e:
            print(f"wc: {e.filename}: No such file or directory", file=sys.stderr)
            continue
        except IsADirectoryError as e:
            print(f"wc: {e.filename}: Is a directory", file=sys.stderr)
            continue

        total_lines += lines
        total_words += words
        total_chars += chars
        print(WORD_COUNT_DISPLAY.format(lines=lines, words=words, chars=chars, name=path))

    if len(files) > 1:
        print(WORD_COUNT_DISPLAY.format(
            lines=total_lines, words=total_words, chars=total_chars, name="total",
        ))


def main() -> None:
    word_count(sys.argv[1:])


if __name__ == "__main__":
    main()
