import sqlite3
import unittest

from database_utilities.implementation.crud import CRUD
from database_utilities.implementation.database import Database


class TestDatabase(unittest.TestCase):
    def setUp(self) -> None:
        self.db = Database()
        self.db.connect()
        self.db.execute(
            "CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT, value INTEGER)"
        )

    def tearDown(self) -> None:
        self.db.close()

    def test_execute_and_fetchall(self) -> None:
        self.db.execute("INSERT INTO test (name, value) VALUES (?, ?)", ("a", 1))
        rows = self.db.fetchall("SELECT * FROM test")
        self.assertEqual(len(rows), 1)
        self.assertEqual(rows[0]["name"], "a")

    def test_fetchone(self) -> None:
        self.db.execute("INSERT INTO test (name, value) VALUES (?, ?)", ("a", 1))
        row = self.db.fetchone("SELECT * FROM test WHERE name=?", ("a",))
        assert row is not None
        self.assertEqual(row["value"], 1)

    def test_fetchone_no_result(self) -> None:
        row = self.db.fetchone("SELECT * FROM test WHERE name=?", ("nonexistent",))
        self.assertIsNone(row)

    def test_executemany(self) -> None:
        data = [("a", 1), ("b", 2), ("c", 3)]
        self.db.executemany(
            "INSERT INTO test (name, value) VALUES (?, ?)",
            [(n, v) for n, v in data],
        )
        rows = self.db.fetchall("SELECT * FROM test ORDER BY id")
        self.assertEqual(len(rows), 3)

    def test_table_exists(self) -> None:
        self.assertTrue(self.db.table_exists("test"))
        self.assertFalse(self.db.table_exists("nonexistent"))

    def test_connect_twice(self) -> None:
        self.db.connect()
        self.assertTrue(self.db.is_connected)

    def test_close(self) -> None:
        self.db.close()
        self.assertFalse(self.db.is_connected)


class TestCRUD(unittest.TestCase):
    def setUp(self) -> None:
        self.db = Database()
        self.db.connect()
        self.db.execute(
            "CREATE TABLE users ("
            "  id INTEGER PRIMARY KEY AUTOINCREMENT,"
            "  name TEXT NOT NULL,"
            "  email TEXT UNIQUE"
            ")"
        )
        self.crud = CRUD(self.db)

    def tearDown(self) -> None:
        self.db.close()

    def test_create(self) -> None:
        user_id = self.crud.create("users", {"name": "Alice", "email": "alice@x.com"})
        self.assertIsInstance(user_id, int)
        self.assertGreater(user_id, 0)

    def test_read_all(self) -> None:
        self.crud.create("users", {"name": "A", "email": "a@x.com"})
        self.crud.create("users", {"name": "B", "email": "b@x.com"})
        rows = self.crud.read("users")
        self.assertEqual(len(rows), 2)

    def test_read_with_conditions(self) -> None:
        self.crud.create("users", {"name": "Alice", "email": "alice@x.com"})
        rows = self.crud.read("users", {"name": "Alice"})
        self.assertEqual(len(rows), 1)
        self.assertEqual(rows[0]["email"], "alice@x.com")

    def test_update(self) -> None:
        uid = self.crud.create("users", {"name": "Alice", "email": "alice@x.com"})
        count = self.crud.update(
            "users", {"id": uid}, {"name": "Alicia"}
        )
        self.assertEqual(count, 1)
        rows = self.crud.read("users", {"id": uid})
        self.assertEqual(rows[0]["name"], "Alicia")

    def test_delete(self) -> None:
        uid = self.crud.create("users", {"name": "Alice", "email": "alice@x.com"})
        count = self.crud.delete("users", {"id": uid})
        self.assertEqual(count, 1)
        rows = self.crud.read("users", {"id": uid})
        self.assertEqual(len(rows), 0)

    def test_missing_table(self) -> None:
        with self.assertRaises(sqlite3.OperationalError):
            self.crud.read("nonexistent")

    def test_transaction_commit(self) -> None:
        with self.db.transaction():
            self.crud.create("users", {"name": "T1", "email": "t1@x.com"})
        rows = self.crud.read("users")
        self.assertEqual(len(rows), 1)

    def test_transaction_rollback(self) -> None:
        try:
            with self.db.transaction():
                self.crud.create("users", {"name": "T2", "email": "t2@x.com"}, commit=False)
                raise RuntimeError("rollback")
        except RuntimeError:
            pass
        rows = self.crud.read("users")
        self.assertEqual(len(rows), 0)

    def test_connection_context_manager(self) -> None:
        db = Database()
        with db.connection():
            db.execute("CREATE TABLE t (x INTEGER)")
            db.execute("INSERT INTO t VALUES (42)")
            row = db.fetchone("SELECT x FROM t")
            assert row is not None
            self.assertEqual(row["x"], 42)
        self.assertFalse(db.is_connected)

    def test_begin_commit(self) -> None:
        self.db.begin()
        self.crud.create("users", {"name": "B", "email": "b@x.com"})
        self.db.commit()
        rows = self.crud.read("users")
        self.assertEqual(len(rows), 1)

    def test_begin_rollback(self) -> None:
        self.db.begin()
        self.crud.create("users", {"name": "C", "email": "c@x.com"}, commit=False)
        self.db.rollback()
        rows = self.crud.read("users")
        self.assertEqual(len(rows), 0)


if __name__ == "__main__":
    unittest.main()
