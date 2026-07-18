package com.qe.test.db;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

@DisplayName("CrudOperations")
class CrudOperationsTest {

    private InMemoryDatabase db;
    private CrudOperations crud;

    @BeforeEach
    void setUp() {
        db = new InMemoryDatabase();
        crud = new CrudOperations(db);
        crud.createTable("items", List.of(
                new InMemoryDatabase.Column("id", Integer.class),
                new InMemoryDatabase.Column("name", String.class),
                new InMemoryDatabase.Column("value", Double.class)
        ));
    }

    @Test
    @DisplayName("inserts and selects all")
    void insertAndSelectAll() {
        crud.insert("items", Map.of("id", 1, "name", "item1", "value", 10.0));
        crud.insert("items", Map.of("id", 2, "name", "item2", "value", 20.0));

        var results = crud.select("items");
        assertEquals(2, results.size());
    }

    @Test
    @DisplayName("selects by id")
    void selectById() {
        crud.insert("items", Map.of("id", 1, "name", "item1", "value", 10.0));
        crud.insert("items", Map.of("id", 2, "name", "item2", "value", 20.0));

        var result = crud.selectById("items", "id", 1);
        assertTrue(result.isPresent());
        assertEquals("item1", result.get().get("name"));
    }

    @Test
    @DisplayName("selects with predicate")
    void selectWhere() {
        crud.insert("items", Map.of("id", 1, "name", "item1", "value", 10.0));
        crud.insert("items", Map.of("id", 2, "name", "item2", "value", 20.0));

        var results = crud.selectWhere("items", row -> (double) row.get("value") > 15.0);
        assertEquals(1, results.size());
        assertEquals(2, results.getFirst().get("id"));
    }

    @Test
    @DisplayName("updates by id")
    void updateById() {
        crud.insert("items", Map.of("id", 1, "name", "item1", "value", 10.0));

        int count = crud.updateById("items", "id", 1, Map.of("name", "updated"));
        assertEquals(1, count);

        var result = crud.selectById("items", "id", 1);
        assertEquals("updated", result.get().get("name"));
    }

    @Test
    @DisplayName("deletes by id")
    void deleteById() {
        crud.insert("items", Map.of("id", 1, "name", "item1", "value", 10.0));
        crud.insert("items", Map.of("id", 2, "name", "item2", "value", 20.0));

        int count = crud.deleteById("items", "id", 1);
        assertEquals(1, count);
        assertEquals(1, crud.select("items").size());
    }

    @Test
    @DisplayName("rejects null database")
    void rejectsNullDatabase() {
        assertThrows(NullPointerException.class, () -> new CrudOperations(null));
    }

    @Test
    @DisplayName("supports transactions via crud")
    void transactionSupport() {
        crud.beginTransaction();
        crud.insert("items", Map.of("id", 1, "name", "tx-item", "value", 1.0));
        crud.rollback();
        assertEquals(0, crud.select("items").size());
    }
}
