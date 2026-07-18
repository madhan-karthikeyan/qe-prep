import os
import tempfile
import unittest

from database_utilities.implementation.crud import CRUD
from database_utilities.implementation.database import Database


class TestDatabaseIntegration(unittest.TestCase):
    def setUp(self) -> None:
        self.tmp = tempfile.NamedTemporaryFile(suffix=".db", delete=False)
        self.tmp.close()
        self.db_path = self.tmp.name

    def tearDown(self) -> None:
        os.unlink(self.db_path)

    def test_sqlite_persistence(self) -> None:
        db = Database(self.db_path)
        with db.connection():
            db.execute(
                "CREATE TABLE IF NOT EXISTS items ("
                "  id INTEGER PRIMARY KEY AUTOINCREMENT,"
                "  name TEXT,"
                "  price REAL"
                ")"
            )
            crud = CRUD(db)
            id1 = crud.create("items", {"name": "widget", "price": 9.99})
            id2 = crud.create("items", {"name": "gadget", "price": 24.99})
            self.assertGreater(id1, 0)
            self.assertGreater(id2, 0)

        db2 = Database(self.db_path)
        with db2.connection():
            crud2 = CRUD(db2)
            rows = crud2.read("items")
            self.assertEqual(len(rows), 2)
            names = {r["name"] for r in rows}
            self.assertIn("widget", names)
            self.assertIn("gadget", names)

    def test_persistence_with_update(self) -> None:
        db = Database(self.db_path)
        with db.connection():
            db.execute(
                "CREATE TABLE IF NOT EXISTS config ("
                "  key TEXT PRIMARY KEY,"
                "  value TEXT"
                ")"
            )
            crud = CRUD(db)
            crud.create("config", {"key": "theme", "value": "dark"})

        db2 = Database(self.db_path)
        with db2.connection():
            crud2 = CRUD(db2)
            crud2.update("config", {"key": "theme"}, {"value": "light"})

        db3 = Database(self.db_path)
        with db3.connection():
            crud3 = CRUD(db3)
            rows = crud3.read("config", {"key": "theme"})
            self.assertEqual(len(rows), 1)
            self.assertEqual(rows[0]["value"], "light")

    def test_concurrent_persistence(self) -> None:
        db = Database(self.db_path)
        with db.connection():
            db.execute(
                "CREATE TABLE IF NOT EXISTS counters ("
                "  id INTEGER PRIMARY KEY,"
                "  count INTEGER"
                ")"
            )
            crud = CRUD(db)
            crud.create("counters", {"id": 1, "count": 0})
            crud.create("counters", {"id": 2, "count": 0})

        import threading
        _ = threading.Lock()

        def increment(cid: int) -> None:
            local_db = Database(self.db_path)
            with local_db.connection():
                local_crud = CRUD(local_db)
                for _ in range(100):
                    rows = local_crud.read("counters", {"id": cid})
                    current = rows[0]["count"]
                    local_crud.update(
                        "counters", {"id": cid}, {"count": current + 1}
                    )

        threads = [
            threading.Thread(target=increment, args=(1,)),
            threading.Thread(target=increment, args=(2,)),
        ]
        for t in threads:
            t.start()
        for t in threads:
            t.join()

        dbv = Database(self.db_path)
        with dbv.connection():
            crudv = CRUD(dbv)
            rows = crudv.read("counters", {"id": 1})
            self.assertEqual(rows[0]["count"], 100)


if __name__ == "__main__":
    unittest.main()
