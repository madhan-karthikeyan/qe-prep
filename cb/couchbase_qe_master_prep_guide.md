# Couchbase QE Interview — Master Prep Guide
*For: Madhan Karthikeyan · VIT Vellore campus drive · Stream: QE*

---

## 0. Reality check: what the JD actually means

The JD you pasted is a shared DEV/SRE/QE posting — Couchbase writes one broad "core engineering" JD and splits candidates by stream afterward. Read "Go, Java, Python, and can hack in C/C++" as **breadth signaling**, not a mastery checklist. In practice, interviewers calibrate to fresher level and focus on:

1. **What's already on your resume** — fair game for deep grilling, since you claimed it (Python, Java, C, C++, SQL).
2. **What the specific QE team actually uses day to day** — worth building fresh, because it's both JD-stated *and* real: Couchbase's India-based QE Cloud team builds test automation in **TypeScript and Go**, using **Cypress**, against **AWS/Azure/GCP/Kubernetes/Couchbase Server**.
3. **Everything else on the wishlist** (deep C++, full Java mastery) — you'll get light conceptual questions at most, not production-level grilling.

So: don't panic-learn five languages to fluency in a week. Triage instead — the plan below is ordered by actual leverage.

---

## 1. Your current stack vs. the gap

| Language | On your resume? | Used in a project? | QE-team relevant? | Action |
|---|---|---|---|---|
| Python | Yes | Yes (all 3 projects, pytest) | Yes (also used in other Couchbase QE frameworks) | **Sharpen**, don't rebuild |
| TypeScript | Not listed, but... | Yes (React frontends in Apex Sprint Planner & JobLens) | **Yes — core QE stack (Cypress)** | **Extend into testing** |
| Java | Yes (skills list) | No | Lightly | **Refresh fundamentals** |
| C / C++ | Yes (skills list) | No | Lightly | **Refresh fundamentals** |
| Go | Not listed | No | **Yes — core QE stack** | **Build from near-zero** |
| SQL | Yes | Implied (Postgres/SQLAlchemy) | Yes (N1QL is SQL-flavored) | **Bridge to N1QL** |

Your honest priority order: **Python (deepen) → TypeScript/Cypress (extend) → Go (build) → Java/C/C++ (refresh, don't rebuild)**.

---

## 2. Per-language testing toolchain

This is the part most guides skip: *what tool do you reach for, in each language, to actually test something.* Know one or two per language, not the whole list.

### Python (your strongest — go deep here)
- **Unit/integration:** `pytest` (your main tool — you already have 363 tests in DecisionDrift, be ready to explain your fixture/parametrize choices), `unittest` (know it exists, stdlib alternative)
- **Mocking:** `unittest.mock`, `pytest-mock`
- **Property-based/fuzz testing:** `Hypothesis` — mentioning this in an interview is a strong QE signal (generates edge cases automatically instead of hand-writing them)
- **Coverage:** `coverage.py` / `pytest-cov`
- **API calls in tests:** `requests`, `httpx`
- **Load testing:** `Locust` (Python-native, good to name-drop given Couchbase cares about "heavy load and stressful conditions")
- **Concurrency:** `threading`, `multiprocessing`, `asyncio` — and know *why* the GIL matters for CPU-bound vs I/O-bound work

### TypeScript / JavaScript (leverage your existing React experience)
- **E2E/UI:** **Cypress** — this is explicitly named in Couchbase's QE job posting. Install it, run `cypress open`, write one spec (visit a page, assert an element, fill a form) before the interview. Know the difference between Cypress and Selenium (Cypress runs in-browser, no WebDriver, faster/flakier-resistant).
- **Alt worth recognizing:** Playwright (cross-browser, increasingly common)
- **Unit:** Jest or Vitest
- **API testing in Node:** Supertest

### Go (build this — genuinely new for you, but highest leverage)
- **Testing:** the built-in `testing` package + `go test` — Go's convention is **table-driven tests**, learn this pattern specifically, it's the first thing an interviewer will look for in your Go test code
- **Assertions/mocks:** `testify` (assert/require), `gomock`
- **HTTP testing:** `net/http/httptest`
- **Race detection:** `go test -race` — directly relevant since the JD explicitly calls out concurrent/multi-threaded programming as something they want you to find "cool"
- **Concurrency primitives to know:** goroutines, channels, `select`, `sync.Mutex`, `sync.WaitGroup`
- Minimum viable Go: syntax, error handling (`if err != nil` pattern), structs/interfaces, and one table-driven test written by hand

### Java (refresh, don't rebuild — it's already claimed on your resume)
- **Testing:** JUnit 5 (know `@Test`, `@BeforeEach`, assertions), TestNG (recognize it, less critical)
- **Mocking:** Mockito
- **API testing:** RestAssured
- **UI:** Selenium WebDriver
- **Concurrency:** `Thread`, `ExecutorService`, `synchronized`, `java.util.concurrent` — at least know these exist and what problem each solves
- Priority: OOP fundamentals (inheritance, interfaces, collections — `List`/`Map`/`Set` internals) over anything exotic

### C / C++ (refresh — "hack in" level, not mastery)
- **Testing:** Google Test (gtest) is the standard; Catch2 as an alternative
- **Memory/concurrency debugging:** Valgrind (memory leaks), **ThreadSanitizer** (race detection — again, ties directly to the concurrency callout in the JD), gdb for live debugging
- **Build:** Make/CMake basics
- Priority: pointers, manual memory management (malloc/free, new/delete), and being able to explain *why* a race condition or memory leak happens — not writing a full C++ project from scratch

---

## 3. Cross-cutting QE toolchain (language-agnostic)

| Category | Tools to know | Priority |
|---|---|---|
| API testing | Postman/Insomnia, Newman (CLI runner) | High — install and use before interview |
| Load/performance | JMeter, Locust, k6 | Medium — know one, name-drop the concept |
| CI/CD | GitHub Actions (you know this), Jenkins (recognize it) | High — you're already strong here |
| Containers | Docker, Docker Compose (you know this), Kubernetes/`kubectl` basics | High — Kubernetes is named in the QE stack |
| Test/bug management | JIRA, TestRail (conceptual awareness only) | Low |
| Version control | Git, branching/PR workflows | Already solid |
| Database/query | Couchbase Query Workbench, **N1QL/SQL++** | High — see below |

**Do this before the interview:** spin up Couchbase Capella's free tier or a local Docker instance, load a sample bucket, and run a basic N1QL query. It's a 20-minute investment that lets you say "I've actually used your product" — which almost no fresher candidate can say.

---

## 4. Distributed systems & concurrency — the JD calls these out by name

The JD explicitly says they want people who think distributed systems and concurrent/multi-threaded programming are "cool." That's not filler — Couchbase's entire product *is* a distributed database, and QE's job is finding where it breaks under concurrency. Know these well enough to discuss, not just define:

- **CAP theorem** — and be ready to say what Couchbase actually prioritizes (it's tunable — strong vs. eventual consistency depending on config)
- **Replication & partitioning/sharding** — why data gets split across nodes, what happens when a node fails
- **Consistency models** — strong vs. eventual consistency, and a real example of when each matters
- **Race conditions** — what one looks like in code, how you'd write a test to catch it (this connects directly to `go test -race` / ThreadSanitizer above)
- **Quorum reads/writes** — the idea of "enough nodes agree" before a write/read is considered valid
- **Idempotency & retries** — why retryable operations need to be safe to run twice (you already touched this in your Celery+Redis pipeline — have that story ready)
- **Testing flaky/non-deterministic systems** — this is a QE-specific challenge worth having an opinion on: how do you write a test for something that only fails 1 in 100 runs?

---

## 5. Your strongest unused card: the AI/vector search angle

Couchbase has spent the last year pushing hard into AI — vector search across their whole product line, a new "AI Data Plane" for agent memory (announced June 2026), Hyperscale Vector Index for billion-scale similarity search. This is a current company priority, not a side feature.

Your NIC internship — a VLM pipeline extracting structured fields and matching them via cosine similarity over 150+ embedded documents, with a custom ingestion/chunking/embedding service — **is a small-scale version of exactly what Couchbase's vector search product does.** Almost no other fresher candidate walks in with hands-on embedding/retrieval pipeline experience.

**Use this deliberately:**
- In "tell me about a project" — lead with the NIC internship, and explicitly connect it: *"I built a retrieval pipeline over embedded documents with metadata-filtered search — which is conceptually close to what your vector search / AI Data Plane does at a much larger scale."*
- In "why Couchbase" — this is a genuine, specific answer instead of a generic one.
- If asked about databases + AI intersection (the JD literally says "intersection of AI, Databases and Data Processing") — you have a real story, not a hypothetical.

---

## 6. Question types to expect, by category

Based on a confirmed Couchbase-Vellore interview report (Jul 2025) plus Couchbase's standing QE job postings — treat category 1–2 as fairly confirmed, the rest as standard SDET/QE convention:

1. **One DSA question** — moderate difficulty. Code it in Python (your fastest language) unless told otherwise. Arrays/strings/hashing fundamentals.
2. **Networking & security basics** — HTTP vs HTTPS, HTTP status code categories, what SQL injection and XSS are and how you'd design a test case to catch each.
3. **Test-design scenario** — "design test cases for [a booking site / login flow / an API]." Structure your answer by test type: functional, boundary/negative, integration, API, performance/load, security, and — given Couchbase's domain — **concurrency** (e.g., "what if two users book the same seat/write the same key simultaneously?").
4. **Language fundamentals** — likely Python in depth (your strongest), lighter Java/C/C++ conceptual questions since they're resume-claimed, possibly a light Go question if the interviewer wants to test adaptability.
5. **Testing concept questions** — test pyramid (unit/integration/e2e), mocking vs. stubbing, what makes a test flaky and how you'd fix it, equivalence partitioning / boundary value analysis / decision tables as test-case design techniques.
6. **System design lite** — "design a key-value store" style questions have come up for other Couchbase roles; have a basic point of view (hashing, partitioning, replication trade-offs).
7. **Deep resume grilling** — expect this specifically; a past candidate's HR round explicitly involved close questioning against the resume. Know cold: why DecisionDrift is deterministic (no LLM in the enforcement path) and what that trade-off buys you, how your 21-patch/95.2% Recall@5 benchmark was built, why Celery+Redis with retries, what your ranking weights in JobLens actually optimize for.
8. **Live tool demo** — your placement notice explicitly said tools must be installed and working; expect to write/run a small test live (pytest or Cypress) or hit an endpoint in Postman.
9. **Behavioral** — why QE and not core dev, why Couchbase, a teamwork/conflict example.

---

## 7. Study priority plan (given limited time)

**Tier 1 — must-do, highest leverage:**
- [ ] Deepen pytest fluency — be ready to defend every design choice in DecisionDrift's 363 tests
- [ ] Write and run one Cypress spec (leverages your existing React/TypeScript experience directly)
- [ ] DSA refresh — arrays, strings, hashing, one or two medium LeetCode-style problems/day
- [ ] Practice structuring test-case-design answers out loud (functional/boundary/security/concurrency framing)
- [ ] Prep your resume deep-dive narrative, especially the "why deterministic, not LLM" story
- [ ] Prep the AI/vector-search bridge story from your NIC internship

**Tier 2 — high-value, build from near-zero:**
- [ ] Go syntax + the `testing` package + one table-driven test written by hand
- [ ] Goroutines/channels — just enough to explain a producer-consumer example
- [ ] Docker Desktop confirmed working; can build/run a container live
- [ ] Postman installed, comfortable sending requests and writing assertions
- [ ] `kubectl` basics — enough to not freeze if asked

**Tier 3 — refresh, don't rebuild (already resume-claimed):**
- [ ] Java: OOP fundamentals, collections (`List`/`Map`/`Set`), one JUnit test written by hand
- [ ] C/C++: pointers, manual memory management, one gtest example

**Tier 4 — nice-to-have if time remains:**
- [ ] Locust or k6 basics for load testing
- [ ] Spin up Couchbase Capella free tier, run one N1QL query
- [ ] ThreadSanitizer/Valgrind — know what they catch even if you don't run them live

---

## 8. Day-of checklist

- [ ] Laptop charged + charger packed
- [ ] IDE (VS Code recommended) open, extensions working, no unresolved errors
- [ ] Docker Desktop running (test with `docker run hello-world` beforehand)
- [ ] Git configured, GitHub accessible
- [ ] Postman installed and logged in if needed
- [ ] Cypress installed in a scratch project (`npm install cypress`, confirm `cypress open` works)
- [ ] Node.js and Python environments both working (`node -v`, `python --version`)
- [ ] Resume printed or accessible — know every number on it without looking

---

*Caveat: the tool stack here is sourced from Couchbase's live QE job postings and one confirmed Vellore campus interview report, not a leaked script for your specific batch. Treat it as strong preparation, not a guarantee of exact questions.*
