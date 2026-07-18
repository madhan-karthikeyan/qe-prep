package database_utilities

import (
	"testing"
)

func TestCreateTable(t *testing.T) {
	db := Open()
	defer db.Close()
	err := db.Execute("CREATE TABLE users (id INT, name TEXT)")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	cols, err := db.TableColumns("users")
	if err != nil {
		t.Fatalf("TableColumns: %v", err)
	}
	if len(cols) != 2 || cols[0] != "id" || cols[1] != "name" {
		t.Errorf("unexpected columns: %v", cols)
	}
}

func TestInsertAndQuery(t *testing.T) {
	db := Open()
	defer db.Close()
	db.Execute("CREATE TABLE users (id INT, name TEXT)")

	err := db.Execute("INSERT INTO users VALUES (?, ?)", 1, "alice")
	if err != nil {
		t.Fatalf("Insert: %v", err)
	}
	err = db.Execute("INSERT INTO users VALUES (?, ?)", 2, "bob")
	if err != nil {
		t.Fatalf("Insert: %v", err)
	}

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		t.Fatalf("Query: %v", err)
	}

	var count int
	for rows.Next() {
		count++
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			t.Fatalf("Scan: %v", err)
		}
	}
	if count != 2 {
		t.Errorf("expected 2 rows, got %d", count)
	}
}

func TestSelectWithWhere(t *testing.T) {
	db := Open()
	defer db.Close()
	db.Execute("CREATE TABLE users (id INT, name TEXT)")
	db.Execute("INSERT INTO users VALUES (?, ?)", 1, "alice")
	db.Execute("INSERT INTO users VALUES (?, ?)", 2, "bob")

	rows, err := db.Query("SELECT * FROM users WHERE name = ?", "alice")
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if !rows.Next() {
		t.Fatal("expected a row")
	}
	var id, name string
	if err := rows.Scan(&id, &name); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if rows.Next() {
		t.Error("expected only one row")
	}
}

func TestUpdate(t *testing.T) {
	db := Open()
	defer db.Close()
	db.Execute("CREATE TABLE users (id INT, name TEXT)")
	db.Execute("INSERT INTO users VALUES (?, ?)", 1, "alice")
	db.Execute("INSERT INTO users VALUES (?, ?)", 2, "bob")

	err := db.Execute("UPDATE users SET name = ? WHERE id = ?", "charlie", 1)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	rows, _ := db.Query("SELECT * FROM users WHERE id = ?", 1)
	if !rows.Next() {
		t.Fatal("expected row")
	}
	var id, name string
	rows.Scan(&id, &name)
	if name != "charlie" {
		t.Errorf("name = %q, want 'charlie'", name)
	}
}

func TestDelete(t *testing.T) {
	db := Open()
	defer db.Close()
	db.Execute("CREATE TABLE users (id INT, name TEXT)")
	db.Execute("INSERT INTO users VALUES (?, ?)", 1, "alice")
	db.Execute("INSERT INTO users VALUES (?, ?)", 2, "bob")

	err := db.Execute("DELETE FROM users WHERE id = ?", 1)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	rows, _ := db.Query("SELECT * FROM users")
	var count int
	for rows.Next() {
		count++
	}
	if count != 1 {
		t.Errorf("expected 1 row, got %d", count)
	}
}

func TestDropTable(t *testing.T) {
	db := Open()
	defer db.Close()
	db.Execute("CREATE TABLE users (id INT)")
	err := db.Execute("DROP TABLE users")
	if err != nil {
		t.Fatalf("Drop: %v", err)
	}
	_, err = db.TableColumns("users")
	if err == nil {
		t.Error("expected error for dropped table")
	}
}

func TestCRUDCreate(t *testing.T) {
	db := Open()
	defer db.Close()
	db.Execute("CREATE TABLE items (id INT, name TEXT, price FLOAT)")

	err := Create(db, "items", map[string]interface{}{
		"id":    1,
		"name":  "widget",
		"price": 9.99,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	results, err := Read(db, "items", nil)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestCRUDReadWithConditions(t *testing.T) {
	db := Open()
	defer db.Close()
	db.Execute("CREATE TABLE items (id INT, name TEXT)")
	Create(db, "items", map[string]interface{}{"id": 1, "name": "a"})
	Create(db, "items", map[string]interface{}{"id": 2, "name": "b"})

	results, err := Read(db, "items", map[string]interface{}{"id": 1})
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestCRUDDelete(t *testing.T) {
	db := Open()
	defer db.Close()
	db.Execute("CREATE TABLE items (id INT, name TEXT)")
	Create(db, "items", map[string]interface{}{"id": 1, "name": "a"})
	Create(db, "items", map[string]interface{}{"id": 2, "name": "b"})

	_, err := Delete(db, "items", map[string]interface{}{"id": 1})
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	results, _ := Read(db, "items", nil)
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestTransactionCommit(t *testing.T) {
	db := Open()
	defer db.Close()
	db.Execute("CREATE TABLE items (id INT, name TEXT)")

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Begin: %v", err)
	}
	tx.Execute("INSERT INTO items VALUES (?, ?)", 1, "tx-item")
	if err := tx.Commit(); err != nil {
		t.Fatalf("Commit: %v", err)
	}

	rows, _ := db.Query("SELECT * FROM items")
	if !rows.Next() {
		t.Error("expected to see committed row")
	}
}

func TestTransactionRollback(t *testing.T) {
	db := Open()
	defer db.Close()
	db.Execute("CREATE TABLE items (id INT, name TEXT)")

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Begin: %v", err)
	}
	tx.Execute("INSERT INTO items VALUES (?, ?)", 1, "tx-item")
	if err := tx.Rollback(); err != nil {
		t.Fatalf("Rollback: %v", err)
	}

	rows, _ := db.Query("SELECT * FROM items")
	if rows.Next() {
		t.Error("expected no rows after rollback")
	}
}

func TestTransactionIsolation(t *testing.T) {
	db := Open()
	defer db.Close()
	db.Execute("CREATE TABLE items (id INT, name TEXT)")
	db.Execute("INSERT INTO items VALUES (?, ?)", 1, "original")

	tx, _ := db.Begin()
	tx.Execute("INSERT INTO items VALUES (?, ?)", 2, "in-tx")

	rows, _ := db.Query("SELECT * FROM items")
	var count int
	for rows.Next() {
		count++
	}
	if count != 1 {
		t.Errorf("expected 1 row outside tx, got %d", count)
	}

	tx.Rollback()
}

func TestDuplicateTableError(t *testing.T) {
	db := Open()
	defer db.Close()
	db.Execute("CREATE TABLE t (col INT)")
	err := db.Execute("CREATE TABLE t (col INT)")
	if err == nil {
		t.Error("expected error for duplicate table")
	}
}

func TestQueryEmptyTable(t *testing.T) {
	db := Open()
	defer db.Close()
	db.Execute("CREATE TABLE t (col INT)")
	rows, _ := db.Query("SELECT * FROM t")
	if rows.Next() {
		t.Error("expected no rows")
	}
}
