package com.qe.test.trie;

import static org.junit.jupiter.api.Assertions.*;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

@DisplayName("Trie")
class TrieTest {

    @Test
    @DisplayName("insert and search basic words")
    void insertAndSearch() {
        var trie = new Trie();
        trie.insert("hello");
        trie.insert("world");
        assertTrue(trie.search("hello"));
        assertTrue(trie.search("world"));
        assertFalse(trie.search("hell"));
        assertFalse(trie.search("worlds"));
    }

    @Test
    @DisplayName("startsWith returns true for prefixes")
    void startsWith() {
        var trie = new Trie();
        trie.insert("apple");
        trie.insert("app");
        assertTrue(trie.startsWith("ap"));
        assertTrue(trie.startsWith("app"));
        assertTrue(trie.startsWith("apple"));
        assertFalse(trie.startsWith("applz"));
    }

    @Test
    @DisplayName("delete removes word but keeps prefixes")
    void delete() {
        var trie = new Trie();
        trie.insert("apple");
        trie.insert("app");
        assertTrue(trie.delete("apple"));
        assertFalse(trie.search("apple"));
        assertTrue(trie.search("app"));
        assertFalse(trie.delete("apple"));
    }

    @Test
    @DisplayName("countPrefix counts words with given prefix")
    void countPrefix() {
        var trie = new Trie();
        trie.insert("dog");
        trie.insert("door");
        trie.insert("dorm");
        trie.insert("cat");
        assertEquals(3, trie.countPrefix("do"));
        assertEquals(1, trie.countPrefix("doo"));
        assertEquals(1, trie.countPrefix("cat"));
        assertEquals(0, trie.countPrefix("z"));
    }

    @Test
    @DisplayName("case insensitive mode")
    void caseInsensitive() {
        var trie = new Trie(true);
        trie.insert("Hello");
        assertTrue(trie.search("hello"));
        assertTrue(trie.search("HELLO"));
        assertTrue(trie.search("Hello"));
    }

    @Test
    @DisplayName("empty string is inserted and found")
    void emptyString() {
        var trie = new Trie();
        assertFalse(trie.search(""));
        trie.insert("");
        assertTrue(trie.search(""));
        assertTrue(trie.startsWith(""));
    }

    @Test
    @DisplayName("null arguments throw")
    void nullArguments() {
        var trie = new Trie();
        assertThrows(IllegalArgumentException.class, () -> trie.insert(null));
        assertThrows(IllegalArgumentException.class, () -> trie.delete(null));
    }

    @Test
    @DisplayName("non-existent delete returns false")
    void deleteNonExistent() {
        var trie = new Trie();
        trie.insert("hello");
        assertFalse(trie.delete("world"));
        assertFalse(trie.delete("hell"));
    }

    @Test
    @DisplayName("search for null returns false")
    void searchNull() {
        var trie = new Trie();
        assertFalse(trie.search(null));
        assertFalse(trie.startsWith(null));
    }
}
