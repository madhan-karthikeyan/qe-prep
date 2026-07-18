# The QE Mindset

## Introduction: What Distinguishes a QE Engineer from a Software Engineer

A software engineer asks: *"How do I build this feature?"*
A QE engineer asks: *"How will this feature break?"*

The difference is not about writing tests — it's about how you think about software. The QE mindset is a systematic approach to finding critical bugs before they reach users, prioritizing test effort where it matters most, and driving quality across the entire development lifecycle.

This guide captures that mindset. It's the most important guide in this playbook.

---

## 1. Thinking Beyond Happy Paths

The happy path is the shortest, most common, and most tested path through a system. A QE engineer systematically explores everything *around* the happy path.

**Systematic exploration techniques:**

| Category | Questions to Ask |
|----------|-----------------|
| **Empty state** | What happens when there's no data? No users? No items? |
| **Error state** | What happens when a dependency fails? A timeout? A crash? |
| **Boundary state** | What happens at min/max values? Exceeding limits? |
| **Edge case** | Unicode? Null characters? Extremely long strings? |
| **Race condition** | Two operations at the same time? Out of order? |
| **Resource exhaustion** | Full disk? No memory? No file descriptors? |
| **Security** | Unauthenticated access? SQL injection? XSS? |

**Mental exercise:** For every feature, ask "What are the three most likely things that will go wrong?" Then test those.

---

## 2. Finding Edge Cases

### Boundary Analysis

Most bugs cluster at boundaries. Always test:
- One below, exactly at, one above every boundary
- Zero, one, and many for collections
- Empty string, single character, max length, overflow

### Combinatorial Testing

When multiple parameters interact, test combinations that are most likely to fail.

```python
# Instead of testing all 3^4 = 81 combinations of {A,B,C} × {X,Y,Z} × ...
# Use pairwise testing to cover all pairs in fewer tests
# Tools: AllPairs, PICT, ACTS

from allpairspy import AllPairs

parameters = [
    ["Chrome", "Firefox", "Safari"],
    ["Windows", "macOS", "Linux"],
    ["Free", "Premium"],
    ["IPv4", "IPv6"],
]

for combo in AllPairs(parameters):
    # Tests each pair of values at least once
    test_login(*combo)
```

### State Transitions

Systems with state (state machines, workflows, multi-step forms) need testing at every transition.

```
                      +--→ [Paid] --+
                      |              ↓
[Draft] --→ [Pending] --→ [Shipped] --→ [Delivered]
                |              |
                ↓              ↓
           [Cancelled]    [Returned]
```

**Test every arrow.** Test forbidden transitions (can you ship a cancelled order?). Test double transitions (click "Ship" twice).

---

## 3. Risk-Based Testing

You cannot test everything. Prioritize where the risk is highest.

**Risk = Probability × Impact**

| Risk Level | Probability | Impact | Test Effort |
|------------|-------------|--------|-------------|
| **Critical** | High | Data loss, security breach, revenue impact | Maximum — automated + manual exploratory |
| **High** | Medium | Major feature broken, slow degradation | High — automated integration + E2E |
| **Medium** | Low | Edge case, cosmetic issue | Moderate — unit tests |
| **Low** | Very low | Rare corner case, minor UX | Minimal — skip or smoke test |

**Process:**
1. Identify all features and changes
2. Assess probability and impact for each
3. Prioritize test effort by risk level
4. Review and adjust as the project evolves
5. Document what you chose NOT to test and why

---

## 4. Failure Injection Mindset

Ask: *"How can I make this break?"* Then deliberately do it.

**Physical layer:** Kill processes, fill disks, disconnect networks
**Application layer:** Pass null, send malformed data, call endpoints out of order
**Timing layer:** Slow down responses, introduce delays, send requests simultaneously
**State layer:** Corrupt data, roll back transactions, modify database directly

**Example:**

```python
def test_how_does_system_handle_corrupted_cache():
    # Inject corrupted data into Redis
    redis.set("session:abc123", ":::garbage:::not:::json:::")
    # Now try to use the session — should it crash? Error gracefully? Recreate?
    response = client.get("/dashboard", headers={"Cookie": "session=abc123"})
    # The answer tells you about the system's resilience
    assert response.status_code == 200  # or 401? or 500?
```

---

## 5. Reproducibility

A bug you cannot reproduce is a bug you cannot fix.

**Making flaky tests deterministic:**
1. Remove shared state — each test creates its own data
2. Inject deterministic clocks instead of `datetime.now()`
3. Mock random number generators with a fixed seed
4. Use synchronization barriers for async code
5. Replace network calls with controlled stubs

**Capture the full context on failure:**

```python
@pytest.fixture
def fail_with_context(request):
    yield
    if request.node.rep_call.failed:
        # Save DB state, logs, thread dump, network captures
        save_debug_artifacts(request.node.name)
```

---

## 6. Automation First

If it can be automated, it should be automated.

| What | Automate? | Why |
|------|-----------|-----|
| Regression tests | ✅ Always | Run on every PR to catch regressions instantly |
| Smoke tests | ✅ Always | Verify deployment succeeded |
| Performance benchmarks | ✅ Always | Detect regressions before release |
| UI exploratory testing | ❌ Never | Human creativity finds novel bugs |
| Security scanning | ✅ Always | Run in pipeline, not manually |
| Environment setup | ✅ Always | `docker compose up`, `terraform apply` |

**Rule:** If you run a manual test more than twice, automate it.

---

## 7. Test Pyramid

```
        ╱╲
       ╱ E2E ╲         ← Few (10%) — critical user journeys
      ╱────────╲
     ╱Integration╲      ← Some (20%) — service boundaries, DB, external APIs
    ╱──────────────╲
   ╱   Unit Tests    ╲  ← Many (70%) — fast, isolated, reliable
  ╱────────────────────╲
```

**Where to invest:**
- **Unit tests** — best ROI for catching logic bugs. Fast feedback, zero flakiness.
- **Integration tests** — catches real interaction bugs (SQL, serialization, contracts).
- **E2E tests** — catches system-level bugs but slow and flaky. Use sparingly.

**Anti-pattern:** Inverted pyramid with too many E2E tests. Expensive to maintain, slow, and unreliable.

---

## 8. Flaky Tests

A flaky test passes and fails without code changes. It erodes trust in the test suite.

### Detection

```bash
# Run tests repeatedly to find flakiness
pytest --count=20 tests/ --random-order
go test -count=20 -race ./...
```

**Quarantine workflow:**
1. Detect flaky test (fails intermittently in CI)
2. Move to a separate "quarantine" suite
3. Investigate root cause — don't just re-run
4. Fix the underlying issue (not the test — the code or the test design)
5. Return to main suite only after 50+ consecutive passes

### Root Causes & Fixes

| Cause | Fix |
|-------|-----|
| Time-dependent | Inject clock, use deterministic timeouts |
| Async race | Add synchronization, use `await` properly |
| Shared data | Fresh data per test |
| Network dependency | Mock or use controlled test server |
| Resource leak | Proper cleanup in `finally`/`defer` |
| Order dependency | Run tests in random order, fix isolation |

---

## 9. Regression Strategy

**What to retest:**
- All code paths that changed (direct impact)
- All code paths that interact with changed code (indirect impact)
- All critical user journeys (always, even if unrelated)

**How often:**

| Cadence | Scope | Who |
|---------|-------|-----|
| Every PR | Unit + integration for changed modules | CI |
| Every merge to main | Full unit + integration + smoke E2E | CI |
| Nightly | Full E2E + performance benchmarks | CI |
| Pre-release | Full regression suite + exploratory | QE team |

**Smart regression:** Only run tests affected by the code change.

```python
# Tools: pytest-testmon, pytest-changed
pytest --testmon  # only re-runs tests affected by changed code
```

---

## 10. Bug Advocacy

Getting a bug fixed is a skill. A well-reported bug with impact analysis gets fixed.

**Strategy:**
1. **Write a great report** (see `bug-reports.md`) — clear steps, expected vs actual, environment
2. **Quantify impact** — "This blocks 5% of registrations, costing ~$1K/day"
3. **Provide a failing test** — a red test is harder to ignore than a bug report
4. **Classify correctly** — don't cry wolf with false P0s
5. **Follow up** — after 2 sprints, escalate with updated data
6. **Celebrate fixes** — thank the developer publicly

**What NOT to do:**
- Don't assign blame ("Your code is broken")
- Don't exaggerate severity (you lose credibility)
- Don't report without reproducing first
- Don't let bugs rot in the backlog — either escalate or close

---

## 11. Testing in Production

Production is the ultimate test environment. But it must be done safely.

| Technique | Description |
|-----------|-------------|
| **Canary deployments** | Roll out to 1% of users, monitor, then ramp |
| **Feature flags** | Toggle features on/off without deployment |
| **Shadow traffic** | Duplicate real requests to a new system without affecting users |
| **Chaos engineering** | Inject failures in production (with guardrails) |
| **Synthetic monitoring** | Automated probes that simulate user behavior 24/7 |
| **Real user monitoring (RUM)** | Collect metrics from actual user browsers |

**Key metrics to monitor:**
- Error rate (5xx, 4xx, client-side JS errors)
- Latency (p50, p95, p99)
- Throughput (requests/sec)
- Business metrics (conversion, signup rate)

**Golden rule:** Build your production monitoring before you release, not after.

---

## 12. Collaboration with Developers

### Shift-Left

Move quality activities earlier in the development cycle.

| Activity | Traditional (Reactive) | Shift-Left (Proactive) |
|----------|----------------------|----------------------|
| Requirements | QE sees spec after dev starts | QE reviews requirements for testability |
| Design | QE sees design doc | QE reviews with failure injection in mind |
| Code | QE tests after merge | QE reviews PRs, suggests testable patterns |
| Testing | QE runs all tests | Dev writes unit tests, QE focuses on integration/E2E |

### Code Review

As a QE, focus on these aspects in code reviews:
- **Test coverage** — are edge cases tested? Are error paths covered?
- **Test design** — are tests readable? Do they follow FIRST?
- **Assertions** — is there an actual assertion, or just smoke testing?
- **Mocking** — are mocks replacing something that should be real?
- **Testability** — is the code easy to test? Dependency injection? Interfaces?

### Building Trust

- Offer to pair with developers on writing tests
- Share test patterns and tools that make testing faster
- When you find a bug, fix the test AND suggest a code fix
- Be the person who finds the hard bugs, not the nitpicker

---

## 13. Continuous Learning

Stay current. The tools change fast.

| Area | What to Learn |
|------|---------------|
| **Languages** | Deepen one language, learn basics of 2–3 others |
| **Testing frameworks** | pytest, JUnit 5, Go testing, Cypress/Playwright |
| **Observability** | OpenTelemetry, Prometheus, Grafana, structured logging |
| **Infrastructure** | Docker, Kubernetes, Terraform basics |
| **Security** | OWASP Top 10, SAST/DAST tools, dependency scanning |
| **Performance** | Profiling (CPU, memory), load testing (k6, Locust) |
| **Chaos engineering** | Litmus, Chaos Mesh, Gremlin |

**One thing per month:** Pick one tool or concept and invest 5 hours.

---

## Summary: The QE Mindset Checklist

Before releasing any feature, ask yourself:

- [ ] Have I tested the happy path?
- [ ] Have I tested every boundary (n-1, n, n+1)?
- [ ] Have I tested every error state (dependency failure, timeout, invalid input)?
- [ ] Have I tested empty and extreme states (zero data, max data)?
- [ ] Have I tested concurrent operations?
- [ ] Have I tested forbidden state transitions?
- [ ] Have I tried to break it deliberately (failure injection)?
- [ ] Is the test result deterministic and reproducible?
- [ ] Is the test automated and in CI?
- [ ] Does the test suite run in < 10 minutes?
- [ ] Are flaky tests quarantined?
- [ ] Are bugs reported with clear reproduction steps and impact?
- [ ] Have I verified in production (canary, monitoring)?

**The QE mindset is not a role — it's how you approach software. Anyone can adopt it. The best engineers do.**
