package com.qe.test.db;

import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.Objects;
import java.util.function.Predicate;

public class CrudOperations {
    private final InMemoryDatabase database;

    public CrudOperations(InMemoryDatabase database) {
        this.database = Objects.requireNonNull(database, "database must not be null");
    }

    public void createTable(String tableName, List<InMemoryDatabase.Column> columns) {
        database.createTable(tableName, columns);
    }

    public void insert(String tableName, Map<String, Object> row) {
        database.insert(tableName, row);
    }

    public List<Map<String, Object>> select(String tableName) {
        return database.executeQuery(tableName);
    }

    public List<Map<String, Object>> selectWhere(String tableName, Predicate<Map<String, Object>> predicate) {
        return database.executeQuery(tableName, predicate);
    }

    public Optional<Map<String, Object>> selectById(String tableName, String idColumn, Object idValue) {
        List<Map<String, Object>> results = database.executeQuery(
                tableName, row -> Objects.equals(row.get(idColumn), idValue));
        return results.isEmpty() ? Optional.empty() : Optional.of(results.getFirst());
    }

    public int update(String tableName, Map<String, Object> updates, Predicate<Map<String, Object>> predicate) {
        return database.update(tableName, updates, predicate);
    }

    public int updateById(String tableName, String idColumn, Object idValue, Map<String, Object> updates) {
        return database.update(tableName, updates, row -> Objects.equals(row.get(idColumn), idValue));
    }

    public int delete(String tableName, Predicate<Map<String, Object>> predicate) {
        return database.delete(tableName, predicate);
    }

    public int deleteById(String tableName, String idColumn, Object idValue) {
        return database.delete(tableName, row -> Objects.equals(row.get(idColumn), idValue));
    }

    public void beginTransaction() {
        database.beginTransaction();
    }

    public void commit() {
        database.commit();
    }

    public void rollback() {
        database.rollback();
    }

    public record Entity(String id, Map<String, Object> data) {
        public Entity {
            Objects.requireNonNull(id, "id must not be null");
            Objects.requireNonNull(data, "data must not be null");
        }
    }
}
