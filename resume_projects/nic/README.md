# NIC (Network Information Collector) — Network Diagnostics Tool

## Overview

NIC is a CLI-based network diagnostics tool that collects connectivity, latency, DNS resolution, and protocol-level information across target hosts. Designed for SRE and network engineering teams, it replaces ad-hoc shell scripting with a structured, extensible framework that produces consistent, machine-readable output for alerting, dashboards, and postmortem analysis.

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                     CLI Entrypoint                    │
│  (argparse, config file, env-based overrides)        │
└──────────┬──────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────┐
│                   Probe Scheduler                     │
│  (asyncio task graph with dependency resolution)     │
└──┬──────────┬──────────┬──────────┬─────────────────┘
   │          │          │          │
   ▼          ▼          ▼          ▼
┌──────┐ ┌──────┐ ┌──────┐ ┌──────────┐
│ ICMP │ │ TCP  │ │ DNS  │ │  HTTP    │
│ Ping │ │ Port │ │ Res. │ │  Probe   │
└──┬───┘ └──┬───┘ └──┬───┘ └────┬─────┘
   │        │        │          │
   ▼        ▼        ▼          ▼
┌─────────────────────────────────────────────────────┐
│                  Result Collector                     │
│  (aggregate, deduplicate, cache with TTL)            │
└──────────────────────┬──────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────┐
│                   Format Pipeline                     │
│  ┌─────────┐  ┌──────────┐  ┌───────────┐          │
│  │ JSON    │  │ YAML     │  │ Table     │          │
│  │ Writer  │  │ Writer   │  │ Renderer  │          │
│  └─────────┘  └──────────┘  └───────────┘          │
└─────────────────────────────────────────────────────┘
```

## Why This Design

The plugin architecture was chosen because network environments are heterogeneous — a probe that works in a cloud environment may be irrelevant for an on-prem data center. By decoupling probe discovery from execution, new collectors can be shipped independently without modifying the core scheduler or output pipeline. The format pipeline is similarly pluggable so that teams can pipe results directly into their existing monitoring stack (JSON for Prometheus exporters, YAML for config management tools, tables for human debugging).

## Key Technical Decisions

- **Plugin system via `importlib.metadata` entry points**: Probes are registered as Python entry points in `pyproject.toml`. The scheduler discovers all installed probes at startup. This means a probe is a standalone package — teams own their collectors, release them on their own cadence, and the core tool stays thin. Downside: dependency conflicts require careful namespace management.

- **Async probes with `asyncio`**: Network operations are I/O-bound. Using `asyncio.gather()` with a semaphore for concurrency control allows running 50 probes against 20 hosts in under 3 seconds. Each probe is a coroutine that yields a typed result dataclass. The scheduler handles timeout via `asyncio.wait_for()` so a single hung probe cannot stall the entire run.

- **Result caching with TTL**: Many probes (e.g., DNS resolution for the same host) produce identical results within a short window. A TTL-based LRU cache avoids redundant network calls. Cache keys are (probe_type, target, params_hash). The TTL is configurable per-probe via metadata, defaulting to 30 seconds for volatile data and 300 seconds for stable data like geolocation.

## Testing Strategy

- **Unit tests with mock network responses**: Each probe is tested against recorded network responses (fixtures stored as JSON files). The mock transport layer returns controlled failures, timeouts, and edge cases (empty responses, malformed packets). This catches parsing bugs without network dependencies.

- **Integration tests against localhost services**: A test harness spins up small Docker containers (nginx on random ports, a DNS stub resolver, a TCP echo server) and runs probes against them. Validates end-to-end collection, result structure, and cache hit/miss behavior.

- **End-to-end tests in Docker with containerized services**: A full stack test using `docker-compose` launches a simulated network topology (3 services, 1 unreachable host, 1 slow host). The CLI is invoked as a subprocess and output is validated across all three formats. This catches regressions in the format pipeline and scheduler.

- **Property-based testing for output formatting**: Using Hypothesis, generated random result structures are fed through each formatter. Properties enforced: valid JSON/YAML output, no duplicate keys, all fields present, timestamps in ISO 8601, latency values non-negative. Catches edge cases like None fields, NaN floats, and Unicode in hostnames.

## Failures & Lessons

- **Initial blocking I/O caused 30s delays on unreachable hosts**: The v1 prototype used `socket.connect()` and `subprocess.run()` for ping. A single unreachable host blocked the entire run for its timeout duration. Switched to `asyncio` with per-probe timeouts. Now a host that times out is reported within 2 seconds and the rest of the probes continue unaffected.

- **Plugin API v1 was too rigid**: The first plugin interface required every probe to implement `run()`, `parse()`, and `format()`. Teams found they needed hooks for pre/post conditions, dependency injection (e.g., shared DNS cache), and custom timeout logic. Redesigned with a versioned protocol: v1 is the minimal interface, v2 adds lifecycle hooks (`setup()`, `teardown()`, `on_timeout()`). The scheduler detects the protocol version at load time and adapts its execution path accordingly.

## Tradeoffs

- **Python speed vs developer productivity**: Python's GIL and interpreter overhead mean that packet-level operations (raw ICMP) require a C extension or `ctypes`. The alternative was Go or Rust, which would have eliminated the plugin ecosystem. Chose Python with C extensions for hot paths because debugging network issues is easier in a scripting language.

- **Async complexity vs blocking simplicity**: Async code is harder to reason about (callback context, task cancellation, shared state). Every probe author must understand coroutines. Mitigated by providing a `BaseProbe` abstract class that handles the async boilerplate so most probes only write synchronous-style logic.

- **Plugin flexibility vs API surface area**: A maximally flexible plugin system makes the core codebase harder to maintain. The versioned protocol (v1/v2) strikes a balance: most probes use v1, power users opt into v2. The scheduler's plugin loader is the most complex part of the codebase and the most likely source of bugs.

## Interview Questions

**Q: How would you add a new collector?**
A: Create a new Python package with a `pyproject.toml` that registers an entry point under `nic.probes`. Implement a class that inherits from `BaseProbe` and defines `probe_type`, `target`, and the async `run()` coroutine that returns a `ProbeResult`. Install the package. The scheduler discovers it automatically at next run. No core changes needed.

**Q: How do you handle a plugin that crashes?**
A: The scheduler wraps every probe execution in a `try/except` with `asyncio.wait_for`. If a probe raises, its result is recorded as an error result with the exception message, stack trace (truncated), and the probe's target metadata. The scheduler continues executing remaining probes. A crash counter is exposed via a health endpoint. If the same probe crashes N times consecutively, it is quarantined (skipped in future runs) and an alert fires.

**Q: How would you distribute this tool?**
A: As a PyPI package with optional extras per probe category: `pip install nic[all]` or `pip install nic[dns,http]`. For air-gapped environments, provide a Docker image with all built-in probes. The plugin system allows internal teams to host a private PyPI server or vend wheel files in a monorepo. CI publishes packages via GitHub Actions on tagged releases.

## Related Problems from This Repository

- **TCP Echo Server**: Shares the same async I/O pattern and timeout management approach used in TCP probe.
- **HTTP Client with Retry**: The retry/backoff logic directly informed NIC's probe retry policy for transient network failures.
- **Connection Pool**: Connection pooling in the HTTP probe was adapted from this repo's generic pool implementation.
- **URL Parser**: Used for parsing and normalizing target specifications across all probe types.
