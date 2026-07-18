# Unit Testing

## FIRST Principles

| Letter | Principle       | Meaning |
|--------|----------------|---------|
| **F**  | **Fast**       | Tests should run quickly (milliseconds). Slow tests discourage frequent runs. |
| **I**  | **Isolated**   | Tests should not depend on each other. Each test sets up and cleans up its own state. |
| **R**  | **Repeatable** | Same test, same environment, same result — every time. No flakiness. |
| **S**  | **Self-validating** | Tests output pass/fail. No manual interpretation of logs. |
| **T**  | **Timely**     | Tests should be written at the right time (ideally before or alongside production code). |

## Arrange-Act-Assert

```
// Arrange — set up the object and inputs
// Act     — call the method under test
// Assert  — verify the result
```

**Python (pytest):**
```python
def test_calculate_total():
    cart = ShoppingCart()
    cart.add_item(Item("shirt", 25.0))
    cart.add_item(Item("shoes", 50.0))
    total = cart.calculate_total()
    assert total == 75.0
```

**Go (testing):**
```go
func TestCalculateTotal(t *testing.T) {
    cart := NewShoppingCart()
    cart.AddItem(Item{Name: "shirt", Price: 25.0})
    cart.AddItem(Item{Name: "shoes", Price: 50.0})
    total := cart.CalculateTotal()
    if total != 75.0 {
        t.Errorf("expected 75.0, got %f", total)
    }
}
```

**Java (JUnit 5):**
```java
@Test
void testCalculateTotal() {
    ShoppingCart cart = new ShoppingCart();
    cart.addItem(new Item("shirt", 25.0));
    cart.addItem(new Item("shoes", 50.0));
    double total = cart.calculateTotal();
    assertEquals(75.0, total, 0.001);
}
```

## Test Doubles

| Type      | Purpose |
|-----------|---------|
| **Dummy** | Passed around but never used. Fills parameter lists. |
| **Fake**  | Working implementation but simplified (e.g., in-memory database). |
| **Stub**  | Returns canned answers to calls made during the test. |
| **Mock**  | Pre-programmed with expectations about which calls will be made. Fails if unexpected calls occur. |
| **Spy**   | Wraps a real object and records calls. Lets you verify happened after the fact. |

## Code Coverage Myths

| Myth | Reality |
|------|---------|
| "80% coverage means 80% of bugs are caught." | Coverage measures which lines ran, not whether they were *checked* correctly. |
| "100% coverage is the goal." | 100% coverage can still have untested behaviors, race conditions, or missing edge cases. |
| "Coverage guarantees quality." | Without assertions and meaningful inputs, high coverage is worthless. |
| "Low coverage means bad code." | Some code (e.g., error handling, config) is hard to test — coverage is just one signal. |

**Focus on:** mutation testing and boundary analysis rather than chasing a coverage percentage.

## Testing Private Methods

Don't test private methods directly — test them through the public interface. If a private method is so complex it needs its own tests, it's a sign it should be extracted into its own class/module.

```python
# BAD: testing private method
def test__calculate_discount():
    result = Order()._calculate_discount(100)
    assert result == 10

# GOOD: test through public interface
def test_order_total_applies_discount():
    order = Order(items=[Item(100)])
    assert order.total == 90  # discount applied internally
```

## Parameterized Tests

**Python (pytest):**
```python
import pytest

@pytest.mark.parametrize("input,expected", [
    ("racecar", True),
    ("hello", False),
    ("", True),
    ("a", True),
    ("ab", False),
])
def test_is_palindrome(input, expected):
    assert is_palindrome(input) == expected
```

**Java (JUnit 5):**
```java
@ParameterizedTest
@CsvSource({
    "racecar, true",
    "hello, false",
    "'', true",
    "a, true"
})
void testIsPalindrome(String input, boolean expected) {
    assertEquals(expected, isPalindrome(input));
}
```

**Go:** Go does not have built-in parameterized tests, but table-driven tests are the idiomatic alternative.

```go
func TestIsPalindrome(t *testing.T) {
    tests := []struct {
        input    string
        expected bool
    }{
        {"racecar", true},
        {"hello", false},
        {"", true},
        {"a", true},
    }
    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            result := IsPalindrome(tt.input)
            if result != tt.expected {
                t.Errorf("IsPalindrome(%q) = %v, want %v", tt.input, result, tt.expected)
            }
        })
    }
}
```

## Flaky Tests

### Causes

| Cause | Fix |
|-------|-----|
| Shared mutable state | Each test gets fresh objects |
| Time-dependent logic | Inject clocks, use deterministic timeouts |
| Network calls | Mock/stub external services |
| Resource leaks | Proper teardown in `finally` / `defer` / context managers |
| Randomness | Seed your random number generator in tests |
| Order dependency | Run tests in random order locally to detect |
| Async race conditions | Use synchronization barriers, await all futures |

### Detection

- Run tests multiple times: `pytest --count=10` or `go test -count=10`
- Tag flaky tests and quarantine them
- Track in CI: rerun failed tests automatically; alert if a test is flaky across builds

### Fix

1. Identify the root cause (isolation failure? time dependency? async?)
2. Fix at the source — don't add retries to the test
3. Verify fix by running the test 50–100 times in CI
