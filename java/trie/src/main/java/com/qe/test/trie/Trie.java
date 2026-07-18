package com.qe.test.trie;

import java.util.HashMap;
import java.util.Map;

public class Trie {
    private final Node root;
    private final boolean caseInsensitive;

    public Trie() {
        this(false);
    }

    public Trie(boolean caseInsensitive) {
        this.root = new Node();
        this.caseInsensitive = caseInsensitive;
    }

    public void insert(String word) {
        if (word == null) throw new IllegalArgumentException("word must not be null");
        var node = root;
        var s = normalize(word);
        for (int i = 0; i < s.length(); i++) {
            char c = s.charAt(i);
            node = node.children.computeIfAbsent(c, k -> new Node());
        }
        node.isEnd = true;
    }

    public boolean search(String word) {
        if (word == null) return false;
        var node = traverse(normalize(word));
        return node != null && node.isEnd;
    }

    public boolean startsWith(String prefix) {
        if (prefix == null) return false;
        return traverse(normalize(prefix)) != null;
    }

    public boolean delete(String word) {
        if (word == null) throw new IllegalArgumentException("word must not be null");
        var found = new boolean[]{false};
        delete(root, normalize(word), 0, found);
        return found[0];
    }

    public int countPrefix(String prefix) {
        if (prefix == null) return 0;
        var node = traverse(normalize(prefix));
        if (node == null) return 0;
        return countWords(node);
    }

    private void delete(Node node, String word, int depth, boolean[] found) {
        if (depth == word.length()) {
            if (node.isEnd) {
                found[0] = true;
                node.isEnd = false;
            }
            return;
        }
        char c = word.charAt(depth);
        var child = node.children.get(c);
        if (child == null) return;
        delete(child, word, depth + 1, found);
        if (found[0] && child.children.isEmpty() && !child.isEnd) {
            node.children.remove(c);
        }
    }

    private Node traverse(String s) {
        var node = root;
        for (int i = 0; i < s.length(); i++) {
            node = node.children.get(s.charAt(i));
            if (node == null) return null;
        }
        return node;
    }

    private int countWords(Node node) {
        int count = node.isEnd ? 1 : 0;
        for (var child : node.children.values()) {
            count += countWords(child);
        }
        return count;
    }

    private String normalize(String s) {
        return caseInsensitive ? s.toLowerCase() : s;
    }

    private static class Node {
        final Map<Character, Node> children = new HashMap<>();
        boolean isEnd;
    }
}
