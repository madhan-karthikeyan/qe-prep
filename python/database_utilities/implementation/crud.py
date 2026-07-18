import sqlite3
from typing import Any

from database_utilities.implementation.database import Database


class CRUD:
    def __init__(self, db: Database) -> None:
        self._db = db

    def create(self, table: str, data: dict[str, Any], commit: bool = True) -> int:
        columns = ", ".join(data.keys())
        placeholders = ", ".join("?" for _ in data)
        query = f"INSERT INTO {table} ({columns}) VALUES ({placeholders})"
        cursor = self._db.execute(query, tuple(data.values()))
        if commit:
            self._db.commit()
        assert cursor.lastrowid is not None
        return cursor.lastrowid

    def read(
        self,
        table: str,
        conditions: dict[str, Any] | None = None,
    ) -> list[sqlite3.Row]:
        query = f"SELECT * FROM {table}"
        params: tuple[Any, ...] = ()
        if conditions:
            where_clause = " AND ".join(f"{k} = ?" for k in conditions)
            query = f"{query} WHERE {where_clause}"
            params = tuple(conditions.values())
        return self._db.fetchall(query, params)

    def update(
        self,
        table: str,
        conditions: dict[str, Any],
        data: dict[str, Any],
        commit: bool = True,
    ) -> int:
        set_clause = ", ".join(f"{k} = ?" for k in data)
        where_clause = " AND ".join(f"{k} = ?" for k in conditions)
        query = f"UPDATE {table} SET {set_clause} WHERE {where_clause}"
        params = tuple(data.values()) + tuple(conditions.values())
        cursor = self._db.execute(query, params)
        if commit:
            self._db.commit()
        return cursor.rowcount

    def delete(
        self,
        table: str,
        conditions: dict[str, Any],
        commit: bool = True,
    ) -> int:
        where_clause = " AND ".join(f"{k} = ?" for k in conditions)
        query = f"DELETE FROM {table} WHERE {where_clause}"
        cursor = self._db.execute(query, tuple(conditions.values()))
        if commit:
            self._db.commit()
        return cursor.rowcount
