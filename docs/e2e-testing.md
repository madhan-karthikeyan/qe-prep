# End-to-End Testing

## E2E vs Integration vs Unit

| Aspect | Unit | Integration | E2E |
|--------|------|-------------|-----|
| **Scope** | Single function/class | Two or more components | Full system |
| **Speed** | Milliseconds | Seconds | Minutes |
| **Dependencies** | All mocked | Real DB/queue, mocks on 3rd party | Everything real (or closest replica) |
| **Flakiness** | Very low | Medium | High |
| **Cost to maintain** | Low | Medium | High |
| **Confidence** | Low (piece works) | Medium (pair works) | High (system works) |

**Test pyramid rule of thumb:** ~70% unit, ~20% integration, ~10% E2E.

## Full Workflow Testing

An E2E test should cover a complete user journey, not isolated UI interactions.

```python
def test_user_registration_and_checkout():
    # 1. Register
    signup("alice@example.com", "password123")
    assert email_contains_verification_link("alice@example.com")

    # 2. Verify email
    click_verification_link(get_verification_url("alice@example.com"))
    assert is_logged_in()

    # 3. Browse and add to cart
    search_product("wireless mouse")
    add_to_cart("Wireless Mouse MX")
    assert cart_badge_shows(1)

    # 4. Checkout
    proceed_to_checkout()
    fill_shipping("123 Main St")
    select_payment("credit_card")
    submit_order()
    assert confirmation_page_shows()
    assert email_received("alice@example.com", "order_confirmation")
```

## Environment Management

| Strategy | How | Pro | Con |
|----------|-----|-----|-----|
| **Dedicated staging** | Long-lived env | Stable, known state | Expensive, configuration drift |
| **Ephemeral env** | Spin up per PR/branch | Identical to prod, fresh state | Slower, resource-heavy |
| **Preview env** | Deploy from PR | Fast feedback | May not match prod scale |
| **Local dev env** | Docker Compose | Fastest iteration | Not realistic for cloud infra |

**Best practice:** Use ephemeral environments with preview URLs for each PR, and a dedicated staging environment for nightly E2E runs.

## Data Setup and Cleanup

```python
@pytest.fixture
def seed_data(api_client):
    # Create test user via API (not UI — faster, more reliable)
    user = api_client.post("/api/users", json={
        "email": "test-{uuid}@example.com",
        "password": "TestPass123!"
    })
    yield user
    # Cleanup via API
    api_client.delete(f"/api/users/{user['id']}")
```

- **Avoid** data setup through the UI in E2E tests — it's slow and fragile
- Use API calls or direct database insertions for test data
- Generate unique data per run (UUID suffixes) to avoid collisions
- Clean up after each test, even on failure

## Flakiness in E2E Tests

### Common Causes

| Cause | Mitigation |
|-------|-----------|
| Network timeouts | Use generous timeouts; retry on transient failures |
| Async rendering | Use explicit wait conditions, not `sleep()` |
| Shared state | Isolate tests; each test creates its own data |
| Browser differences | Run against a single browser in CI; expand later |
| Database contention | Unique data per test; transactions with rollback |

### Handling Flaky Tests

1. **Detect** — Track test pass rate over last 100 runs
2. **Quarantine** — Move flaky tests to a separate pipeline
3. **Fix** — Investigate root cause (not symptom)
4. **Re-integrate** — Only return fixed tests to main suite

## Page Object Model (UI)

Separate page structure from test logic for maintainability.

```python
class LoginPage:
    def __init__(self, page):
        self.page = page

    def navigate(self):
        self.page.goto("https://app.example.com/login")

    def login(self, email, password):
        self.page.fill("#email", email)
        self.page.fill("#password", password)
        self.page.click("#login-button")

    def error_message(self):
        return self.page.text_content(".error-message")

    def is_logged_in(self):
        return self.page.is_visible(".user-avatar")


# Test stays clean
def test_invalid_login(login_page):
    login_page.navigate()
    login_page.login("bad@example.com", "wrong")
    assert "Invalid credentials" in login_page.error_message()
```

## API E2E Testing Patterns

### Pattern 1: Chained API calls

```python
def test_order_lifecycle(api):
    # Create
    order = api.post("/orders", json={"item": "laptop"})
    assert order.status_code == 201
    order_id = order.json()["id"]

    # Process
    payment = api.post(f"/orders/{order_id}/pay", json={"method": "card"})
    assert payment.status_code == 200

    # Verify state
    status = api.get(f"/orders/{order_id}")
    assert status.json()["status"] == "paid"
```

### Pattern 2: State verification

```python
def test_payment_creates_invoice(api, db):
    api.post("/orders", json={"item": "book", "price": 20})
    invoice = db.query("SELECT * FROM invoices WHERE amount = 20")
    assert invoice is not None
    assert invoice["status"] == "pending"
```

### Pattern 3: Contract validation

```python
def test_create_user_returns_expected_schema(api):
    resp = api.post("/users", json={"name": "Alice"})
    assert resp.status_code == 201
    assert set(resp.json().keys()) == {"id", "name", "created_at", "links"}
```

### Best Practices for API E2E

1. Use a client that handles auth, base URL, and retries
2. Validate response schemas with tools like Pydantic, OpenAPI validators
3. Test error paths: 4xx, 5xx, malformed bodies, missing headers
4. Verify side effects (database state, email sent, event emitted)
5. Clean up created resources after test
