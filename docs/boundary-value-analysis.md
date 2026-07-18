# Boundary Value Analysis

## Equivalence Partitioning

Divide inputs into groups (partitions) where each value in a group behaves the same way. Test one representative from each partition.

**Example:** A field accepts ages 18–65.

| Partition | Representative |
|-----------|---------------|
| Below range | 17 |
| In range | 30 |
| Above range | 66 |

No need to test 19, 20, 21... — they're all equivalent.

## Boundary Value Analysis (BVA)

Test the edges of equivalence partitions. Most bugs live at boundaries.

For the same age field (18–65), test:

| Boundary | Values | Why |
|----------|--------|-----|
| Lower boundary | 17, 18, 19 | Off-by-one below, exactly on, just above |
| Upper boundary | 64, 65, 66 | Just below, exactly on, just above |

That's 6 test cases instead of the 48 values in the entire range.

## Off-by-One Errors

The most common boundary bug. Occurs when using `<` instead of `<=` or vice versa.

```python
# Bug: allows values up to 100 but should be < 100
def is_valid(items):
    return len(items) < 100   # Bug: should be <= 100

# Fixed
def is_valid(items):
    return len(items) <= 100
```

**Interview tip:** Always write tests for `n-1`, `n`, and `n+1` where `n` is the boundary.

## Common Boundary Types

### Numeric

| Scenario | Boundaries |
|----------|------------|
| 0 < x < 100 | -1, 0, 1, 99, 100, 101 |
| x >= 18 | 17, 18, 19 |
| Discount tiers (10% off orders > $100) | $100, $100.01 |
| Array index | 0, 1, length-2, length-1, length |

### String Length

| Scenario | Boundaries |
|----------|------------|
| min=1, max=50 | "", len=1, len=50, len=51 |
| max=255 (DB field) | len=254, 255, 256 |
| Unicode/charset | ASCII, multi-byte UTF-8, emoji |

### Date Ranges

| Scenario | Boundaries |
|----------|------------|
| Start must be before end | same date, start=end+1, start=end-1 |
| Subscription period | expires today, expires tomorrow, expired yesterday |
| Leap year | Feb 28, Feb 29 (leap/non-leap), Mar 1 |
| Timezone boundaries | UTC midnight, DST transitions |

## Examples with Code

### Python

```python
import pytest

def calculate_discount(items):
    """10% off for 5+ items, 20% off for 10+ items"""
    if len(items) >= 10:
        return 0.20
    elif len(items) >= 5:
        return 0.10
    return 0.0

@pytest.mark.parametrize("count,expected", [
    (0, 0.0),
    (4, 0.0),
    (5, 0.10),   # boundary — exactly at threshold
    (9, 0.10),   # boundary — just below next threshold
    (10, 0.20),  # boundary — exactly at next threshold
    (100, 0.20),
])
def test_discount_boundaries(count, expected):
    items = [None] * count
    assert calculate_discount(items) == expected
```

### Java

```java
@ParameterizedTest
@CsvSource({
    "0, 0",
    "4, 0",
    "5, 10",
    "9, 10",
    "10, 20",
    "100, 20"
})
void testDiscountBoundaries(int count, int expectedPercent) {
    List<Item> items = Collections.nCopies(count, new Item("x", 1.0));
    assertEquals(expectedPercent, calculateDiscount(items));
}

int calculateDiscount(List<Item> items) {
    if (items.size() >= 10) return 20;
    if (items.size() >= 5) return 10;
    return 0;
}
```

### Go

```go
func TestDiscountBoundaries(t *testing.T) {
    tests := []struct {
        count    int
        expected int
    }{
        {0, 0},
        {4, 0},
        {5, 10},
        {9, 10},
        {10, 20},
        {100, 20},
    }
    for _, tt := range tests {
        t.Run(fmt.Sprintf("count_%d", tt.count), func(t *testing.T) {
            items := make([]Item, tt.count)
            result := calculateDiscount(items)
            if result != tt.expected {
                t.Errorf("got %d, want %d", result, tt.expected)
            }
        })
    }
}
```

## Practical Interview Tips

1. **Start with equivalence partitions** — divide input space, then pick boundaries.
2. **Test each boundary three times** — at, just below, just above.
3. **Don't forget output boundaries** — not just inputs. Does a function return values near its limits?
4. **Check off-by-one in loops** — `<` vs `<=`, `i++` vs `++i`.
5. **Watch for hidden boundaries** — string encoding, timezone transitions, floating point precision.
6. **Ask about constraints** — before writing tests, clarify min/max, null/empty, valid ranges.

**Standard interview response:** "I would partition the input space into equivalence classes, then test values at each boundary — n-1, n, and n+1 — since most bugs cluster at the edges."
