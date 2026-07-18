package com.qe.test.db;

import java.util.ArrayList;
import java.util.Collection;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.locks.ReentrantReadWriteLock;
import java.util.function.Predicate;
import java.util.stream.Collectors;

public class InMemoryDatabase implements AutoCloseable {
    private final Map<String, Table> tables = new ConcurrentHashMap<>();
    private final ThreadLocal<Transaction> currentTransaction = new ThreadLocal<>();

    public void createTable(String name, List<Column> columns) {
        Objects.requireNonNull(name, "name must not be null");
        Objects.requireNonNull(columns, "columns must not be null");
        if (columns.isEmpty()) {
            throw new IllegalArgumentException("Table must have at least one column");
        }
        if (tables.putIfAbsent(name, new Table(name, columns)) != null) {
            throw new IllegalArgumentException("Table already exists: " + name);
        }
    }

    public void dropTable(String name) {
        Objects.requireNonNull(name, "name must not be null");
        if (tables.remove(name) == null) {
            throw new IllegalArgumentException("Table not found: " + name);
        }
    }

    public List<Map<String, Object>> executeQuery(String tableName, Predicate<Map<String, Object>> predicate) {
        Objects.requireNonNull(tableName, "tableName must not be null");
        Table table = tables.get(tableName);
        if (table == null) {
            throw new IllegalArgumentException("Table not found: " + tableName);
        }
        table.readLock.lock();
        try {
            return table.rows.stream()
                    .filter(predicate)
                    .map(LinkedHashMap::new)
                    .collect(Collectors.toList());
        } finally {
            table.readLock.unlock();
        }
    }

    public List<Map<String, Object>> executeQuery(String tableName) {
        return executeQuery(tableName, row -> true);
    }

    public void executeUpdate(String tableName, Operation operation) {
        Objects.requireNonNull(tableName, "tableName must not be null");
        Objects.requireNonNull(operation, "operation must not be null");
        Table table = tables.get(tableName);
        if (table == null) {
            throw new IllegalArgumentException("Table not found: " + tableName);
        }

        table.writeLock.lock();
        try {
            Transaction tx = currentTransaction.get();
            if (tx != null) {
                tx.saveSnapshot(table);
            }
            operation.execute(table);
        } finally {
            table.writeLock.unlock();
        }
    }

    public void insert(String tableName, Map<String, Object> row) {
        executeUpdate(tableName, new InsertOperation(row));
    }

    public int update(String tableName, Map<String, Object> updates, Predicate<Map<String, Object>> predicate) {
        final int[] count = {0};
        executeUpdate(tableName, table -> {
            for (int i = 0; i < table.rows.size(); i++) {
                Map<String, Object> row = table.rows.get(i);
                if (predicate.test(row)) {
                    row.putAll(updates);
                    count[0]++;
                }
            }
        });
        return count[0];
    }

    public int delete(String tableName, Predicate<Map<String, Object>> predicate) {
        final int[] count = {0};
        executeUpdate(tableName, table -> {
            var iter = table.rows.iterator();
            while (iter.hasNext()) {
                if (predicate.test(iter.next())) {
                    iter.remove();
                    count[0]++;
                }
            }
        });
        return count[0];
    }

    public void beginTransaction() {
        if (currentTransaction.get() != null) {
            throw new IllegalStateException("Transaction already in progress");
        }
        currentTransaction.set(new Transaction());
    }

    public void commit() {
        Transaction tx = currentTransaction.get();
        if (tx == null) {
            throw new IllegalStateException("No active transaction");
        }
        tx.clear();
        currentTransaction.remove();
    }

    public void rollback() {
        Transaction tx = currentTransaction.get();
        if (tx == null) {
            throw new IllegalStateException("No active transaction");
        }
        tx.undo(this);
        tx.clear();
        currentTransaction.remove();
    }

    public boolean isInTransaction() {
        return currentTransaction.get() != null;
    }

    public Collection<String> getTableNames() {
        return Collections.unmodifiableSet(tables.keySet());
    }

    @Override
    public void close() {
        tables.clear();
    }

    public record Column(String name, Class<?> type) {
        public Column {
            Objects.requireNonNull(name, "name must not be null");
            Objects.requireNonNull(type, "type must not be null");
        }
    }

    private static class Table {
        final String name;
        final List<Column> columns;
        final List<Map<String, Object>> rows = new ArrayList<>();
        final ReentrantReadWriteLock.ReadLock readLock;
        final ReentrantReadWriteLock.WriteLock writeLock;
        private final ReentrantReadWriteLock rwLock = new ReentrantReadWriteLock();

        Table(String name, List<Column> columns) {
            this.name = name;
            this.columns = List.copyOf(columns);
            this.readLock = rwLock.readLock();
            this.writeLock = rwLock.writeLock();
        }
    }

    @FunctionalInterface
    public interface Operation {
        void execute(Table table);
    }

    private static class InsertOperation implements Operation {
        private final Map<String, Object> row;

        InsertOperation(Map<String, Object> row) {
            this.row = new LinkedHashMap<>(row);
        }

        @Override
        public void execute(Table table) {
            table.rows.add(new LinkedHashMap<>(row));
        }
    }

    private static class Transaction {
        private final List<Runnable> undoSnapshots = new ArrayList<>();

        void saveSnapshot(Table table) {
            var snapshot = table.rows.stream()
                    .map(LinkedHashMap::new)
                    .toList();
            undoSnapshots.add(() -> {
                table.rows.clear();
                table.rows.addAll(snapshot);
            });
        }

        void undo(InMemoryDatabase db) {
            for (int i = undoSnapshots.size() - 1; i >= 0; i--) {
                undoSnapshots.get(i).run();
            }
        }

        void clear() {
            undoSnapshots.clear();
        }
    }
}
