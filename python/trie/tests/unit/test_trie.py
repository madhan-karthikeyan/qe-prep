from __future__ import annotations

from trie.implementation.trie import Trie


class TestTrie:
    def test_insert_and_search(self) -> None:
        t = Trie()
        t.insert("hello")
        assert t.search("hello")
        assert not t.search("world")

    def test_starts_with(self) -> None:
        t = Trie()
        t.insert("hello")
        assert t.starts_with("hel")
        assert not t.starts_with("wor")

    def test_delete_word(self) -> None:
        t = Trie()
        t.insert("hello")
        assert t.delete("hello")
        assert not t.search("hello")

    def test_delete_non_existent(self) -> None:
        t = Trie()
        t.insert("hello")
        assert not t.delete("world")

    def test_overlapping_prefixes(self) -> None:
        t = Trie()
        t.insert("app")
        t.insert("apple")
        assert t.search("app")
        assert t.search("apple")
        t.delete("app")
        assert not t.search("app")
        assert t.search("apple")

    def test_empty_string(self) -> None:
        t = Trie()
        t.insert("")
        assert t.search("")
        assert t.starts_with("")

    def test_case_sensitive_default(self) -> None:
        t = Trie()
        t.insert("Hello")
        assert t.search("Hello")
        assert not t.search("hello")

    def test_case_insensitive(self) -> None:
        t = Trie(case_insensitive=True)
        t.insert("Hello")
        assert t.search("hello")
        assert t.search("HELLO")
        assert t.search("Hello")

    def test_count_prefix(self) -> None:
        t = Trie()
        words = ["app", "apple", "application", "appetite", "banana"]
        for w in words:
            t.insert(w)
        assert t.count_prefix("app") == 4
        assert t.count_prefix("ban") == 1
        assert t.count_prefix("xyz") == 0

    def test_size(self) -> None:
        t = Trie()
        assert t.size == 0
        t.insert("a")
        assert t.size == 1
        t.insert("a")
        assert t.size == 1
        t.insert("b")
        assert t.size == 2
        t.delete("a")
        assert t.size == 1
