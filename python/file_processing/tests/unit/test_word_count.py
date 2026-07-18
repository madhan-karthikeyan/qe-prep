from __future__ import annotations

import io

from file_processing.implementation.word_count import count_stream


class TestWordCount:
    def test_empty(self) -> None:
        lines, words, chars = count_stream(io.StringIO(""))
        assert (lines, words, chars) == (0, 0, 0)

    def test_single_line(self) -> None:
        lines, words, chars = count_stream(io.StringIO("hello world\n"))
        assert (lines, words, chars) == (1, 2, 12)

    def test_multiple_lines(self) -> None:
        text = "one two three\nfour five\nsix\n"
        lines, words, chars = count_stream(io.StringIO(text))
        assert (lines, words, chars) == (3, 6, 28)

    def test_unicode_chars(self) -> None:
        text = "héllo wörld\n"
        lines, words, chars = count_stream(io.StringIO(text))
        assert chars == 12

    def test_leading_trailing_spaces(self) -> None:
        text = "  hello   world  \n"
        lines, words, chars = count_stream(io.StringIO(text))
        assert (lines, words, chars) == (1, 2, 18)

    def test_blank_lines(self) -> None:
        text = "foo\n\n\nbar\n"
        lines, words, chars = count_stream(io.StringIO(text))
        assert (lines, words, chars) == (4, 2, 10)
