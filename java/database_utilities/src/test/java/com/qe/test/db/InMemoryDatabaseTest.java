package com.qe.test.db;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

@DisplayName("InMemoryDatabase")
class InMemoryDatabaseTest {

    private InMemoryDatabase db;

    @BeforeEach
    void setUp() {
        db = new InMemoryDatabase();
        db.createTable("users", List.of(
                new InMemoryDatabase.Column("id", Integer.class),
                new InMemoryDatabase.Column("name", String.class)
        ));
    }

    @Test
    @DisplayName("creates table and inserts rows")
    void createAndInsert() {
        db.insert("users", Map.of("id", 1, "name", "Alice"));
        db.insert("users", Map.of("id", 2, "name", "Bob"));

        var results = db.executeQuery("users");
        assertEquals(2, results.size());
    }

    @Test
    @DisplayName("selects rows with predicate")
    void selectWithPredicate() {
        db.insert("users", Map.of("id", 1, "name", "Alice"));
        db.insert("users", Map.of("id", 2, "name", "Bob"));

        var results = db.executeQuery("users", row -> "Alice".equals(row.get("name")));
        assertEquals(1, results.size());
        assertEquals(1, results.getFirst().get("id"));
    }

    @Test
    @DisplayName("updates matching rows")
    void updateRows() {
        db.insert("users", Map.of("id", 1, "name", "Alice"));
        db.insert("users", Map.of("id", 2, "name", "Bob"));

        int updated = db.update("users", Map.of("name", "Updated"), row -> true);
        assertEquals(2, updated);
        var results = db.executeQuery("users");
        assertEquals("Updated", results.getFirst().get("name"));
    }

    @Test
    @DisplayName("deletes matching rows")
    void deleteRows() {
        db.insert("users", Map.of("id", 1, "name", "Alice"));
        db.insert("users", Map.of("id", 2, "name", "Bob"));

        int deleted = db.delete("users", row -> (int) row.get("id") == 1);
        assertEquals(1, deleted);
        assertEquals(1, db.executeQuery("users").size());
    }

    @Test
    @DisplayName("throws when creating duplicate table")
    void duplicateTable() {
        assertThrows(IllegalArgumentException.class,
                () -> db.createTable("users", List.of(new InMemoryDatabase.Column("x", String.class))));
    }

    @Test
    @DisplayName("throws when dropping non-existent table")
    void dropNonExistentTable() {
        assertThrows(IllegalArgumentException.class, () -> db.dropTable("nonexistent"));
    }

    @Test
    @DisplayName("throws when querying non-existent table")
    void queryNonExistentTable() {
        assertThrows(IllegalArgumentException.class, () -> db.executeQuery("nonexistent"));
    }

    @Test
    @DisplayName("rejects null column list")
    void rejectsNullColumns() {
        assertThrows(NullPointerException.class, () -> db.createTable("t", null));
    }

    @Test
    @DisplayName("rejects empty column list")
    void rejectsEmptyColumns() {
        assertThrows(IllegalArgumentException.class, () -> db.createTable("t", List.of()));
    }

    @Test
    @DisplayName("supports begin/commit transaction")
    void transactionCommit() {
        db.beginTransaction();
        db.insert("users", Map.of("id", 1, "name", "Alice"));
        db.commit();
        assertEquals(1, db.executeQuery("users").size());
    }

    @Test
    @DisplayName("supports begin/rollback transaction")
    void transactionRollback() {
        db.beginTransaction();
        db.insert("users", Map.of("id", 1, "name", "Alice"));
        db.rollback();
        assertEquals(0, db.executeQuery("users").size());
    }

    @Test
    @DisplayName("throws when no active transaction on commit")
    void commitWithoutTransaction() {
        assertThrows(IllegalStateException.class, () -> db.commit());
    }

    @Test
    @DisplayName("throws when no active transaction on rollback")
    void rollbackWithoutTransaction() {
        assertThrows(IllegalStateException.class, () -> db.rollback());
    }

    @Test
    @DisplayName("tracks transaction state")
    void tracksTransactionState() {
        assertFalse(db.isInTransaction());
        db.beginTransaction();
        assertTrue(db.isInTransaction());
        db.commit();
        assertFalse(db.isInTransaction());
    }
}
