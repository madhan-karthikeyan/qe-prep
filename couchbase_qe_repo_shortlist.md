# qe-prep Repo — What Actually Matters (Shortlist)
*683 files, ~840 directories. You don't need most of it — you already learned it by building it. This is the ~40-file subset worth actually opening before the interview. Total time: roughly 3.5 hrs.*

---

## Tier 0 — Pre-flight (not content — just verify it runs, ~15 min)
- [ ] `docker --version` && `docker run hello-world`
- [ ] `cd docker && docker-compose up` — dry-run your own `docker-compose.yml` once so it's proven working on your machine, not just present
- [ ] `git --version`, confirm push/pull auth actually works
- [ ] IDE extensions load without errors
- [ ] Any AI tool access mentioned in the invite is logged in and usable

---

## Tier 1 — Conceptual docs (`docs/`) — cheap, high yield (~35 min)
- [ ] `docs/qe_mindset.md` — sets the vocabulary for "how do you think about quality"
- [ ] `docs/unit-testing.md` — FIRST principles
- [ ] `docs/integration-testing.md` — testcontainers, real dependencies
- [ ] `docs/concurrency-testing.md` — race detection, stress testing (direct tie to your H1 checklist)
- [ ] `docs/debugging-techniques.md`
- [ ] `docs/docker-basics.md` — Docker fluency is explicitly called out as required day-of

Skip unless time remains: `boundary-value-analysis.md`, `mocking-stubbing.md` (already covered in your Python pytest section), `profiling.md`, `distributed-systems-testing.md`, `e2e-testing.md`, `bug-reports.md`.

---

## Tier 2 — Interview guides (`interview/`) — direct prep (~30 min)
- [ ] `interview/resume.md` — how to narrate your resume_projects
- [ ] `interview/testing_round.md`
- [ ] `interview/coding_round.md`
- [ ] `interview/behavioral.md`

Skim only if time remains: `system_round.md`, `networking.md`, `sql.md`, `distributed_systems.md`.

---

## Tier 3 — Resume projects — non-negotiable, know cold (~20 min)
- [ ] `resume_projects/decision_drift/README.md`
- [ ] `resume_projects/nic/README.md`
- [ ] `resume_projects/sprint_planner/README.md`

---

## Tier 4 — Flagship concurrency code (~45 min)
These ARE the "do this cold" exercises from your language checklist — re-open to refresh syntax, don't re-read the tests:
- [ ] `go/producer_consumer/producer_consumer.go`
- [ ] `go/rate_limiter/token_bucket.go`
- [ ] `go/thread_pool/pool.go`
- [ ] `python/producer_consumer/implementation/producer_consumer.py`
- [ ] `python/thread_pool/implementation/thread_pool.py`
- [ ] `java/producer_consumer/.../BlockingQueue.java`
- [ ] `c/producer_consumer/src/producer_consumer.c` — your one C concurrency touchpoint

---

## Tier 5 — Debugging scenarios (~45 min) — the actual QE skill signal
Read the README + hints, not a line-by-line broken/solution diff:
- [ ] `debugging/go/data_race/README.md` (+ hints) — rehearses your H1 checklist live
- [ ] `debugging/python/race_condition/README.md` (+ hints)
- [ ] `debugging/java/oom/README.md` (+ hints)
- [ ] `debugging/database/deadlock/README.md` (+ hints) — DB-specific, high relevance for Couchbase
- [ ] `debugging/docker/container_exit/README.md` (+ hints) — ties to the explicit Docker requirement
- [ ] `bug_reports/*.md` — all 5, they're short; this is literally "how to write a bug report," a core QE artifact

---

## Tier 6 — Patterns worth a glance (~25 min)
Couchbase is a database — these map directly. README + `implementation.py` only, skip the tests:
- [ ] `patterns/object_pool` — connection pooling, directly relevant to a DB client
- [ ] `patterns/retry`
- [ ] `patterns/exponential_backoff`
- [ ] `patterns/circuit_breaker`
- [ ] `patterns/pub_sub`

---

## Explicitly skip
- Every `target/`, `checkstyle-*`, `surefire-reports`, `generated-sources` path — Maven build byproducts, zero content value
- Every full test suite you already wrote (unit/integration/stress/fuzz) — you wrote them; re-reading isn't re-learning. If you can pass the self-tests in your language checklist from memory, the tests already did their job
- `distributed-systems/` and `fault_injection/` subfolders — no files show under these in your tree; if they're genuinely empty scaffolding, skip entirely
- `benchmarks/` — bonus only if every other tier is done

---

## Bottom line
This repo's real value right now isn't as a textbook — it's as proof you already did the reps, plus a couple dozen touchpoints to refresh. Don't try to re-earn what building it already earned you.
