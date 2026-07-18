# Coding Round — QE Engineer Interview Guide

## Overview

The coding round evaluates your ability to write correct, efficient, and testable code. Unlike pure SWE interviews, QE coding rounds emphasize edge cases, error handling, testability, and understanding system behavior under failure.

## Top 25 Questions by Category

### String Manipulation (Difficulty: ★☆☆ Easy–★★☆ Medium)

1. **Reverse words in a string**
   - Approach: Split, reverse, join. Handle multiple spaces.
   - Test: Empty string, leading/trailing spaces, single word.

2. **Find first non-repeating character**
   - Approach: Two-pass hash map (counts). O(n) time, O(1) space (26 letters).
   - Test: All repeating, single char, empty string.

3. **Check if two strings are anagrams**
   - Approach: Character count array (256) or sort-and-compare.
   - Test: Unicode, case sensitivity, different lengths.

4. **String to integer (atoi)**
   - Approach: Handle whitespace, sign, overflow, invalid chars.
   - Test: `"   -42"`, `"4193 with words"`, `"   "`, `"2147483648"`.

5. **Longest substring without repeating characters**
   - Approach: Sliding window with hash set/map. O(n).
   - Test: All unique, all same, empty.

6. **Valid palindrome** (with non-alphanumeric chars)
   - Approach: Two-pointer with `isalnum` skip.
   - Test: `"A man, a plan, a canal: Panama"`, `"race a car"`.

7. **Implement strStr() / indexOf()**
   - Approach: Sliding window O(n*m). For interviews, KMP is overkill unless asked.
   - Test: Needle longer than haystack, empty needle.

### File Processing (Difficulty: ★★☆ Medium)

8. **Read a large file line by line**
   - Approach: `with open(...)` + iterator; never `readlines()` on huge files.
   - Test: File with 10M lines, encoding issues.

9. **Count word frequency in a file**
   - Approach: `defaultdict(int)`, normalize case, strip punctuation.
   - Test: Case sensitivity, punctuation, empty file.

10. **Merge multiple sorted log files**
    - Approach: Min-heap of (line, file_handle). O(N log k).
    - Test: Files of different lengths, empty files.

11. **Parse a structured log format (e.g., JSON lines)**
    - Approach: `json.loads()` per line, handle malformed lines gracefully.
    - Test: Corrupted lines, varying schemas, missing keys.

12. **Find top-K frequent items in a stream**
    - Approach: Counter + heap or quickselect. O(n log k).
    - Test: Stream with uniform distribution, small K.

### Data Structures (Difficulty: ★★☆ Medium–★★★ Hard)

13. **Valid parentheses** (stack)
    - Approach: Push opening brackets, pop on closing, check matching.
    - Test: `"([{}])"`, `"([)]"`, empty string.

14. **LRU Cache**
    - Approach: Doubly-linked list + hash map. O(1) get/put.
    - Test: Capacity=1, evict least recently used, concurrent access.

15. **Implement a Min Stack** (getMin in O(1))
    - Approach: Two stacks (values + current min) or tuple stack.
    - Test: Push after pop, duplicates.

16. **Binary search in rotated array**
    - Approach: Find pivot, then binary search on appropriate half.
    - Test: No rotation, single element, all duplicates.

17. **Design a thread-safe counter**
    - Approach: `threading.Lock` or `atomic` module.
    - Test: 1000 threads incrementing 10k times each.

### Concurrency Basics (Difficulty: ★★★ Hard)

18. **Producer-Consumer with bounded buffer**
    - Approach: `threading.Condition` or `queue.Queue`.
    - Test: Multiple producers/consumers, empty buffer, full buffer.

19. **Wait for all threads to complete** (join, barrier)
    - Approach: `Thread.join()`, `concurrent.futures`, or barrier.
    - Test: Thread raises exception — does main wait?

20. **Implement a rate limiter** (token bucket)
    - Approach: Timestamp-based refill with mutex.
    - Test: Burst behavior, concurrent requests, idle period.

21. **Race condition reproduction**
    - Approach: Write a test that reliably demonstrates the race (threading + sleep/barrier).
    - Test: The test fails without fix, passes with fix.

### Debugging & Testing (Difficulty: ★★☆ Medium)

22. **Fix a buggy implementation** (given code with off-by-one, null pointer, etc.)
    - Approach: Identify symptoms, trace inputs, fix root cause.
    - Test: Add tests that cover the edge case.

23. **Write a test for a given function**
    - Approach: Normal cases, edge cases, error cases, boundary values.
    - Test: Parametrize to cover equivalence partitions.

24. **Explain what this code does** (code review)
    - Approach: Read top-down, identify assumptions, potential bugs, performance issues.

25. **Design a test strategy for a function/API**
    - Approach: Input domains, invariants, property-based testing (Hypothesis/QuickCheck).

---

## How to Approach

| Step | Action | Time |
|------|--------|------|
| 1 | Clarify requirements: input size, constraints, edge cases | 2 min |
| 2 | Outline approach verbally before writing code | 3 min |
| 3 | Write solution with clean code and meaningful names | 15 min |
| 4 | Walk through example — trace with sample input | 3 min |
| 5 | Discuss test cases: normal, edge, error | 5 min |
| 6 | Analyze time/space complexity | 2 min |
| 7 | Discuss tradeoffs and alternatives | 5 min |

## Communication Strategy

- **Think aloud**: Say "I'm considering using a hash map here because we need O(1) lookups."
- **Ask clarifying questions**: "Should I handle integer overflow? What character encoding?"
- **Test your assumptions**: "I'm assuming input fits in memory — is that correct?"
- **Admit gaps**: "I don't remember the exact syntax for the threading module, but the logic is…"

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Jumping to code without clarifying | Ask about constraints first |
| Ignoring edge cases | Always test empty, single, duplicates |
| Writing O(n²) when O(n) is possible | Discuss tradeoffs before coding |
| Not testing the solution | Walk through the example manually |
| Staying silent | Narrate your thought process |
| Getting stuck on one approach | Say "This isn't working, let me try another approach" |

## Difficulty Levels

| Level | Description | Examples |
|-------|-------------|----------|
| ★☆☆ Easy | 1 data structure, straightforward logic | Reverse string, palindrome check |
| ★★☆ Medium | 2+ data structures or moderate complexity | LRU cache, word frequency |
| ★★★ Hard | Concurrency, complex algorithms | Producer-consumer, rate limiter |
