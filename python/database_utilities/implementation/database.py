import sqlite3
from collections.abc import Generator
from contextlib import contextmanager
from typing import Any


class Database:
    def __init__(self, db_path: str = ":memory:") -> None:
        self._path = db_path
        self._conn: sqlite3.Connection | None = None
        self._in_memory_default = db_path == ":memory:"

    def connect(self) -> None:
        if self._conn is not None:
            return
        self._conn = sqlite3.connect(self._path)
        self._conn.row_factory = sqlite3.Row
        self._conn.execute("PRAGMA journal_mode=WAL")

    def close(self) -> None:
        if self._conn is not None:
            self._conn.close()
            self._conn = None

    def execute(self, query: str, params: tuple[Any, ...] = ()) -> sqlite3.Cursor:
        self._ensure_connected()
        assert self._conn is not None
        return self._conn.execute(query, params)

    def executemany(self, query: str, params: list[tuple[Any, ...]]) -> sqlite3.Cursor:
        self._ensure_connected()
        assert self._conn is not None
        return self._conn.executemany(query, params)

    def fetchone(self, query: str, params: tuple[Any, ...] = ()) -> sqlite3.Row | None:
        self._ensure_connected()
        assert self._conn is not None
        cursor = self._conn.execute(query, params)
        row = cursor.fetchone()
        return row if isinstance(row, sqlite3.Row) else None

    def fetchall(self, query: str, params: tuple[Any, ...] = ()) -> list[sqlite3.Row]:
        self._ensure_connected()
        assert self._conn is not None
        cursor = self._conn.execute(query, params)
        return cursor.fetchall()

    def begin(self) -> None:
        self._ensure_connected()
        assert self._conn is not None
        self._conn.execute("BEGIN")

    def commit(self) -> None:
        self._ensure_connected()
        assert self._conn is not None
        self._conn.commit()

    def rollback(self) -> None:
        self._ensure_connected()
        assert self._conn is not None
        self._conn.rollback()

    @contextmanager
    def connection(self) -> Generator["Database", None, None]:
        self.connect()
        try:
            yield self
        finally:
            self.close()

    @contextmanager
    def transaction(self) -> Generator["Database", None, None]:
        self._ensure_connected()
        assert self._conn is not None
        self.begin()
        try:
            yield self
            self.commit()
        except Exception:
            self.rollback()
            raise

    def table_exists(self, table_name: str) -> bool:
        result = self.fetchone(
            "SELECT name FROM sqlite_master WHERE type='table' AND name=?",
            (table_name,),
        )
        return result is not None

    def _ensure_connected(self) -> None:
        if self._conn is None:
            self.connect()

    @property
    def is_connected(self) -> bool:
        return self._conn is not None
