# Resume Guide — QE Engineer

## Overview

Your resume must convey that you are a **technical quality engineer**, not a manual tester. Recruiters scan for keywords, metrics, and impact — not duties.

## Principles

1. **Every bullet is a result, not a task.**  
   ❌ "Responsible for writing test cases"  
   ✅ "Designed and automated 200+ test cases, reducing regression cycle from 3 days to 4 hours"

2. **Quantify everything.**  
   Use numbers: bugs found, time saved, coverage percentage, team size, release cycles.

3. **Match job description keywords.**  
   If the JD mentions "Kafka," "distributed systems," "pytest" — those words should appear in your resume (honestly).

4. **Show technical depth.**  
   Mention frameworks, languages, tools, protocols.

## Before & After Examples

### QA Engineer (Before — Weak)

```
QA Engineer — CompanyX (2019-2022)
- Wrote test cases for web application
- Performed manual testing
- Reported bugs in JIRA
- Attended daily standups
```

### QE Engineer (After — Strong)

```
Quality Engineering — CompanyX (2019-2022)
- Built automated test framework using Python + Selenium, covering 500+ test scenarios
- Reduced regression cycle from 5 days to 6 hours via parallel execution on 16-node cluster
- Authored 80+ bug reports; 95% confirmed by dev team, 3 critical P0 bugs caught pre-production
- Implemented CI/CD pipeline integration (Jenkins) with flaky test detection, reducing false failures by 40%
- Led performance testing initiative — identified 2 N+1 query regressions via k6 load tests
- Introduced mutation testing (MutPy) increasing test suite reliability score from 0.72 to 0.94
```

### Test Lead (Before — Weak)

```
Test Lead — CompanyY (2020-2023)
- Managed team of 3 testers
- Created test plans
- Reported to QA Manager
- Coordinated with dev team
```

### Test Lead (After — Strong)

```
Engineering Lead, Quality — CompanyY (2020-2023)
- Led 3-person QE team embedded across 4 dev pods; drove shift-left testing strategy
- Designed risk-based test plans for 12 major releases; zero P0/P1 production escapes in 2 years
- Reduced test suite execution time by 70% through test parallelization and containerization (Docker + k6)
- Established quality dashboards (Grafana) tracking escaped defect rate, MTTR, and test pass rate per service
- Championed chaos engineering program — ran weekly GameDays injecting pod/network failures; found 8 resilience gaps
- Mentored 2 junior engineers to independent test-design capability within 3 months
```

### Entry Level / Bootcamp Grad (Before — Weak)

```
QA Intern — CompanyZ (2022)
- Helped with testing
- Learned JIRA
- Attended meetings
```

### Entry Level (After — Strong)

```
QE Intern — CompanyZ (2022)
- Wrote 150+ automated API tests (Python + pytest) for order-management microservice
- Automated deployment smoke test reducing manual verification from 45 min to 3 min
- Documented 20 common failure scenarios with reproduction steps for dev onboarding
- Contributed to open-source testing tool (1 merged PR)
```

## Keywords That Catch Recruiters' Eyes

| Category | Keywords |
|----------|---------|
| **Languages** | Python, Java, Go, JavaScript, TypeScript, SQL |
| **Test Frameworks** | pytest, JUnit, TestNG, unittest, Playwright, Cypress, Selenium |
| **CI/CD** | Jenkins, GitHub Actions, GitLab CI, CircleCI, ArgoCD |
| **Performance** | k6, Locust, JMeter, Gatling, vegeta |
| **Infrastructure** | Docker, Kubernetes, Terraform, Helm, Ansible |
| **DB/Storage** | PostgreSQL, MySQL, Cassandra, Couchbase, MongoDB, Redis |
| **Observability** | Prometheus, Grafana, ELK, Datadog, OpenTelemetry |
| **Protocols** | HTTP, gRPC, WebSocket, TCP/IP, MQTT |
| **Concepts** | contract testing, property-based testing, chaos engineering, mutation testing, risk-based testing, shift-left |
| **Soft** | cross-functional collaboration, mentoring, stakeholder communication, triage facilitation |

## Metrics & Impact Framework

| What You Did | How to Quantify |
|-------------|----------------|
| Found bugs | "Found 45 bugs; 98% accepted, 5 critical" |
| Automated tests | "Automated 300 tests covering 85% of critical paths" |
| Reduced time | "Shortened regression from 3 days → 45 minutes" |
| Improved quality | "Reduced escaped defect rate from 5% → 0.5%" |
| Led initiative | "Led 3-person team delivering X" |
| Built framework | "Designed test framework used by 4 teams" |
| Prevented incidents | "Caught 2 P0 regressions pre-release" |
| Increased coverage | "Increased code coverage from 45% → 78%" |

## Project Description Template

```
[Project Name] — [Brief context, 1 sentence]

Tech: Python, pytest, PostgreSQL, Docker, Jenkins

- [Bullet: what YOU did, with metric]
- [Bullet: technical challenge overcome]
- [Bullet: impact with number]
```

### Example

```
Checkout System — E-commerce payment processing with 99.99% uptime requirement

Tech: Java, JUnit, WireMock, k6, Kubernetes

- Designed contract tests (Pact) for 3 payment providers; caught 2 integration mismatches before GA
- Wrote load tests simulating Black Friday traffic (10K req/s); identified and fixed 2 connection pool bottlenecks
- Reduced false-positive test failures by 60% via deterministic test data seeding and database snapshots
```

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Listing duties instead of impact | Rewrite every bullet as "action + metric" |
| Using passive voice | "Built framework" not "Framework was built" |
| Being vague | "Reduced time by 70%" not "Made testing faster" |
| No technical keywords | Include languages, tools, protocols explicitly |
| Too many pages | 1 page if <10 years exp, 2 pages max |
| Typos | Use Grammarly. Have a friend review. |

## Resume Sections Order

1. **Summary** (2 lines max) — role title + years + key expertise
2. **Skills** — keyword-based, grouped by category
3. **Experience** — reverse chronological, most impact first
4. **Projects** — if <3 years experience or relevant side work
5. **Education** — degree, school, year (omit GPA if >3 years ago)
6. **Certifications** — optional: ISTQB, AWS, k8s, etc.
