# Couchbase QE — Multi-Language Prep Checklist
*Companion to `couchbase_qe_master_checklist.md`. That doc covers DSA/LLD/OOP/concurrency content once — this doc covers WHICH LANGUAGE to use HOW DEEP for each. Don't relearn concepts here; just calibrate depth per language.*

---

## 0. The rule

Python is your primary problem-solving language and default for coding rounds — general syntax fluency is already there, so it doesn't need a "learn this" section. What it does need is the **QE-specific layer**: pytest, mocking, and testing idioms — because "how would you automate this" is likely regardless of which language a round is nominally in.

Go/C/C++/Java exist because Couchbase's actual stack spans them (C/C++ = core KV engine, Go = query/service layer, Java = SDKs and OOP rounds), not because you're expected to match your Python fluency in each.

**Depth budget, total ~9.5–10.5 hrs:**

| Language | Priority | Target depth | Budget |
|---|---|---|---|
| Python (testing layer) | Primary | pytest + mocking fluency — general syntax already strong | ~2–2.5 hrs |
| Go | 2nd | Working fluency — can write real code cold | ~3 hrs |
| C | Conceptual | Can explain + write small, focused snippets | ~1.5 hrs |
| C++ | Conceptual | Mostly "C + what OOP/RAII adds" | ~1 hr |
| Java | 3rd | OOP fluency — one pattern coded cold | ~1.5–2 hrs |

---

## Part 1 — Python (testing layer)

**Target level:** not "learn Python" — you already know it. This is the QE-idiom layer: pytest fluency, mocking, and being able to translate any Part B scenario from your master checklist into an actual automated test, cold.

### 1.1 pytest core
- [ ] Test discovery/naming (`test_*.py`, `Test*` classes, `test_*` functions)
- [ ] Plain `assert` — pytest rewrites it, no need for `self.assertEqual` style
- [ ] Fixtures: `@pytest.fixture`, scope (`function`/`class`/`module`/`session`)
- [ ] `yield` fixtures for setup/teardown
- [ ] `conftest.py` — shared fixtures across files
- [ ] `@pytest.mark.parametrize` — data-driven tests
- [ ] `@pytest.mark.skip` / `xfail`
- [ ] `pytest.raises()` for exception testing
- [ ] Custom markers + marker-based selection (`-m` flag)

### 1.2 Mocking / test doubles
- [ ] `unittest.mock`: `Mock`, `MagicMock`, `patch`, `patch.object`
- [ ] `side_effect` vs `return_value`
- [ ] Mock vs stub vs spy vs fake — the actual distinction, and when each fits
- [ ] `pytest-mock` (the `mocker` fixture) as cleaner syntax over raw `unittest.mock`

### 1.3 Broader QE-relevant Python
- [ ] `threading` module + how to actually write a test that catches a race condition — ties directly to your H1 checklist; this is probably your strongest concurrency-testing story since it's your home language
- [ ] `requests` for API testing + `responses`/`httpretty` for mocking HTTP calls
- [ ] Context managers (`with`, `__enter__`/`__exit__`) — relevant for setup/teardown and resource cleanup
- [ ] Decorators — understand how `@pytest.fixture` and `@pytest.mark.parametrize` work under the hood; be able to write a simple one from scratch if asked
- [ ] Type hints (`typing` module) — a "do you write production-quality code" signal interviewers increasingly check for
- [ ] `pytest-cov` for coverage reporting
- [ ] `pytest-xdist` for parallel test execution (relevant if load/performance testing comes up)

### 1.4 Do these
- [ ] Write a parametrized pytest suite for one Part B scenario (e.g. Q1, the login flow) with at least 3 parametrized cases
- [ ] Mock a "slow or down" dependency with `patch` and assert your code degrades gracefully — this is a direct, working demo of B1 category 4 (Integration)
- [ ] Write one test that deliberately triggers *and catches* a race condition using `threading` — proves your H1 checklist item is real, not theoretical

### 1.5 Skip this
Deep pytest internals/plugin authoring, tox environment matrices, indirect parametrization, Django/Flask test clients (unless the role specifically needs them).

**Self-test:** given any one Part B scenario, write the pytest skeleton — fixtures + parametrize + at least one mock — cold, in under 15 minutes.

---

## Part 2 — Go

**Target level:** comfortable enough to write a correct, idiomatic producer-consumer or a re-implementation of a DSA problem you already solved in Python, from memory, in under 10 minutes.

### 1.1 Syntax checklist
- [ ] Variables, structs, slices vs arrays, maps
- [ ] Functions — multiple return values, named returns
- [ ] Pointers (`&` / `*`) — note: no pointer arithmetic, unlike C
- [ ] Error handling idiom (`if err != nil`) — know *why* Go does this instead of exceptions
- [ ] `defer` / `panic` / `recover`

### 1.2 Concurrency checklist (this is the part that actually matters most)
- [ ] Goroutines (`go func(){}()`)
- [ ] Channels — unbuffered vs buffered, send/receive, `close`
- [ ] `select` statement
- [ ] `sync.Mutex`, `sync.WaitGroup`
- [ ] `go test -race` — know what it catches and why (ties to your H1 checklist)

### 1.3 Do these
- [ ] Write a bounded producer-consumer from memory, cold
- [ ] Re-implement 2 problems from your DSA list (LRU cache is a good one) in Go
- [ ] State 2–3 concrete differences from Python: static typing, compiled not interpreted, goroutines+channels vs GIL-limited threads/asyncio

### 1.4 Skip this
Generics, module tooling, full stdlib, advanced channel patterns (fan-in/fan-out beyond the basics).

**Self-test:** bounded producer-consumer in Go, cold, <10 min.

---

## Part 3 — C

**Target level:** precise explanation of memory mechanics + can write small, focused programs. Not expected to be fluent at Python speed.

### 2.1 Checklist
- [ ] Pointers and pointer arithmetic
- [ ] `malloc` / `free` / `calloc` / `realloc` — what each actually does
- [ ] Stack vs heap allocation
- [ ] Structs
- [ ] Arrays vs pointers (decay)
- [ ] Manual string handling (no built-in string type)
- [ ] Common bug classes: dangling pointer, leak, buffer overflow, double free
- [ ] Can sketch how a simple `malloc` might work internally (free list, block headers) — this is your own A2 question #4, so it's near-confirmed territory

### 2.2 Do these
- [ ] Write a singly linked list (insert/delete) from memory
- [ ] Write a basic hash table (array of buckets + chaining) from memory

### 2.3 Skip this
Makefiles, preprocessor macro depth, multi-file compilation, anything beyond `stdio.h`/`stdlib.h`/`string.h`.

**Self-test:** malloc-based linked list insert/delete, cold.

---

## Part 3 — C++

**Target level:** "C plus what OOP/RAII adds." If C is solid, this is a short add-on, not a separate track.

### 3.1 Checklist
- [ ] Classes vs C structs — what OOP adds mechanically
- [ ] Constructors/destructors — why destructors matter for RAII
- [ ] `new`/`delete` vs `malloc`/`free`
- [ ] References vs pointers
- [ ] STL basics: know `vector` and `unordered_map` exist and roughly how they behave (don't need internals)
- [ ] Smart pointers (`unique_ptr`/`shared_ptr`) — conceptual only: what problem they solve

### 3.2 Skip this
Templates, operator overloading, multiple inheritance, move semantics, STL algorithms.

**Self-test:** explain RAII and why destructors help prevent resource leaks in a DB server context, in under 60 seconds, out loud.

---

## Part 4 — Java

**Target level:** OOP-fluent enough to code one design pattern cold and discuss SOLID naturally using it as the vehicle.

### 4.1 Checklist
- [ ] Class/interface syntax — `extends` vs `implements`
- [ ] Abstract class vs interface — when each fits
- [ ] Access modifiers
- [ ] Generics basics (`List<T>`)
- [ ] `try`/`catch`/`finally`, checked vs unchecked exceptions
- [ ] Common collections: `ArrayList`, `HashMap`
- [ ] `synchronized` keyword — how it maps to the mutex concept you already know
- [ ] Thread-safe Singleton (double-checked locking or eager init) — this is your own G2 Q4, near-confirmed
- [ ] One more pattern from memory: Factory or Observer

### 4.2 Skip this
Streams API, lambda syntax depth, JVM internals, any framework (Spring, etc.)

**Self-test:** thread-safe Singleton in Java, cold, <5 min — plus explain out loud why the naive version breaks under concurrency.

---

## Part 5 — Cross-language talking points (know these regardless of which language comes up)

Have a crisp 1–2 sentence answer ready for each, since interviewers may ask "how would this differ in X" even if you don't write the code:

- [ ] **Memory management:** garbage collected (Python/Java/Go) vs manual (C/C++)
- [ ] **Typing:** dynamic (Python) vs static (Go/Java/C/C++)
- [ ] **Concurrency model:** GIL-limited threads/asyncio (Python) vs goroutines+channels (Go) vs threads+locks (Java/C++) vs pthreads (C)
- [ ] **Execution:** interpreted (Python) vs compiled (Go/C/C++/Java-to-bytecode)

---

## Part 6 — Fallback protocol (use this live, in the interview)

If your syntax in a weaker language isn't clean under pressure: **narrate the logic instead of going silent.** This is already your own C3 rule ("narrate reasoning even when stuck"), it just applies double here.

> *"In Go I'd reach for a channel here instead of a mutex — let me sketch the shape, flag me if the syntax is off."*

This signals adaptability, which is the actual thing being evaluated — not perfect recall of a language you didn't grow up in.

---

## Self-audit

- [ ] Go: wrote a producer-consumer cold
- [ ] Go: re-implemented at least 1 Python DSA solve
- [ ] C: wrote a linked list AND a hash table cold
- [ ] C++: can explain RAII in under 60 seconds
- [ ] Java: wrote a thread-safe Singleton cold
- [ ] Cross-language: can answer all 4 talking points above without hesitation

Whatever's unchecked is where your remaining minutes go — not evenly across all four languages.
