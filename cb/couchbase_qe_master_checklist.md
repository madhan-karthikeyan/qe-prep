# Couchbase QE — Master Topic Checklist
*Companion to your 58-hour hourly plan. That doc tells you WHEN. This one tells you WHAT, exhaustively, per subject — and hands you a question bank instead of answers where the whole point is that you generate the answers.*

---

## 0. How the two documents fit together

| This checklist's section | Where it slots into your hourly plan |
|---|---|
| Part A — System Design/LLD | Sun 8:00–10:15 + 10:15–12:15 |
| Part B — Test-Suite Design (the big one) | Sun 1:15–3:15 |
| Part C — DSA | Sat 9:00–11:00 |
| Part D — CN | Sun 3:30–5:00 |
| Part E — DBMS/SQL→N1QL | Sun 3:30–5:00 (fold in) or a Sat evening gap |
| Part F — OS | Sat/Sun gaps, or fold into Sun 3:30–5:00 block |
| Part G — OOP | Fold into Sat 7:15–9:15 Java block or Sun LLD block |
| Part H — Concurrency | Sat 4:30–6:30 |

Resume numbers, behavioral prep, day-of logistics, and per-language execution targets are **already fully covered** in your hourly plan (§5.4, §5.5, §8, §4) — not repeated here to avoid drift between two versions of the same thing.

**How to use this doc:** don't read it. Attempt it. Cover a box, say the answer out loud, uncover, check yourself, tick it.

---

## Part A — System Design / LLD

### A1. Universal checklist — run this for *every* design prompt, no exceptions
- [ ] Restate the problem in your own words before touching pen/keyboard
- [ ] Ask 1–2 clarifying questions (scale? read-heavy or write-heavy? single machine or distributed? persistence needed?)
- [ ] State assumptions explicitly if no clarification is given
- [ ] Separate functional requirements from non-functional (latency, scale, durability)
- [ ] Define core entities/data model first
- [ ] Sketch the API/interface signatures before the internals
- [ ] Draw a box-and-arrow diagram, even crude, on paper
- [ ] State time/space complexity of the core operations
- [ ] Name at least one explicit trade-off ("I chose X over Y because...")
- [ ] Name at least one failure mode and how the design handles it
- [ ] If asked to extend/scale it up, have a next step ready (sharding, caching, replication)

### A2. Confirmed Couchbase-reported prompts — do these first, cold, <15 min each
1. Key-value store — how would you store and look up keys optimally?
2. Parking lot — classic OOD: classes, relationships, edge cases (full lot, multiple vehicle types)
3. Hash table from scratch — collision handling strategy + when/how you'd resize
4. `malloc`/`free` — what does a simple allocator actually do internally?
5. `tail -f` on a file much larger than RAM — how do you avoid loading the whole file?
6. Parallel sort across a multi-core machine, memory-constrained

### A3. Adjacent likely prompts (same difficulty tier, common at infra/DB companies — not confirmed for your batch, but worth one pass each if A2 feels solid)
7. LRU cache (hash map + doubly linked list — this one doubles as a DSA problem, see Part C)
8. Rate limiter (token bucket vs sliding window — pick one, defend it)
9. Bounded connection pool (thread-safety: what happens when all connections are checked out?)
10. In-memory pub/sub or simple event bus
11. URL shortener (hashing + collision handling — good warm-up if the above feel done)

### A4. Self-grading — after each attempt, ask yourself
- [ ] Did I state a complexity number, not just "it's fast"?
- [ ] Did I name a trade-off out loud, unprompted?
- [ ] Did I mention what breaks at 10x scale?
- [ ] Did I avoid silence for more than ~10 seconds anywhere?

---

## Part B — Test-Suite Design (your real differentiator — go deep here)

### B1. The universal test-category framework
Memorize this list well enough that you run through it automatically, in order, for *any* feature — this ordering itself is a signal to the interviewer that you have a system, not just instincts.

| # | Category | The trigger question you ask yourself |
|---|---|---|
| 1 | Functional | Does it do the basic job correctly, happy path? |
| 2 | Boundary/edge values | What happens exactly at the limit — zero, max, one-off-max? |
| 3 | Negative/error handling | What happens on wrong, missing, or malformed input? |
| 4 | Integration | What if the thing this depends on is slow, down, or returns garbage? |
| 5 | API-level | Status codes, headers, payload shape, versioning, idempotency |
| 6 | Performance/load | What happens at 10x or 100x expected traffic? |
| 7 | Security | AuthN/authZ, injection, data exposure, abuse/rate limits |
| 8 | **Concurrency** | What if two (or more) actors do this at the *exact same time*? — Couchbase's signature angle, never skip this one |
| 9 | Data consistency/state | What's the actual DB/system state afterward, not just the HTTP response? |
| 10 | Recovery/resilience | What if it fails halfway through — partial write, crash mid-request? |
| 11 | Observability | If this broke in production, how would anyone find out? |

### B2. Practice scenario bank — attempt each cold, out loud, ~10–12 min, running the *entire* B1 list against it

**Auth & identity**
- **Q1.** Design a complete test suite for a login flow with username/password + MFA. What are all the categories of test cases you can think of?
- **Q2.** Design a test suite for a "forgot password" / reset-password flow.
- **Q3.** Design a test suite for session/token expiry and refresh (think JWT specifically — what's special about testing token expiry vs a session cookie?).
- **Q4.** Design a test suite for role-based access control on an admin panel (e.g., a support agent shouldn't see billing data).

**Transactional / business-critical**
- **Q5.** Design a test suite for a train/flight seat booking system — specifically, what happens when two users try to book the *same seat* at the *same time*?
- **Q6.** Design a test suite for an e-commerce checkout/payment flow.
- **Q7.** Design a test suite for a shopping cart where stock is shared across concurrent users.
- **Q8.** Design a test suite for a voting/polling system — how do you prevent (and test for) a double vote?

**API & data**
- **Q9.** Design a test suite for a generic REST endpoint, e.g. `POST /orders`.
- **Q10.** Design a test suite for a paginated list API over a large dataset (what breaks at page 1 vs the last page vs a page that no longer exists because data changed mid-scroll?).
- **Q11.** Design a test suite for a file upload feature.
- **Q12.** Design a test suite for a file download/streaming feature.
- **Q13.** Design a test suite for a webhook delivery system (hint: at-least-once delivery — what does that imply about the receiver's test cases?).

**Infra / distributed-systems flavored — extra weight here, it's Couchbase's actual domain**
- **Q14.** Design a test suite for a distributed cache with TTL and eviction.
- **Q15.** Design a test suite for database replication/failover — does a write survive a node going down mid-request? How would you even construct that test?
- **Q16.** Design a test suite for a rate limiter under sudden burst traffic.
- **Q17.** Design a test suite for a search feature (full-text or vector search) — think relevance *and* freshness (newly indexed doc should become searchable within some SLA — how do you test that?).
- **Q18.** Design a test suite for multi-tenant data isolation — one tenant should never see another's data, even under load.

### B3. Self-practice protocol (do this, don't just read the list)
1. Pick one Q, set a 10-minute timer.
2. Say the B1 categories out loud in order. Force yourself to generate **at least 2 concrete cases per category** before moving to the next.
3. When time's up, grade yourself: which categories did you skip or rush? Which felt automatic?
4. Track your weak categories across attempts — for almost everyone this is **concurrency** and **security**, since they're the least intuitive to generate on the fly. That tracking *is* your real signal for what to drill more, not raw scenario count.
5. Bonus rep: for Q5, Q6, Q7, or Q15, after listing cases, pick your single best concurrency case and actually explain **how you'd implement it** as an automated test (two threads/processes hitting the same resource, a lock or transaction to assert against).

---

## Part C — DSA

### C1. Pattern coverage — check off as "can solve without notes and state complexity unprompted"
- [ ] Two pointers
- [ ] Sliding window
- [ ] Hash map / frequency counting
- [ ] Binary search (+ search-on-answer variants)
- [ ] BFS/DFS graph traversal
- [ ] Topological sort
- [ ] DP — 1D and 2D
- [ ] Heap / priority queue
- [ ] Intervals / merge
- [ ] Backtracking
- [ ] Trie
- [ ] Union-find
- [ ] Linked list manipulation
- [ ] Stack / monotonic stack

### C2. Curated practice set — one representative problem per pattern
Pick your own exact problem per pattern from LeetCode/GFG if you already have go-tos; if not, these are solid defaults:
1. Two Sum / Two Sum II (two pointers)
2. Longest Substring Without Repeating Characters (sliding window)
3. Top K Frequent Elements (hashmap + heap)
4. Search in Rotated Sorted Array (binary search variant)
5. Number of Islands (BFS/DFS)
6. Course Schedule (topological sort)
7. Word Break (DP)
8. Kth Largest Element in a Stream (heap)
9. Merge Intervals
10. Subsets / Combination Sum (backtracking)
11. Implement Trie (autocomplete-style — ties into "search" from Part B)
12. Number of Provinces / Redundant Connection (union-find)
13. Detect Cycle in a Linked List
14. **LRU Cache** (deliberately overlaps with Part A — do this one properly, it's a two-for-one)
15. **Design a Hit Counter** or **Design Twitter feed** (OOD + DSA hybrid — good if a round blends the two, which yours reportedly does)

### C3. Execution checklist per problem
- [ ] State the brute-force approach first, out loud, before jumping to optimal
- [ ] State time/space complexity before *and* after optimizing
- [ ] Narrate reasoning even when stuck — silence is the actual failure mode, not a wrong first attempt
- [ ] Once code is written, dry-run it against your own edge cases before declaring done — this is a QE round, showing this habit unprompted is worth real points

---

## Part D — Computer Networks

### D1. Topic checklist
- [ ] OSI vs TCP/IP model layers
- [ ] TCP vs UDP — and when you'd choose each
- [ ] TCP three-way handshake
- [ ] HTTP vs HTTPS, what TLS actually buys you
- [ ] HTTP methods and which are idempotent (GET/PUT/DELETE) vs not (POST)
- [ ] HTTP status code categories (1xx/2xx/3xx/4xx/5xx) and what triggers 401 vs 403
- [ ] DNS resolution flow, end to end
- [ ] Load balancing algorithms (round robin, least connections, consistent hashing)
- [ ] REST vs WebSocket vs gRPC — when each fits
- [ ] Cookies vs sessions vs JWT
- [ ] CORS — what it protects against
- [ ] Reverse proxy vs forward proxy
- [ ] Connection pooling / keep-alive, and why it matters for a DB client

### D2. Sample questions to rehearse out loud
1. Walk me through what happens when you type a URL into a browser and hit enter.
2. TCP vs UDP — if you were designing a replication protocol between two DB nodes, which would you lean toward and why?
3. Why does idempotency matter when a client retries a failed request in a distributed system?
4. What's the difference between a 401 and a 403? Give an example where an API should return each.
5. How would you test an endpoint for SQL injection? For XSS?
6. Explain CORS, and how you'd test for a misconfigured CORS policy.
7. What problem does consistent hashing solve in load balancing, specifically for a distributed cache?

---

## Part E — DBMS / SQL → N1QL

### E1. Topic checklist
- [ ] ACID properties, with a concrete example of what breaks if each is violated
- [ ] Isolation levels and the anomaly each one prevents (dirty read, non-repeatable read, phantom read)
- [ ] Indexing — B-tree vs hash index, when each wins
- [ ] Normalization vs denormalization trade-offs
- [ ] Join types (inner/left/right/full) and a case where the wrong one silently drops data
- [ ] Transactions and rollback behavior
- [ ] Replication — leader-follower basics
- [ ] Sharding/partitioning strategies
- [ ] CAP theorem, and where a document store like Couchbase typically sits (tunable consistency)
- [ ] Eventual vs strong consistency — a scenario where "eventual" is actually fine

### E2. Sample questions
1. Explain ACID with an example of a transaction that could violate each property if not handled correctly.
2. What are the SQL isolation levels, and what specific anomaly does each one prevent?
3. How would you test a database replication setup for consistency *after* a node failure? (Directly relevant to a distributed-DB company.)
4. Where does CAP theorem apply to a system like Couchbase, and what does "tunable consistency" mean in practice?
5. Write a query (SQL or N1QL-style) to find duplicate rows/documents in a dataset.
6. How would you design a test to catch an index performance regression, not just a correctness bug?

---

## Part F — Operating Systems

### F1. Topic checklist
- [ ] Process vs thread — memory model, creation cost, communication method
- [ ] CPU scheduling algorithms (basic overview — FCFS, round robin, priority)
- [ ] Deadlock — four necessary conditions, plus one prevention strategy
- [ ] Race condition vs deadlock — the actual difference
- [ ] Paging vs segmentation
- [ ] Virtual memory basics
- [ ] Mutex vs semaphore vs monitor
- [ ] Context switch cost — why it matters for a highly concurrent server
- [ ] System calls — what a syscall actually is at a mechanical level

### F2. Sample questions
1. What's the difference between a process and a thread, in terms of memory and communication?
2. What are the four conditions required for deadlock? How would you go about detecting one in a running system?
3. Explain paging vs segmentation.
4. Give an example of a race condition and an example of a deadlock — make sure they're clearly distinct.
5. At a high level, how does an OS scheduler decide what runs next?

---

## Part G — OOP

### G1. Concept checklist
- [ ] Encapsulation, abstraction, inheritance, polymorphism — one clean example of each, not textbook definitions
- [ ] Composition over inheritance — a case where composition clearly wins
- [ ] SOLID principles — one sentence each
- [ ] Interface vs abstract class
- [ ] At least one design pattern you can code from memory (Singleton, Factory, Observer, Strategy, Builder)

### G2. Sample questions — tie these to your own projects where you can
1. Explain SOLID with an example from something you've actually built (e.g., how DecisionDrift's rule categories reflect single-responsibility, or how JobLens's ingestion sources are open for extension).
2. When would you choose composition over inheritance? Give a real example, not a hypothetical.
3. Design a parking lot with proper class relationships (this doubles as Part A's OOD prompt — do it here focusing purely on the class design, then again in Part A focusing on the system-design framing).
4. Explain Singleton and why it's risky in multi-threaded code. How would you actually *test* that a Singleton implementation is thread-safe?

---

## Part H — Concurrency (cross-language)

### H1. Checklist
- [ ] Can explain mutex vs semaphore vs monitor, precisely, not interchangeably
- [ ] Can code a producer-consumer with a bounded buffer in at least two languages (Go channels + Python `threading`)
- [ ] Can explain atomic variables / compare-and-swap
- [ ] Can deliberately write a race condition, then fix it three different ways (mutex, atomic, redesign to avoid shared state)
- [ ] Can write a test that would actually *catch* a race condition (ties to `go test -race`, ThreadSanitizer)
- [ ] Has a real answer ready for: "how do you test something that only fails 1 in 100 runs?" (repeated/stress runs, seeded randomization to force interleavings, chaos/fault injection, sanitizer tooling)
- [ ] Can explain starvation as distinct from deadlock

---

## Quick self-audit before Sunday night

Go section by section and mark honestly:
- [ ] Part A — attempted all 6 confirmed prompts at least once
- [ ] Part B — attempted at least 8 of the 18 scenarios, full framework each time
- [ ] Part C — solved at least 10 of the 15 curated problems without notes
- [ ] Part D/E/F — can answer every sample question above out loud, no hesitation
- [ ] Part G — have one project-grounded example ready for each concept
- [ ] Part H — can code producer-consumer in both Go and Python from memory

Whatever's still unchecked Sunday evening is your Monday-morning 6:15–6:45 review list — but by then it should be a short list, not this whole document.
