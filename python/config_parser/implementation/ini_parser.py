

class IniParser:
    def __init__(self, source: str | None = None) -> None:
        self._sections: dict[str, dict[str, str]] = {"DEFAULT": {}}
        if source is not None:
            self.parse(source)

    def parse(self, source: str) -> dict[str, dict[str, str]]:
        self._sections = {"DEFAULT": {}}
        current_section = "DEFAULT"
        for line in source.splitlines():
            stripped = line.strip()
            if not stripped or stripped.startswith(";") or stripped.startswith("#"):
                continue
            if stripped.startswith("[") and stripped.endswith("]"):
                current_section = stripped[1:-1].strip()
                if current_section not in self._sections:
                    self._sections[current_section] = {}
                continue
            if "=" in stripped or ":" in stripped:
                sep = "=" if "=" in stripped else ":"
                key, _, value = stripped.partition(sep)
                key = key.strip()
                value = value.strip()
                if key:
                    if value.startswith('"') and value.endswith('"'):
                        value = value[1:-1]
                    elif value.startswith("'") and value.endswith("'"):
                        value = value[1:-1]
                    self._sections[current_section][key] = value

        return dict(self._sections)

    def parse_file(self, path: str) -> dict[str, dict[str, str]]:
        with open(path) as f:
            return self.parse(f.read())

    def get(self, section: str, key: str, default: str | None = None) -> str | None:
        if section in self._sections and key in self._sections[section]:
            return self._sections[section][key]
        if key in self._sections.get("DEFAULT", {}):
            return self._sections["DEFAULT"][key]
        return default

    def getint(self, section: str, key: str, default: int | None = None) -> int | None:
        value = self.get(section, key)
        if value is None:
            return default
        try:
            return int(value)
        except ValueError:
            raise ValueError(f"Cannot convert {key!r}={value!r} to int")

    def getfloat(self, section: str, key: str, default: float | None = None) -> float | None:
        value = self.get(section, key)
        if value is None:
            return default
        try:
            return float(value)
        except ValueError:
            raise ValueError(f"Cannot convert {key!r}={value!r} to float")

    def getbool(self, section: str, key: str, default: bool | None = None) -> bool | None:
        value = self.get(section, key)
        if value is None:
            return default
        lower = value.lower()
        if lower in ("true", "yes", "1", "on"):
            return True
        if lower in ("false", "no", "0", "off"):
            return False
        raise ValueError(f"Cannot convert {key!r}={value!r} to bool")

    def sections(self) -> list[str]:
        return [s for s in self._sections if s != "DEFAULT"]

    def __getitem__(self, section: str) -> dict[str, str]:
        if section not in self._sections:
            raise KeyError(f"Section {section!r} not found")
        return dict(self._sections[section])

    def __repr__(self) -> str:
        return f"IniParser({dict(self._sections)!r})"
