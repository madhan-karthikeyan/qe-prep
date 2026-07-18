# Couchbase QE — Final 58-Hour Plan
*Jul 17, 9:00 PM → Jul 20, 7:30 AM · Whole-day on-campus process, VIT Vellore*

---

## 0. TL;DR — your 5 highest-leverage moves

You said it yourself: projects and system design are your edge; CGPA and LC breadth are "moderate." Don't fight that — lean into it.

1. **Make system design and project narrative untouchable.** This is where you actually out-perform the field. Everything else just needs to clear the bar, not win the round.
2. **Don't try to become fluent in Go/Java/C++ in 2 days.** You can't, and trying will burn time you need elsewhere. Get to "won't embarrass myself" level: syntax, one idiomatic test written by hand, one concurrency primitive explained.
3. **Practice test-case design out loud, not just DSA.** This is a QE interview. An interviewer who sees you structure an answer as functional → boundary → negative → concurrency → security, unprompted, is scoring you higher than someone who just solves the DSA question silently.
4. **Know every number on your resume cold**, with the "why," not just the "what." This is confirmed as a real HR-round pattern for this company specifically.
5. **Protect your sleep on Sunday night.** A foggy brain on a 10+ hour interview day loses you more points than one extra hour of Go syntax ever gains you.

---

## 1. What the interviewer is actually scoring (the honest answer)

If I'm the one deciding who's in the 12–16, here's my actual scorecard, roughly in order of weight:

| Dimension | What "pass" looks like | What "fail" looks like |
|---|---|---|
| **Depth in one language** | You can go 4–5 layers deep under follow-ups in Python without hesitating | You know the syntax of five languages but can't defend a design choice in any of them |
| **Edge-case/test-design instinct** | You generate boundary, negative, concurrency, and security cases *before* being asked | You wait for the interviewer to prompt every category |
| **Thinking out loud** | Silence is filled with reasoning, even when stuck | Long silences, or jumping straight to code with no narration |
| **Resume defensibility** | Every number has a "why we chose this" story attached | You say "I don't remember the exact number" |
| **Grace under an unfamiliar task** | You state assumptions, ask clarifying questions, degrade gracefully | You freeze, or silently guess and hope |
| **Tool fluency** | Docker/Postman/Cypress/IDE open and used without fumbling | Spending 5 minutes just getting the tool to work |
| **Culture fit signals** | Concrete "why QE not dev," a real teamwork story, a genuine "why Couchbase" | Generic, rehearsed-sounding answers |

The unfamiliar-language rounds specifically exist to test #5, not your Go mastery. Nobody expects fresher-level Go fluency. They're watching *how* you handle not knowing something — which is arguably the single most QE-relevant trait there is.

---

## 2. What the day itself likely looks like

Grounded in recent Couchbase Vellore/Chennai campus-drive reports (treat as calibration, not a leaked script):

- **Structure:** resume/GitHub screen (already done — that's you) → 3–4 technical rounds → 1 HR round, each ~30–45 min, run back-to-back through the day with waiting gaps in between candidates.
- **Round flavor, roughly:** one round blends a project/role discussion + one DSA question + a paper/whiteboard LLD problem. Another blends resume review + core CS subjects (OS/CN/DBMS) + more DSA + a system-design implementation question. A later round leans system design + CS fundamentals + HR.
- **QE-specific pattern reported:** one DSA question (language of your choice, usually), then networking basics (HTTP vs HTTPS, status code categories, SQL injection/XSS), then a scenario question — "design test cases for [a booking site / login flow / API]."
- **Recurring Couchbase design prompts:** key-value store design, hash table implementation, parking lot OOD, tail command under low-RAM constraints, parallel sort on a multi-core machine with limited memory, malloc implementation.
- **Concurrency is a real theme, not JD filler:** mutex/semaphore code-ups, producer-consumer, race conditions, atomic variables have come up directly.
- **HR round grills the resume specifically** — keep every claim truthful and know the numbers cold.
- **The email said "AI tools" should be installed too** — this suggests you may be allowed/expected to use an AI coding assistant (Copilot, Claude Code, whatever you have) during a live task. If so, they still want to see *your* reasoning, not silent acceptance of suggestions — practice narrating while using one.

---

## 3. The 58-Hour Timetable

Adapt block boundaries to your energy — the point is coverage and rest discipline, not clock-watching.

### Fri Jul 17 — tonight (light setup only, brain is already tired)
- [ ] 9:00–11:30 PM: Install/verify every tool — Docker Desktop, Node + npm, Go, JDK, GCC/G++, Python venv, VS Code extensions, Postman, Git
- [ ] Scratch repos created for: Cypress spec, Go table-driven test, JUnit test, gtest example
- [ ] Print resume (2 copies) or confirm offline access; re-read it once, out loud
- [ ] **Sleep by 12:30 AM.** No new learning tonight — orientation only.

### Sat Jul 18 — full study day
| Time | Focus | Deliverable |
|---|---|---|
| 8:00–9:00 | Wake, breakfast | — |
| 9:00–11:00 | DSA sharpening (timed) | 2–3 medium problems, arrays/strings/hashing, narrate out loud while solving |
| 11:00–11:15 | Break | — |
| 11:15–1:15 | Go from zero | Syntax, error handling, structs/interfaces, one hand-written table-driven test with `go test -race` |
| 1:15–2:15 | Lunch + walk | — |
| 2:15–4:15 | Cypress/TypeScript | Install, write 2–3 specs (visit page, assert element, fill form), articulate Cypress vs Selenium |
| 4:15–4:30 | Break | — |
| 4:30–6:30 | Concurrency + distributed systems theory | CAP, replication, quorum, idempotency; code a producer-consumer in both Go and Python |
| 6:30–7:15 | Dinner | — |
| 7:15–9:15 | Java refresh | OOP, collections internals (`List`/`Map`/`Set`), one hand-written JUnit test |
| 9:15–9:30 | Break | — |
| 9:30–11:00 | C/C++ refresh | Pointers, manual memory management, one gtest example, explain what ThreadSanitizer catches |
| 11:00–11:30 | Wind down | — |
| **11:30 PM** | **Sleep (target 7h)** | — |

### Sun Jul 19 — full study day
| Time | Focus | Deliverable |
|---|---|---|
| 7:00–8:00 | Wake, breakfast | — |
| 8:00–10:00 | System design, pen-and-paper | Key-value store, parking lot OOD, hash table from scratch |
| 10:00–10:15 | Break | — |
| 10:15–12:15 | Memory-constrained design drills | Implement `malloc`, design `tail` under low RAM, parallel sort on multi-core with limited memory — these are *actual reported Couchbase questions* |
| 12:15–1:15 | Lunch | — |
| 1:15–3:15 | Test-case-design rehearsal (out loud, recorded) | Run all 5 scenarios in §5.3 below, structured functional/boundary/negative/security/concurrency |
| 3:15–3:30 | Break | — |
| 3:30–5:00 | Networking & security basics | HTTP vs HTTPS, status code categories, how you'd test for SQLi/XSS |
| 5:00–6:30 | Resume deep-dive rehearsal | Recite every number cold (§5.4); get someone to grill you like HR would |
| 6:30–7:15 | Dinner | — |
| 7:15–8:30 | Behavioral prep | Why QE not dev, why Couchbase, teamwork/conflict story, the AI/vector-search bridge story (§5.5) |
| 8:30–9:15 | Full mock run | Pick 2 categories, time-box them back-to-back to simulate whole-day fatigue |
| 9:15–10:00 | Pack + final tool check | `docker run hello-world`, `cypress open`, Postman logged in, both resume copies in bag |
| 10:00–10:30 | Unplug, relax | No screens — let it settle |
| **10:30 PM** | **Sleep (target 6.5–7h)** | Non-negotiable — see §7 |

### Mon Jul 20 — interview day
- [ ] 5:30 AM wake, shower
- [ ] 6:00 breakfast — something that won't spike/crash you (protein + slow carbs, not just sugar)
- [ ] 6:15–6:45 final review of your one-page numbers cheat-sheet only — no new material
- [ ] 6:45 leave, arrive with buffer
- [ ] 7:30 process begins

---

## 4. Language execution targets (cross-referenced to your master guide's tiers)

Your existing prep guide already has the full toolchain per language — this is just the "what must be *done*, not just read" version, mapped to when.

| Language | Minimum bar by Sun night | Stretch if time allows |
|---|---|---|
| **Python** | Can defend every fixture/parametrize choice in DecisionDrift's 363 tests without notes | Mention `Hypothesis` unprompted when discussing edge cases |
| **Go** | Wrote one table-driven test by hand, ran `go test -race`, can explain goroutines/channels via producer-consumer | Comfortable reading unfamiliar Go and guessing at `testify` syntax |
| **TypeScript/Cypress** | 2–3 working specs, can explain Cypress vs Selenium in one sentence | Know Playwright exists as the alt |
| **Java** | One hand-written JUnit test, can explain `List`/`Map`/`Set` internals | Know what `ExecutorService`/`synchronized` solve |
| **C/C++** | Can explain pointers + manual memory management out loud, one gtest example | Know what Valgrind and ThreadSanitizer each catch |
| **SQL → N1QL** | Comfortable with standard SQL joins/aggregations | 20-min Couchbase Capella free-tier spin-up, run one N1QL query |

---

## 5. Rehearsal bank

### 5.1 System design / LLD — pen and paper, no IDE
Practice these cold, with a diagram, in under 15 minutes each:
- [ ] Key-value store — how you'd optimally store and look up keys
- [ ] Parking lot — classic OOD, classes + relationships
- [ ] Hash table implementation from scratch (collision handling, resize)
- [ ] `malloc`/`free` implementation — what a simple allocator does
- [ ] `tail -f` on a file much larger than RAM — how you'd do it without loading the whole file
- [ ] Parallel sort across a multi-core machine, memory-constrained

### 5.2 Concurrency code-ups
- [ ] Producer-consumer with a bounded buffer (in Go using channels, and in Python using `threading`)
- [ ] Mutex vs semaphore vs atomic variable — code a small race condition, then fix it three ways
- [ ] Write a test that would *catch* a race condition (ties to `go test -race` / ThreadSanitizer)
- [ ] Have an opinion ready on "how do you test something that fails 1 in 100 runs?"

### 5.3 Test-case-design scenarios — rehearse out loud, structured
For each, walk through: functional → boundary/negative → integration → API → performance/load → security → **concurrency** (Couchbase's signature angle — "what if two users do this at the exact same time?"):
1. Login flow (username/password + MFA)
2. Train/flight ticket booking (the confirmed real example — two users booking the same seat)
3. A generic REST API endpoint (e.g., `POST /orders`)
4. File upload feature
5. Payment/checkout flow

### 5.4 Resume numbers — know these cold, with the "why"
- [ ] NIC: 150+ embedded docs, resolution time cut from >24h to <1 min, 4GB corpus, sub-second search — *why cosine similarity, why no dedicated search infra*
- [ ] GloballyGI: 5,000 images / 50 classes, YOLOv8n + D-FINE ensemble, 80–85% mAP, ~25% latency cut, ~35% mAP variance reduction — *why YOLOv8n over larger variants*
- [ ] DecisionDrift: 5 rule categories, AST + Tree-sitter across 12 languages, 363 tests, 21-patch benchmark, 95.2% Recall@5, 92.3% classification precision — *why deterministic, no LLM in the enforcement path — what trade-off that buys*
- [ ] Apex Sprint Planner: 6-stage Celery/Redis pipeline, 6 containerized services — *why retryable, why decouple from HTTP layer*
- [ ] JobLens: 7 ingestion sources, weights 35/30/20/10%, 90+ skills — *what each weight optimizes for*

### 5.5 Behavioral / narrative
- [ ] Why QE, not core dev (have a real answer, not a rehearsed-sounding one)
- [ ] Why Couchbase specifically (the AI Data Plane / vector search angle — your NIC internship is a genuinely rare fresher story here: a retrieval pipeline over embedded documents with metadata-filtered search is conceptually close to what their vector search product does at scale)
- [ ] One teamwork/conflict story, structured (situation → tension → what you did → outcome)
- [ ] "Do it right vs. move fast" — have a real opinion, not a platitude

---

## 6. Whole-day task-round protocol

If handed a live task in an unfamiliar language:

1. **Clarify first** — restate the problem, ask about constraints/edge cases before touching the keyboard
2. **Write test cases before code** — this is your QE-mindset differentiator; say it out loud: "before I implement, here's what I'd want this to handle"
3. **Pseudocode** — narrate structure before syntax
4. **Implement the minimal working version** — don't gold-plate
5. **Walk through your own test cases against it**
6. **State what you'd add given more time** (error handling, concurrency safety, etc.)
7. If genuinely stuck on syntax: it's fine to ask "can I write the logic in Python and translate the structure?" — one past DSA round explicitly allowed language of choice. Don't assume it's always allowed, but it's a reasonable ask if you're floundering.
8. If using an AI tool during the task: narrate why you're accepting/rejecting each suggestion — the interviewer is watching your judgment, not the tool's output.

**Stamina management for a 10+ hour day:** eat something small between every round if there's a gap, keep water on hand, use bathroom breaks to reset mentally rather than replaying the last round in your head, and treat each round as a fresh start — a rough round doesn't predict the next one.

---

## 7. Sleep discipline (non-negotiable)

You said not to treat time as a limiting factor — this is the one place I'll push back. A 10+ hour interview day on 4 hours of sleep loses you more composure and recall than any extra hour of Go syntax gains you. Hold the line on:
- Fri night: sleep by 12:30 AM
- Sat night: sleep by 11:30 PM (target 7h)
- **Sun night: sleep by 10:30 PM — this one matters most.** Everything in §3's Sunday evening block is designed to end early enough to make this possible.

---

## 8. Day-of carry list

- [ ] Laptop + charger (charged to 100%, pack the charger anyway)
- [ ] Power bank if you have one
- [ ] Docker Desktop confirmed working (`docker run hello-world`)
- [ ] Cypress installed and `cypress open` confirmed
- [ ] Postman installed and logged in
- [ ] Node.js and Python environments both verified (`node -v`, `python --version`)
- [ ] Git configured, GitHub accessible, repos pulled and buildable offline if possible
- [ ] Resume — 2 printed copies + accessible on laptop
- [ ] College ID / any admit card or permission letter mentioned in the CDC mail
- [ ] Pen + small notepad (for whiteboard/paper LLD rounds)
- [ ] Water bottle + snacks (explicitly flagged in the CDC email — no extended breaks)
- [ ] Comfortable layered clothing (interview rooms run cold with AC)
- [ ] Earbuds/charger cable spares if you use them for anything

---

*Built from your existing master prep guide + recent Couchbase Vellore/Chennai campus-drive reports. Treat the specific past questions as calibration, not a script for your exact batch.*
