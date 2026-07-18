from __future__ import annotations


class _TrieNode:
    __slots__ = ("children", "is_end")

    def __init__(self) -> None:
        self.children: dict[str, _TrieNode] = {}
        self.is_end: bool = False


class Trie:
    def __init__(self, case_insensitive: bool = False) -> None:
        self._root = _TrieNode()
        self._case_insensitive = case_insensitive
        self._size = 0

    def _normalize(self, word: str) -> str:
        return word.lower() if self._case_insensitive else word

    def insert(self, word: str) -> None:
        word = self._normalize(word)
        node = self._root
        for ch in word:
            if ch not in node.children:
                node.children[ch] = _TrieNode()
            node = node.children[ch]
        if not node.is_end:
            node.is_end = True
            self._size += 1

    def search(self, word: str) -> bool:
        word = self._normalize(word)
        node = self._root
        for ch in word:
            if ch not in node.children:
                return False
            node = node.children[ch]
        return node.is_end

    def starts_with(self, prefix: str) -> bool:
        prefix = self._normalize(prefix)
        node = self._root
        for ch in prefix:
            if ch not in node.children:
                return False
            node = node.children[ch]
        return True

    def count_prefix(self, prefix: str) -> int:
        prefix = self._normalize(prefix)
        node = self._root
        for ch in prefix:
            if ch not in node.children:
                return 0
            node = node.children[ch]
        return self._count_words(node)

    def _count_words(self, node: _TrieNode) -> int:
        count = 1 if node.is_end else 0
        for child in node.children.values():
            count += self._count_words(child)
        return count

    def delete(self, word: str) -> bool:
        word = self._normalize(word)
        return self._delete(self._root, word, 0)

    def _delete(self, node: _TrieNode, word: str, depth: int) -> bool:
        if depth == len(word):
            if not node.is_end:
                return False
            node.is_end = False
            self._size -= 1
            return len(node.children) == 0
        ch = word[depth]
        if ch not in node.children:
            return False
        should_delete_child = self._delete(node.children[ch], word, depth + 1)
        if should_delete_child:
            del node.children[ch]
            return len(node.children) == 0 and not node.is_end
        return False

    @property
    def size(self) -> int:
        return self._size
