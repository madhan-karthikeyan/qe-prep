# Testing Round — QE Engineer Interview Guide

## Overview

The testing round assesses how you think about quality holistically. You'll be asked to design tests for a feature, critique existing tests, explain your automation strategy, and advocate for bug fixes against pushback. The interviewer wants to see structured thinking, risk awareness, and communication skills.

## Top 20 Questions & How to Answer Them

### Test Planning

1. **How would you test a file upload feature?**
   - **Structure**: Normal (valid file, various sizes, formats) → Edge (empty file, max size, special chars in name) → Error (network timeout, disk full, corrupted file) → Security (path traversal, oversized file, malware) → Performance (1K concurrent uploads, large files).
   - **Hint**: Mention smoke test boundary (1 byte, 1 MB, 10 GB) and resume-ability.

2. **Design a test strategy for a chat application.**
   - **Expected**: Message delivery (at-least-once, at-most-once), ordering, offline messages, group chat, media sharing, typing indicators. Include: unit → integration → E2E → performance → chaos.

3. **How do you decide what to automate vs what to test manually?**
   - **Expected**: Automate high-ROI: regression, smoke, performance, data validation. Keep manual: exploratory, usability, accessibility, one-time scenarios. Rule: automate if you'll run it >3 times.

4. **What is risk-based testing? Walk through an example.**
   - **Expected**: Identify features by risk (impact × probability). Test highest-risk items first. E.g., for a payment system: charge flow = critical risk (P0), receipt generation = lower risk (P2).

5. **How would you test a system that processes 1M events/day?**
   - **Expected**: Load test (peak throughput), stress test (2x-10x load), durability test (does data survive restart?), latency SLO verification, backpressure handling, schema evolution.

### Automation Strategy

6. **Describe your ideal test automation pyramid for a microservices architecture.**
   - **Expected**: 
     - **Unit** (60%): Service logic, domain models, utils
     - **Integration** (25%): API contracts, database queries, message queues
     - **E2E** (10%): Critical user journeys
     - **Chaos/Performance** (5%): Fault injection, load testing

7. **How do you handle flaky tests?**
   - **Expected**: Detect (rerun with `--flaky-fail`), quarantine, investigate root cause (race condition, timing, environment dependency), fix or delete. Never ignore flaky tests.

8. **What's your approach to API contract testing?**
   - **Expected**: Use Pact or OpenAPI-based contract tests. Consumer-driven contracts ensure provider changes don't break consumers. Run in CI on both sides.

9. **How do you test asynchronous systems (queues, event streams)?**
   - **Expected**: Deduplication testing (send same message twice), ordering tests (sequence numbers), timeout/dead-letter queue tests, message loss with fault injection, consumer rebalancing.

10. **How do you measure test coverage meaningfully?**
    - **Expected**: Line coverage is a floor, not a goal. Track mutation coverage (pitest/stryker), condition coverage, boundary coverage. Monitor coverage of changed code in PRs.

### Bug Advocacy

11. **A developer says "that edge case will never happen." How do you respond?**
    - **Expected**: "I agree it's unlikely, but the impact is high (data loss/crash). Let's measure: how hard is the fix? If it's 2 lines and a test, we should do it. If it's expensive, we can document and add monitoring."

12. **How do you prioritize which bugs to fix before release?**
    - **Expected**: Severity × Likelihood × User Impact. Create a triage matrix. Critical/P0 must fix; Major/P1 should fix; Minor/P2 fix if time; Trivial/P3 backlog or never.

13. **How do you communicate a release-blocking bug to stakeholders?**
    - **Expected**: Clear summary, reproduction steps, impact scope (users affected, data risk). Provide workaround, ETA if known. Use data: "This affects 30% of logins."

14. **What do you do if a bug is reproducible but developers can't find the cause?**
    - **Expected**: Narrow down preconditions (bisect configs, reduce input). Add logging. Capture thread dumps, heap dumps, network traces. Simplify reproduction to minimal code. Pair debug with developer.

15. **How do you track quality metrics across releases?**
    - **Expected**: Track escaped bug rate, MTTR (mean time to resolve), test pass rate, code coverage delta, flaky test count. Present as a dashboard, not a spreadsheet.

### Presenting Your Testing Strategy

16. **Walk through how you'd test a feature from spec to deployment.**
    - **Expected**:
      1. Review spec — identify ambiguities, missing scenarios
      2. Write test plan — risk-based, prioritized
      3. Write unit + integration tests alongside development
      4. Manual exploratory testing on feature branch
      5. E2E tests on staging
      6. Canary testing in production
      7. Post-release monitoring (dashboards, error budgets)

17. **What's the difference between positive and negative testing? Give examples.**
    - **Positive**: Valid inputs produce expected outputs. E.g., username "alice" is accepted.
    - **Negative**: Invalid inputs are properly rejected. E.g., username "alice\n" triggers validation error.

18. **How do you test a recommendation engine?**
    - **Expected**: Historical replay (run against past data), A/B testing framework, diversity metrics, cold-start scenarios, bias detection, performance (latency P99 < 100ms).

19. **What property-based testing tools have you used?**
    - **Expected**: Hypothesis (Python), QuickCheck (Haskell/Erlang), fast-check (JS). Write properties: "sorted(list) ≡ list but ordered", "JSON serialization roundtrips".

20. **How would you test a distributed transaction (2PC/Saga)?**
    - **Expected**: Coordinator crash after prepare, participant crash during commit, timeout handling, idempotency, compensating transactions for rollback. Inject failures at each phase.

---

## Answer Framework

| Step | Content |
|------|---------|
| 1 | Restate the feature/system to test (confirm understanding) |
| 2 | Identify key quality attributes (correctness, performance, security, reliability) |
| 3 | List risk areas — what could go wrong? |
| 4 | Propose test levels (unit → integration → E2E → performance → chaos) |
| 5 | Describe specific test cases (include normal, edge, failure) |
| 6 | Discuss tooling and automation approach |

---

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Only mentioning positive test cases | Always include edge and failure cases |
| Not quantifying coverage | "We test the high-risk paths" not "we have 80% coverage" |
| Ignoring non-functional testing | Mention performance, security, chaos |
| Being too vague | Specific test cases show depth of thinking |
| Not discussing prioritization | Not all tests are equal — explain tradeoffs |

## Difficulty Levels

| Topic | Difficulty |
|-------|-----------|
| Test planning basics | ★☆☆ |
| Automation strategy | ★★☆ |
| Bug advocacy | ★★☆ |
| Property-based testing | ★★★ |
| Distributed systems testing | ★★★ |
| Chaos engineering | ★★★ |
