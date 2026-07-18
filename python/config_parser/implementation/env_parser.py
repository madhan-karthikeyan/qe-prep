import os
import re


class EnvParser:
    def __init__(self, source: str | None = None) -> None:
        self._data: dict[str, str] = {}
        if source is not None:
            self.parse(source)

    def parse(self, source: str) -> dict[str, str]:
        self._data = {}
        for line in source.splitlines():
            line = line.strip()
            if not line or line.startswith("#"):
                continue
            key, _, value = line.partition("=")
            key = key.strip()
            if not key:
                continue
            value = value.strip()
            value = self._unquote(value)
            value = self._substitute(value)
            self._data[key] = value
        return dict(self._data)

    def parse_file(self, path: str) -> dict[str, str]:
        with open(path) as f:
            return self.parse(f.read())

    def _unquote(self, value: str) -> str:
        if len(value) >= 2 and value[0] == value[-1] and value[0] in ("'", '"'):
            return self._unescape(value[1:-1])
        return value

    def _unescape(self, s: str) -> str:
        result: list[str] = []
        i = 0
        while i < len(s):
            if s[i] == "\\" and i + 1 < len(s):
                n = s[i + 1]
                esc_map = {"n": "\n", "r": "\r", "t": "\t", "\\": "\\", '"': '"', "'": "'"}
                result.append(esc_map.get(n, n))
                i += 2
            else:
                result.append(s[i])
                i += 1
        return "".join(result)

    def _substitute(self, value: str) -> str:
        def replacer(m: re.Match[str]) -> str:
            var_name = m.group(1) or m.group(2)
            return self._data.get(var_name, os.environ.get(var_name, m.group(0)))

        pattern = r"\$(\w+)|(?:\$\{(\w+)\})"
        return re.sub(pattern, replacer, value)

    def get(self, key: str, default: str | None = None) -> str | None:
        return self._data.get(key, default)

    def __getitem__(self, key: str) -> str:
        return self._data[key]

    def __contains__(self, key: str) -> bool:
        return key in self._data

    def __repr__(self) -> str:
        return f"EnvParser({self._data!r})"
