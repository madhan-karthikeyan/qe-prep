# Docker Basics for QE

## Dockerfile Best Practices

### Multi-Stage Builds

Separate build environment from runtime to produce small, secure images.

```dockerfile
# Stage 1: Build
FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/server

# Stage 2: Runtime
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /app/server /server
EXPOSE 8080
USER nobody
ENTRYPOINT ["/server"]
```

### .dockerignore

Prevents sending unnecessary files to the Docker daemon, speeding up builds.

```
node_modules
.git
*.log
.env
__pycache__
*.pyc
.venv
dist
```

### Layer Caching

Order matters. Place less-frequently-changing layers first (dependencies before code).

```dockerfile
# BAD — cache invalidated on every code change
COPY . .
RUN pip install -r requirements.txt

# GOOD — dependencies layer cached as long as requirements.txt doesn't change
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY . .
```

## docker-compose for Multi-Service Testing

```yaml
version: "3.9"
services:
  app:
    build: .
    ports: ["8080:8080"]
    environment:
      DB_URL: "postgres://user:pass@db:5432/testdb"
      REDIS_URL: "redis://redis:6379"
    depends_on:
      db: { condition: service_healthy }
      redis: { condition: service_healthy }

  db:
    image: postgres:16
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: testdb
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U user -d testdb"]
      interval: 2s
      retries: 10
    tmpfs: /var/lib/postgresql/data  # ephemeral for tests

  redis:
    image: redis:7-alpine
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 2s
      retries: 10
```

```bash
# Run tests against full stack
docker compose up -d
docker compose run app pytest tests/
docker compose down -v
```

## Container Health Checks

Always define health checks for services in test environments.

```yaml
healthcheck:
  test: curl -f http://localhost:8080/health || exit 1
  interval: 5s
  timeout: 3s
  retries: 5
  start_period: 10s
```

In Dockerfile:
```dockerfile
HEALTHCHECK --interval=5s --timeout=3s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1
```

## Debugging Containers

```bash
# Exec into a running container
docker exec -it <container> sh

# View logs
docker logs <container>
docker logs --tail=100 -f <container>

# Inspect config
docker inspect <container>

# Copy files out
docker cp <container>:/app/logs/error.log ./error.log

# Run a debug sidecar for network debugging
docker run -it --network container:<target-container> nicolaka/netshoot

# View resource usage
docker stats <container>
```

## Volume and Network Management

### Volumes

```yaml
services:
  app:
    volumes:
      - ./src:/app/src          # Bind mount (live code reload)
      - test_data:/data/test    # Named volume

volumes:
  test_data:
```

### Networks

```yaml
services:
  app:
    networks:
      - test_net
      - isolated_net

networks:
  test_net:
    driver: bridge
  isolated_net:
    internal: true  # no external access
```

**Important for tests:** Each test run should use a unique project name to avoid container naming conflicts in CI.

```bash
docker compose -p "test-${CI_JOB_ID}" up -d
```

## Resource Limits

```yaml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 256M
        reservations:
          cpus: '0.1'
          memory: 128M
```

Set limits to:
- Prevent one test suite from starving others in CI
- Reproduce resource-constrained environments
- Detect memory leaks early

## CI Integration

### GitHub Actions

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_PASSWORD: testpass
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
      - uses: actions/checkout@v4
      - run: docker compose up -d
      - run: docker compose run app make test
```

### Cleanup

```bash
# Always clean up after tests
docker compose down -v --remove-orphans

# Prune unused resources periodically (in CI cron)
docker system prune -af --volumes
```

### Tips

1. Use `--abort-on-container-exit` with `docker compose up` to stop when a test container exits
2. Use `--exit-code-from` to get exit codes from the test container
3. Use unique project names (`-p`) for parallel CI jobs
4. Use `tmpfs` mounts for test databases — faster and no cleanup needed
5. Pin image versions (don't use `latest`) to avoid unexpected changes in CI
