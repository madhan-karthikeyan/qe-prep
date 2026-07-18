# Container crashes on startup with non-root user

**Severity:** Major
**Priority:** P1
**Environment:** Docker 20.10+, Kubernetes 1.22+, all Linux distros
**Component:** Containerization / Dockerfile

## Summary

The application container exits immediately on startup when running as a non-root user (UID 1001). The app tries to write logs to `/var/log/app/`, which is owned by `root:root` (UID 0) in the container image. Since `PodSecurityPolicy` and security best practices require non-root containers, this breaks deployments in production.

## Steps to Reproduce

1. Build the Docker image using the provided Dockerfile
2. Run container with `--user 1001` (or deploy on Kubernetes with `runAsUser: 1001`)
   ```bash
   docker build -t myapp:latest .
   docker run --rm --user 1001 myapp:latest
   ```
3. Observe container exit immediately

## Expected Behavior

Container should start successfully and the application should run correctly as a non-root user (UID != 0).

## Actual Behavior

Container exits with error:

```
Error: Failed to create log directory /var/log/app/: Permission denied
```

The exit code is 1 (generic error). The container status shows `CrashLoopBackOff` in Kubernetes.

## Logs / Screenshots

```
$ docker run --rm --user 1001 myapp:latest
Error: Failed to create log directory /var/log/app/: Permission denied

$ docker run --rm myapp:latest
# (runs fine as root)

$ docker run --rm -it --user 1001 myapp:latest /bin/sh
$ ls -la /var/log/app/
drwxr-xr-x 2 root root 4096 /var/log/app/   # ← owned by root
$ touch /var/log/app/test.log
touch: /var/log/app/test.log: Permission denied
```

Kubernetes pod events:
```
Events:
  Type     Reason     Age   From               Message
  ----     ------     ----  ----               -------
  Warning  BackOff    5s    kubelet            Back-off restarting failed container
  Normal   Pulling    12s   kubelet            Pulling image "myapp:latest"
  Warning  Unhealthy  5s    kubelet            Startup probe failed: stat /tmp/healthz: permission denied
```

## Root Cause Analysis

The Dockerfile creates the log directory but **does not grant write permissions** to non-root users:

```dockerfile
FROM alpine:3.18

RUN mkdir -p /var/log/app

COPY myapp /usr/local/bin/myapp

USER 1001
CMD ["myapp"]
```

Issues:
1. `mkdir -p /var/log/app` creates the directory owned by `root:root` (default)
2. `chmod` is not called — the directory has default permissions (755), which only allows the **owner** (root) to write
3. The application needs to create/write log files in this directory
4. The `USER 1001` directive only changes the runtime user — it does **not** fix existing directory ownership

This is an extremely common Dockerfile mistake. The fix is to create the directory **before** switching users and grant appropriate permissions.

## Fix

```dockerfile
FROM alpine:3.18

# Create directory and set ownership BEFORE switching user
RUN mkdir -p /var/log/app && \
    chown -R 1001:1001 /var/log/app && \
    chmod 755 /var/log/app

COPY myapp /usr/local/bin/myapp

USER 1001
CMD ["myapp"]
```

Additional hardening:

```dockerfile
FROM alpine:3.18

RUN addgroup -S myapp && \
    adduser -S -G myapp -u 1001 myapp && \
    mkdir -p /var/log/app && \
    chown -R myapp:myapp /var/log/app && \
    chmod 750 /var/log/app

COPY --chown=myapp:myapp myapp /usr/local/bin/myapp

USER myapp
CMD ["myapp"]
```

**Key changes:**
- Added `chown` to set log directory ownership to the runtime user
- Used named user (`myapp`) instead of raw UID for readability
- `--chown` flag on COPY ensures application binary is also owned correctly
- `chmod 750` restricts log access to the user and group only (security best practice)

## Regression Tests

### 1. Container Health Check with Non-Root User

```python
import subprocess
import docker

def test_container_runs_as_non_root():
    client = docker.from_env()
    
    # Build and run with non-root user
    container = client.containers.run(
        "myapp:latest",
        user="1001",
        detach=True,
        remove=True,
        entrypoint=["/bin/sh", "-c", "sleep 5 && stat /tmp/healthz"]
    )
    
    try:
        result = container.wait(timeout=10)
        logs = container.logs().decode()
        
        assert result["StatusCode"] == 0, f"Container exited with {result['StatusCode']}: {logs}"
        assert "Permission denied" not in logs
    finally:
        container.remove(force=True)

def test_container_creates_log_as_non_root():
    client = docker.from_env()
    
    container = client.containers.run(
        "myapp:latest",
        user="1001",
        detach=True,
        remove=True,
        entrypoint=["touch", "/var/log/app/test.log"]
    )
    
    result = container.wait(timeout=10)
    assert result["StatusCode"] == 0, f"Cannot create log file as non-root"
```

### 2. Dockerfile Best Practice Check (Static Analysis)

```python
def test_dockerfile_chowns_log_directory():
    with open("Dockerfile") as f:
        content = f.read()
    
    # Check that the log directory is chowned
    assert "chown" in content.lower() or "chmod" in content.lower(), \
        "Dockerfile must set ownership/permissions for log directory"
    
    # Check that USER is not root
    user_lines = [l for l in content.splitlines() if l.startswith("USER")]
    assert len(user_lines) > 0, "Dockerfile must specify a non-root USER"
    assert "root" not in user_lines[-1].lower(), "USER must not be root"
```

### 3. CI Integration

Add to CI pipeline:
- Build Docker image
- Run container with `--user 1001:1001` and verify process is running after 5s
- Run container as non-root and verify log file creation succeeds
- Use `hadolint` to lint Dockerfile for security issues (missing USER, directory permissions)
- Run Kubernetes `pod-security-admission` pod test with `restricted` profile
