# QE Interview Playbook

[![CI](https://github.com/madhan/qe-interview-playbook/actions/workflows/ci.yml/badge.svg)](https://github.com/madhan/qe-interview-playbook/actions/workflows/ci.yml)
![Python](https://img.shields.io/badge/Python-3.12-3776AB)
![Go](https://img.shields.io/badge/Go-1.23-00ADD8)
![Java](https://img.shields.io/badge/Java-21-ED8B00)

A production-quality repository for **Software/Quality Engineering** interview preparation.

Focuses on real engineering tasks — logging, concurrency, networking, testing, debugging,
distributed systems — rather than pure algorithms.

---

## Repository Structure

```
qe-interview-playbook/
├── python/               ← Full implementations + tests (pytest)
├── go/                   ← Full implementations + tests (testing)
├── java/                 ← Full implementations + tests (JUnit 5)
├── c/                    ← Minimal reference implementations
├── cpp/                  ← Minimal reference implementations
├── patterns/             ← Reusable engineering patterns
├── debugging/            ← Interactive debugging exercises
├── bug_reports/          ← Bug report templates + examples
├── interview/            ← Round-specific interview guides
├── resume_projects/      ← Personal project deep-dives
├── benchmarks/           ← Performance benchmarks
├── distributed-systems/  ← Distributed systems simulations
├── fault_injection/      ← Chaos engineering scenarios
└── docs/                 ← Testing & engineering guides
```

---

## Quick Start

```bash
# Python
cd python && pip install -r requirements-dev.txt && pytest -v

# Go
cd go && go test -race -count=1 ./...

# Java
cd java && mvn test
```

---

## Coding Problems

| # | Module | Category | Python | Go | Java |
|---|--------|----------|--------|----|------|
| 1 | Logging Toolkit | Logging | ✓ | ✓ | ✓ |
| 2 | File Processing | File I/O | ✓ | ✓ | ✓ |
| 3 | LRU Cache | Data Structures | ✓ | ✓ | ✓ |
| 4 | Linear Data Structures | Data Structures | ✓ | ✓ | ✓ |
| 5 | Trie | Data Structures | ✓ | ✓ | ✓ |
| 6 | Producer-Consumer | Concurrency | ✓ | ✓ | ✓ |
| 7 | Thread Pool | Concurrency | ✓ | ✓ | ✓ |
| 8 | Rate Limiter | Concurrency | ✓ | ✓ | ✓ |
| 9 | Networking | Networking | ✓ | ✓ | ✓ |
| 10 | URL Parser | Parsing | ✓ | ✓ | ✓ |
| 11 | Expression Evaluator | Parsing | ✓ | ✓ | ✓ |
| 12 | Config Parser | Parsing | ✓ | ✓ | ✓ |
| 13 | Database Utilities | Database | ✓ | ✓ | ✓ |
| 14 | REST API Tester | Automation | ✓ | ✓ | ✓ |

---

## Documentation

| Guide | Description |
|-------|-------------|
| [QE Mindset](docs/qe_mindset.md) | Thinking like a quality engineer |
| [Unit Testing](docs/unit-testing.md) | FIRST principles, test structure |
| [Integration Testing](docs/integration-testing.md) | Real dependencies, testcontainers |
| [Debugging Techniques](docs/debugging-techniques.md) | Print vs debugger, profiling |
| [Concurrency Testing](docs/concurrency-testing.md) | Race detection, stress testing |
| [Distributed Systems Testing](docs/distributed-systems-testing.md) | Jepsen-style, fault injection |
| [Docker Basics](docs/docker-basics.md) | Containers, compose, multi-stage |
| [Bug Reports](docs/bug-reports.md) | Writing reproducible reports |
| [Profiling](docs/profiling.md) | cProfile, pprof, JFR |

---

## CI Pipeline

Each push runs:

- **Python**: ruff → mypy → pytest
- **Go**: gofmt → go vet → golangci-lint → go test -race
- **Java**: checkstyle → spotbugs → JUnit 5

---

## Prerequisites

- Python 3.11+
- Go 1.22+
- Java 21+
- Docker (for integration tests and fault injection)

---

## License

MIT
