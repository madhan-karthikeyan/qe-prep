# Mocking & Stubbing

## When to Mock vs When to Use Real Objects

| Use Mock When… | Use Real Object When… |
|----------------|----------------------|
| External service (API, email, SMS) | Simple data structures (DTOs, configs) |
| Slow dependency (database, file I/O) | Fast, deterministic utils (string parsing, math) |
| Non-deterministic (random, clock, network) | Value objects (Money, Date, Email) |
| Hard to set up (third-party SDK) | Objects with no external dependencies |
| Expensive to create (ML model, large cache) | Simple repositories (can use fake/in-memory) |

**Rule of thumb:** Mock across process boundaries, not within your own code.

## Mocking Frameworks

### Python: `unittest.mock`

```python
from unittest.mock import Mock, patch

# Stub — return a fixed value
def test_get_user_returns_user():
    mock_db = Mock()
    mock_db.query.return_value = {"id": 1, "name": "Alice"}
    service = UserService(mock_db)
    assert service.get_user(1).name == "Alice"

# Mock — verify interaction happened
def test_send_notification_called():
    mock_notifier = Mock()
    service = NotificationService(mock_notifier)
    service.send_welcome("alice@example.com")
    mock_notifier.send.assert_called_once_with(
        "alice@example.com", "Welcome!"
    )

# Spy — wrap a real object
def test_logger_records_calls():
    logger = Logger()
    spy = Mock(wraps=logger)
    spy.info("hello")
    spy.info.assert_called_once_with("hello")

# Patch — replace a module's dependency
@patch("myapp.services.send_email")
def test_signup_sends_email(mock_send):
    signup("alice@example.com")
    mock_send.assert_called_once()
```

### Go: `gomock`

```go
// Generate mock: mockgen -source=notifier.go -destination=mock_notifier.go -package=notifier
type Notifier interface {
    Send(to, message string) error
}

func TestSendWelcome(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockNotifier := NewMockNotifier(ctrl)
    mockNotifier.EXPECT().Send("alice@example.com", "Welcome!").Return(nil)

    service := NewNotificationService(mockNotifier)
    service.SendWelcome("alice@example.com")
}
```

### Java: Mockito

```java
// Stub
@Test
void testGetUser() {
    UserRepository mockRepo = mock(UserRepository.class);
    when(mockRepo.findById(1)).thenReturn(new User("Alice"));

    UserService service = new UserService(mockRepo);
    User user = service.getUser(1);
    assertEquals("Alice", user.getName());
}

// Verify interaction
@Test
void testSendNotification() {
    Notifier mockNotifier = mock(Notifier.class);
    NotificationService service = new NotificationService(mockNotifier);

    service.sendWelcome("alice@example.com");

    verify(mockNotifier).send("alice@example.com", "Welcome!");
}

// Argument matchers
verify(mockNotifier).send(anyString(), contains("Welcome"));
```

## Interface-Based Design for Testability

Depend on abstractions, not concretions. This lets you swap real implementations with mocks.

```python
# Hard to test — concrete dependency
class OrderService:
    def __init__(self):
        self.db = PostgreSQLConnection()  # fixed

# Testable — depends on abstraction
class OrderService:
    def __init__(self, db: Database):  # any Database works
        self.db = db
```

```go
// Go interfaces are implicit — makes mocking natural
type Database interface {
    Save(order Order) error
    Find(id string) (Order, error)
}

type OrderService struct {
    db Database
}
```

```java
public class OrderService {
    private final OrderRepository repository;

    public OrderService(OrderRepository repository) {  // interface
        this.repository = repository;
    }
}
```

## Over-Mocking Antipattern

**Signs you are over-mocking:**
- Tests break when you refactor internal implementation
- Tests mock every single method call, even internal ones
- Setup code is longer than the test logic itself
- You mock value objects or simple data structures
- Tests pass but the real system fails

**Fix:** Use real objects for internal collaborators, integration tests for database/network code, and only mock across module/process boundaries.

| Symptom | Solution |
|---------|----------|
| Mocking `User.name` getter | Use a real `User` object |
| Mocking `String.format()` | Just call the real method |
| Mocking a private helper | Test through public API |
| 50-line mock setup | Write an integration test instead |

## Mocking External Services

### HTTP Services

```python
# Python: responses library
@responses.activate
def test_payment_gateway():
    responses.add(
        responses.POST,
        "https://payments.example.com/charge",
        json={"status": "success"},
        status=200,
    )
    result = PaymentService.charge(100)
    assert result.success
```

```go
// Go: httptest
func TestPaymentGateway(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"success"}`))
    }))
    defer server.Close()

    service := NewPaymentService(server.URL)
    result, _ := service.Charge(100)
    assert.True(t, result.Success)
}
```

```java
// Java: WireMock
@Test
void testPaymentGateway() {
    wireMockServer.stubFor(post(urlEqualTo("/charge"))
        .willReturn(aResponse()
            .withStatus(200)
            .withBody("{\"status\":\"success\"}")));

    PaymentService service = new PaymentService(wireMockServer.baseUrl());
    assertTrue(service.charge(100).isSuccess());
}
```

### Message Queues

Mock the publisher/subscriber interface and assert messages are sent/received with correct content.
