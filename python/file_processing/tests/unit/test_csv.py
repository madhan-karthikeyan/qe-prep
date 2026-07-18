from __future__ import annotations

import io

from file_processing.implementation.csv_parser import parse_csv


class TestCSVParser:
    def test_simple_csv(self) -> None:
        text = "a,b,c\n1,2,3\n4,5,6\n"
        rows = list(parse_csv(io.StringIO(text)))
        assert rows == [["a", "b", "c"], ["1", "2", "3"], ["4", "5", "6"]]

    def test_quoted_fields(self) -> None:
        text = 'name,"hello, world",age\nfoo,"bar, baz",42\n'
        rows = list(parse_csv(io.StringIO(text)))
        assert rows == [["name", "hello, world", "age"], ["foo", "bar, baz", "42"]]

    def test_escaped_quotes(self) -> None:
        text = 'a,"say ""hello""",b\n'
        rows = list(parse_csv(io.StringIO(text)))
        assert rows == [["a", 'say "hello"', "b"]]

    def test_custom_delimiter(self) -> None:
        text = "a|b|c\n1|2|3\n"
        rows = list(parse_csv(io.StringIO(text), delimiter="|"))
        assert rows == [["a", "b", "c"], ["1", "2", "3"]]

    def test_no_header(self) -> None:
        text = "1,2,3\n4,5,6\n"
        rows = list(parse_csv(io.StringIO(text), has_header=False))
        assert rows == [["1", "2", "3"], ["4", "5", "6"]]

    def test_empty_fields(self) -> None:
        text = "a,,c\n,2,\n"
        rows = list(parse_csv(io.StringIO(text)))
        assert rows == [["a", "", "c"], ["", "2", ""]]

    def test_trailing_comma(self) -> None:
        text = "a,b,\n1,2,\n"
        rows = list(parse_csv(io.StringIO(text)))
        assert rows == [["a", "b", ""], ["1", "2", ""]]

    def test_empty_file(self) -> None:
        rows = list(parse_csv(io.StringIO("")))
        assert rows == []

    def test_header_only_file(self) -> None:
        rows = list(parse_csv(io.StringIO("a,b,c\n")))
        assert rows == [["a", "b", "c"]]

    def test_streaming(self) -> None:
        text = "a,b\n1,2\n3,4\n5,6\n"
        it = parse_csv(io.StringIO(text))
        assert next(it) == ["a", "b"]
        assert next(it) == ["1", "2"]
        assert next(it) == ["3", "4"]
        assert next(it) == ["5", "6"]
