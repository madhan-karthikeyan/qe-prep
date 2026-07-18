from dataclasses import dataclass


@dataclass
class ParsedURL:
    scheme: str
    userinfo: str | None
    host: str
    port: int | None
    path: str
    query: str | None
    fragment: str | None

    def reconstruct(self) -> str:
        parts: list[str] = []
        parts.append(f"{self.scheme}://")
        if self.userinfo is not None:
            parts.append(f"{self.userinfo}@")
        if ":" in self.host:
            parts.append(f"[{self.host}]")
        else:
            parts.append(self.host)
        default_port = _default_ports.get(self.scheme)
        if self.port is not None and self.port != default_port:
            parts.append(f":{self.port}")
        parts.append(self.path)
        if self.query is not None:
            parts.append(f"?{self.query}")
        if self.fragment is not None:
            parts.append(f"#{self.fragment}")
        return "".join(parts)

    def validate(self) -> bool:
        if not self.scheme:
            return False
        if not self.host:
            return False
        if self.port is not None and not (0 < self.port <= 65535):
            return False
        return True


_default_ports: dict[str, int] = {
    "http": 80,
    "https": 443,
    "ftp": 21,
    "ssh": 22,
    "smtp": 25,
    "ws": 80,
    "wss": 443,
}


class URLParser:
    @staticmethod
    def parse(url: str) -> ParsedURL:
        if not url:
            raise ValueError("Empty URL")

        fragment: str | None = None
        if "#" in url:
            url, fragment = url.rsplit("#", 1)
            if fragment == "":
                fragment = None

        query: str | None = None
        if "?" in url:
            url, query = url.rsplit("?", 1)
            if query == "":
                query = None

        scheme = ""
        userinfo: str | None = None
        host = ""
        port: int | None = None
        path = ""

        if "://" in url:
            scheme, rest = url.split("://", 1)
            if not scheme:
                raise ValueError("Missing scheme")
            if "/" in rest:
                authority, path = rest.split("/", 1)
                path = "/" + path
            else:
                authority = rest
                path = ""

            if "@" in authority:
                userinfo, authority = authority.rsplit("@", 1)
                if not userinfo:
                    userinfo = None

            if authority.startswith("["):
                if "]" not in authority:
                    raise ValueError("Unclosed IPv6 bracket")
                host, _, maybe_port = authority[1:].partition("]")
                if maybe_port:
                    if not maybe_port.startswith(":"):
                        raise ValueError("Invalid IPv6 syntax")
                    port_str = maybe_port[1:]
                    port = URLParser._parse_port(port_str) if port_str else None
                else:
                    port = _default_ports.get(scheme)
            else:
                if ":" in authority:
                    host, port_str = authority.rsplit(":", 1)
                    port = URLParser._parse_port(port_str)
                else:
                    host = authority
                    port = _default_ports.get(scheme)

            if not host:
                raise ValueError("Missing host")

        else:
            path = url

        return ParsedURL(
            scheme=scheme,
            userinfo=userinfo,
            host=host,
            port=port,
            path=path or "/",
            query=query,
            fragment=fragment,
        )

    @staticmethod
    def _parse_port(port_str: str) -> int:
        try:
            port = int(port_str)
        except ValueError:
            raise ValueError(f"Invalid port: {port_str}")
        if not (0 < port <= 65535):
            raise ValueError(f"Port out of range: {port}")
        return port
