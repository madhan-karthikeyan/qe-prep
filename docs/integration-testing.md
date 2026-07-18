# Integration Testing

## What to Integration Test

| Component | Examples | Why |
|-----------|----------|-----|
| **Database** | CRUD operations, migrations, constraint enforcement | SQL dialects, unique constraints, transactions behave differently from mocks |
| **Network** | REST/gRPC calls, message queues | Serialization, timeouts, retries, status codes |
| **File System** | File upload/download, config file parsing | Permissions, encoding, path resolution |
| **External APIs** | Payment gateways, auth providers | Contract mismatches, rate limits, error responses |

## Testcontainers

Testcontainers spins up real dependencies (PostgreSQL, Redis, Kafka) in disposable Docker containers for tests.

**Python (testcontainers):**
```python
from testcontainers.postgres import PostgresContainer

def test_user_repository():
    with PostgresContainer("postgres:16") as postgres:
        db_url = postgres.get_connection_url()
        repo = UserRepository(db_url)
        repo.create_table()
        repo.save(User("alice@example.com"))
        assert repo.find_by_email("alice@example.com") is not None
```

**Java (Testcontainers):**
```java
@Testcontainers
class UserRepositoryTest {
    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:16");

    @Test
    void testSaveAndFindUser() {
        UserRepository repo = new UserRepository(postgres.getJdbcUrl());
        repo.createTable();
        repo.save(new User("alice@example.com"));
        assertNotNull(repo.findByEmail("alice@example.com"));
    }
}
```

**Go (testcontainers-go):**
```go
func TestUserRepository(t *testing.T) {
    ctx := context.Background()
    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: testcontainers.ContainerRequest{
            Image: "postgres:16",
            ExposedPorts: []string{"5432/tcp"},
        },
        Started: true,
    })
    require.NoError(t, err)
    defer container.Terminate(ctx)

    port, _ := container.MappedPort(ctx, "5432")
    repo := NewUserRepository(fmt.Sprintf("postgres://localhost:%s/...", port.Port()))
    // test here
}
```

## Docker in Tests

- Use Testcontainers (don't manage containers manually in test code)
- Set resource limits on test containers (memory, CPU) to avoid host exhaustion
- Use a reusable network for containers that need to communicate
- Always terminate containers in `defer` / `@AfterEach` / context managers

## External Dependency Management

| Strategy | Best For | Tradeoff |
|----------|----------|----------|
| **Testcontainers** | Databases, message brokers, caches | Slower than mocks but real behavior |
| **WireMock / MockServer** | HTTP APIs | Fast, controllable responses, but drift from real API |
| **Contract tests (Pact)** | Service-to-service integrations | Catches contract breaks, but needs infra |
| **Sandbox/Test accounts** | Payment gateways, 3rd party APIs | Requires network, rate-limited |

## Setup and Teardown

```python
@pytest.fixture(scope="module")
def database():
    container = PostgresContainer("postgres:16")
    container.start()
    url = container.get_connection_url()
    run_migrations(url)
    yield url
    container.stop()

@pytest.fixture(autouse=True)
def clean_db(database):
    truncate_all_tables(database)
    yield
```

```java
@BeforeEach
void setUp() {
    jdbcTemplate.execute("TRUNCATE TABLE users CASCADE");
}

@AfterAll
static void tearDown() {
    postgresContainer.stop();
}
```

```go
func setup(t *testing.T) *UserRepository {
    t.Helper()
    container := startPostgres(t)
    t.Cleanup(func() { container.Terminate(context.Background()) })
    repo := NewUserRepository(container.ConnectionString)
    repo.RunMigrations()
    return repo
}
```

## Real vs In-Memory Dependencies

| Aspect | Real (Testcontainers) | In-Memory (H2 / SQLite) |
|--------|----------------------|------------------------|
| **Fidelity** | High — same SQL, same constraints | Low — dialect differences, missing features |
| **Speed** | Slow (container startup) | Fast (in-process) |
| **Setup complexity** | Requires Docker | Zero infra |
| **Flakiness** | Docker daemon issues | Deterministic |
| **CI support** | Needs Docker-in-Docker | Works everywhere |

**Rule of thumb:** Use Testcontainers for critical data paths (payments, auth). Use in-memory for quick feedback during development. Never use in-memory as a full replacement for the real database in CI.

## Patterns by Language

### Python (pytest + testcontainers)
```python
def test_order_flow():
    with RabbitMqContainer("rabbitmq:3") as rmq:
        with PostgresContainer("postgres:16") as db:
            app = create_app(db_url=db.get_connection_url(), broker_url=rmq.get_connection_url())
            client = app.test_client()
            resp = client.post("/orders", json={"item": "book"})
            assert resp.status_code == 201
```

### Go (testcontainers-go)
```go
func TestOrderFlow(t *testing.T) {
    dbContainer := startPostgres(t)
    redisContainer := startRedis(t)
    defer dbContainer.Terminate(ctx)
    defer redisContainer.Terminate(ctx)

    app := NewApp(dbContainer.URL, redisContainer.URL)
    resp := app.CreateOrder(Order{Item: "book"})
    assert.Equal(t, 201, resp.StatusCode)
}
```

### Java (Spring Boot + Testcontainers)
```java
@SpringBootTest
@Testcontainers
class OrderFlowTest {
    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:16");

    @DynamicPropertySource
    static void configure(DynamicPropertyRegistry reg) {
        reg.add("spring.datasource.url", postgres::getJdbcUrl);
    }

    @Test
    void testCreateOrder() {
        var resp = restTemplate.postForEntity("/orders", new Order("book"), String.class);
        assertThat(resp.getStatusCode()).isEqualTo(201);
    }
}
```
