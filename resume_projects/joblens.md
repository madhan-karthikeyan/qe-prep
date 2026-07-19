# JobLens — Complete Interview Preparation Guide

> For Madhan Karthikeyan · VIT Vellore B.Tech CSE 2026  
> Project: Personal Placement Operating System for Indian Freshers

---

# Phase 1 — Architectural Overview

## Project Purpose

JobLens is a **single-user, local-first "Personal Placement Operating System"** that automatically discovers, ranks, and tracks entry-level tech jobs for an Indian B.Tech final-year student. It replaces manual spreadsheet-based job tracking with automated scraping across 5 ATS platforms + 3 API aggregators, a deterministic scoring engine, a Kanban React UI, and a daily Discord digest.

## Problem Being Solved

Indian freshers face:
- Roles scattered across LinkedIn, Wellfound, company career portals, and aggregator APIs
- Irrelevant results (senior roles, wrong location, recruiter spam)
- No unified system to rank opportunities against personal skills
- Tedious manual application tracking

**Key insight**: The candidate has two career tracks (ML and Software Engineering) with different skill sets — the system needs to rank the same job pool differently for each profile.

## Architecture (High-Level)

```
┌─────────────────────────────────────────────────────────────────┐
│                    OS Cron (8:00 AM)                            │
│  python -m scripts.run_daily_pipeline                          │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│  COLLECTOR LAYER                                                │
│  ┌─────────────────┐  ┌──────────────────────────────────────┐  │
│  │ API Scrapers    │  │ ATS Portal Providers                 │  │
│  │ ─ Adzuna        │  │ ─ Greenhouse (12+ companies)        │  │
│  │ ─ JSearch       │  │ ─ Lever (3 companies)               │  │
│  │ ─ YC Work/Startup│  │ ─ Ashby (3 companies)              │  │
│  └────────┬────────┘  │ ─ Workday (5 companies)             │  │
│           │           │ ─ Custom (7 companies, config-driven)│  │
│           │           └──────────────────┬───────────────────┘  │
│           └──────────────┬───────────────┘                      │
│                          ▼                                      │
│              Deduplicate → Merge → Normalize                    │
└──────────────────────────┬──────────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────────┐
│  CLASSIFICATION LAYER (ingestion-time)                          │
│  ┌────────────────┐ ┌──────────────┐ ┌────────────────┐        │
│  │ Role Family    │ │ Experience   │ │ Fresher        │        │
│  │ (9 families)   │ │ Level (7)    │ │ Eligibility    │        │
│  └────────────────┘ └──────────────┘ └────────────────┘        │
│  ┌────────────────┐ ┌──────────────────────────────────┐       │
│  │ Relevance      │ │ Years Required Extraction        │       │
│  │ Score (0-100)  │ │ (regex from title/description)   │       │
│  └────────────────┘ └──────────────────────────────────┘       │
└──────────────────────────┬──────────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────────┐
│  DATABASE LAYER                                                 │
│  PostgreSQL 16 + pgvector                                       │
│  10 tables (job_posts, saved, applied, feedback, etc.)         │
│  Alembic migrations (14 versions)                               │
└──────────────────────────┬──────────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────────┐
│  RANKING ENGINE                                                  │
│  ┌──────────────────────────────────────────────────────────────┐│
│  │  Score = role×35% + skill×30% + behavior×20% +              ││
│  │          freshness×10% + bonuses(±20 capped)                ││
│  │                                                            ││
│  │  Plus: garbage filter, role filter, years filter,          ││
│  │  skill overlap, seniority penalties, dedup, explainability ││
│  └──────────────────────────────────────────────────────────────┘
└──────────┬───────────────────────────────────────┬───────────────┘
           │                                       │
┌──────────▼──────────┐           ┌─────────────────▼───────────┐
│  FastAPI REST API   │           │  NOTIFICATION LAYER          │
│  15 route modules   │           │  Discord Webhook Digest      │
│  16 service modules │           │  ─ new opportunities         │
│  Pydantic v2 schemas│           │  ─ pipeline summary          │
│  CORS: localhost:8080│          │  ─ deadlines                 │
└──────────┬──────────┘           │  ─ top recommendations       │
           │                      └─────────────────────────────┘
┌──────────▼──────────┐
│  React 19 Frontend  │
│  TanStack Start     │
│  TanStack Router    │
│  TanStack Query     │
│  shadcn/ui          │
│  Tailwind CSS v4    │
│  dnd-kit Kanban     │
│  recharts charts    │
└─────────────────────┘
```

## Key Architectural Decisions

| Decision | Rationale | Tradeoff |
|---|---|---|
| **Deterministic ranking instead of ML** | Zero cost, deterministic, fully explainable. Every score reason is traceable to specific keywords. | Cannot capture semantic similarity. "PyTorch" spelled "Pytorch" won't match. Misses synonyms. |
| **Async throughout** | All I/O is async (httpx, asyncpg). Suitable for 30+ concurrent HTTP requests to ATS portals. | Adds complexity for the scrape pipeline. Error handling is more complex. |
| **Ingestion-time classification** | Classifications (role_family, experience_level, freshness) computed once at scrape time, stored in DB. Ranking never reclassifies. | If classification rules change, old jobs must be reclassified (scripts exist for this). |
| **Two profiles (ML + Software)** | The same candidate targets two distinct career tracks. Ranking must be profile-aware. | Duplicates work — two sets of tests, two sets of skill weights. |
| **Local-first, single-user** | No auth, no multi-tenancy, no cloud deployment. Simplifies everything. | Cannot share with peers. No backups. No access from outside local network. |
| **Discord webhook (no bot)** | Zero infra — just an HTTP POST to a webhook URL. No bot token management. | No interactive features. One-way notification only. |

## Technology Stack Summary

| Layer | Technology | Version |
|---|---|---|
| Language | Python | 3.12+ |
| API Framework | FastAPI | 0.116+ |
| ASGI Server | Uvicorn | 0.35+ |
| ORM | SQLAlchemy | 2.0+ (async) |
| DB Driver | asyncpg | — |
| Database | PostgreSQL | 16 + pgvector |
| Validation | Pydantic v2 + pydantic-settings | — |
| HTTP Client | httpx | 0.28+ |
| HTML Parsing | BeautifulSoup4 + lxml | — |
| Frontend | React 19 + TanStack Start | — |
| UI | shadcn/ui + Tailwind CSS v4 | — |
| Testing | pytest + pytest-asyncio | — |
| Linting | Ruff | 0.12+ |
| Container | Docker + Docker Compose | v2+ |

## External Services

| Service | Free Tier | Purpose |
|---|---|---|
| Adzuna API | 250 req/day | Job aggregator for India |
| JSearch API (Open Web Ninja) | Paid | Job search API |
| YC Work at a Startup | Free | HTML scraping |
| Greenhouse | Free | ATS career portal API |
| Lever | Free | ATS career portal API |
| Ashby | Free | ATS career portal API |
| Workday | Free | ATS career portal API |
| Discord Webhook | Free | Daily digest notification |

---

# Phase 2 — Knowledge Dependency Graph

Concepts ordered from most foundational → most project-specific:

```
Level 1: Programming Fundamentals
├── Python 3.12+ (async/await, match statements, type hints, dataclasses)
├── TypeScript 5.8+
└── HTTP Protocol (methods, status codes, headers)

Level 2: Web Fundamentals
├── REST API Design
├── JSON Serialization
├── URL Structure & Query Parameters
├── CORS (Cross-Origin Resource Sharing)
└── Web Scraping Ethics & Rate Limiting

Level 3: Async Programming
├── asyncio (event loop, coroutines, tasks)
├── async/await patterns
├── Connection Pooling
├── Exponential Backoff & Retry
└── Timeout Handling

Level 4: Web Frameworks
├── FastAPI
│   ├── ASGI vs WSGI
│   ├── Path/Query Parameters
│   ├── Dependency Injection
│   ├── Request Validation
│   └── Middleware
├── React 19
│   ├── Components & Hooks
│   ├── TanStack Router
│   ├── TanStack Query (server state)
│   └── TanStack Start (SSR)
├── Uvicorn (ASGI Server)
└── Tailwind CSS v4

Level 5: Database
├── PostgreSQL 16
│   ├── Indexing (B-tree, composite)
│   ├── ON CONFLICT (upsert)
│   ├── Full-Text Search (optional)
│   └── LIMIT/OFFSET Pagination
├── pgvector Extension
│   └── Vector embeddings (384-d)
├── SQLAlchemy 2.0+
│   ├── Declarative Mapping
│   ├── Async Session
│   ├── Relationship Loading
│   └── Type Annotations (Mapped[])
├── Alembic
│   ├── Revision Files
│   ├── Upgrade/Downgrade
│   └── Autogeneration
├── asyncpg
│   ├── Parameter Limits (32767)
│   └── Prepared Statements
└── Pydantic v2
    ├── model_validate
    ├── model_dump
    ├── field_validator
    ├── computed_field
    └── ConfigDict

Level 6: Scraping & Data Collection
├── httpx (async HTTP)
│   ├── Timeout Configuration
│   ├── Redirect Following
│   ├── User-Agent Rotation
│   └── Response raise_for_status
├── BeautifulSoup4
│   ├── HTML Parsing
│   ├── get_text()
│   └── lxml Backend
├── API Key Authentication
├── HTML Regex Extraction
├── Sitemap-based URL Discovery
└── JSON API Interaction Patterns

Level 7: Classification & Ranking
├── Regex-Based Classification
│   ├── Pattern Compilation
│   ├── Word Boundary (\b)
│   └── Case-Insensitive Matching
├── Weighted Scoring Algorithms
│   ├── Component Weighting (%, flat)
│   ├── Bonus Pool Capping
│   └── Score Clamping (0-100)
├── Normalization Techniques
│   ├── Text Lowercasing
│   ├── Alias Expansion (ML → machine learning)
│   ├── Company Name Normalization
│   └── Word Substitutions
├── Deduplication Strategies
│   ├── Composite Key (company, title)
│   ├── Set-Based O(n) Dedup
│   └── Highest-Score Wins
└── Explainability
    ├── Score Breakdown
    ├── Match Reasons
    └── Confidence Levels

Level 8: Design Patterns
├── Abstract Base Class
├── Registry Pattern
├── Strategy Pattern
├── Template Method
├── Singleton (cached)
├── DTO (Data Transfer Object)
├── Fire-and-Forget
├── Factory Method
└── Builder Pattern

Level 9: DevOps & Infrastructure
├── Docker
│   ├── Multi-stage Builds
│   ├── Dockerfile Best Practices
│   └── Docker Compose Networking
├── Docker Compose
│   ├── Service Dependencies
│   ├── Volume Mounts
│   ├── Health Checks
│   └── Environment Variables
├── Cron Scheduling
├── uvloop (fast asyncio event loop)
└── Ruff Linting

Level 10: Project-Specific
├── ATS Portal Architecture
│   ├── Greenhouse Board API
│   ├── Lever.co API
│   ├── Ashby API
│   ├── Workday CXS API
│   └── Custom Portal Scraping
├── Fresher-Specific Job Classification
│   ├── India Job Detection (50 cities)
│   ├── Seniority Detection (Sr., Staff, Principal)
│   ├── Experience Level (II, III, IV, V)
│   └── Garbage Recruiter Filter
├── Behavior Learning (no ML)
│   ├── Affinity Scoring
│   ├── Interaction Weighting
│   └── Minimum Interaction Threshold
├── Discord Digest
│   ├── Markdown Formatting
│   └── Webhook HTTP POST
└── TanStack Start (React SSR)
```

---

# Phase 3 — Concept Curriculum

Each concept explained: what, why, problem solved, how it works internally, where it appears in JobLens, alternatives, tradeoffs, misconceptions, and interview follow-ups.

---

## 3.1 FastAPI

### What
FastAPI is a modern Python web framework for building REST APIs. It's built on Starlette (for ASGI) and Pydantic (for data validation).

### Why it exists
Django and Flask were designed before async Python was mainstream. FastAPI fills the gap for high-performance async APIs with automatic OpenAPI documentation.

### Problem it solves
- Request parsing and validation without boilerplate
- Automatic OpenAPI/Swagger docs
- Async request handling without thread pool overhead
- Type safety via Pydantic integration

### How it works internally
1. Uvicorn reads ASGI scope (connection info, headers, body) from the network socket
2. FastAPI's `ASGIApp` routes the request to the matching path operation
3. Before the handler runs, **request validation** happens:
   - Path parameters extracted via regex from the path template
   - Query parameters parsed from the query string
   - Body parsed as JSON and validated against the Pydantic model
   - All fields go through Pydantic's type coercion and validation pipeline
4. The validated parameters are injected into the handler function
5. The handler returns a Python dict/list/Pydantic model
6. FastAPI serializes the return value to JSON via `jsonable_encoder`
7. Response headers (including CORS) are added by middleware

### Where it appears
- `api/main.py`: App factory with CORS middleware, 15 routers
- `api/routes/*.py`: 15 route modules
- `api/schemas/*.py`: 11 Pydantic request/response schemas
- `api/services/*.py`: 16 service modules

### Why JobLens chose it
- Async native (all scraping is async httpx calls)
- Automatic validation of scrape configurations and job data
- OpenAPI docs are useful for debugging scrapers
- Lightweight — no unnecessary abstraction

### Alternatives
- **Django REST Framework**: Too heavy, sync-only (needs additional async layers), more boilerplate
- **Flask**: Sync-only, requires manual validation, no built-in OpenAPI
- **Starlette directly**: Too low-level, would need to reimplement validation

### Tradeoffs
- **+** Automatic validation eliminates an entire class of bugs
- **+** Async support is first-class
- **-** Pydantic v2 model compilation adds import-time overhead
- **-** Dependency injection is less flexible than manual wiring

### Common Misconceptions
- "FastAPI is faster than Flask because it's written in C" — No, it's Python. Speed comes from async I/O and Pydantic's Rust-based validation core (pydantic-core is Rust)
- "FastAPI requires async" — Sync routes work too (run in thread pool)
- "FastAPI generates production UIs with Swagger" — Swagger UI is for development/testing only

### Interview Follow-ups
1. How does FastAPI differ from Starlette?
2. What happens when you make a sync route in FastAPI?
3. How does dependency caching work (`Depends` with `use_cache=True`)?
4. How does `jsonable_encoder` handle edge cases like `datetime` or `Decimal`?
5. What's the `request` object and how do you access raw body?

---

## 3.2 Async/Await & asyncio

### What
`asyncio` is Python's library for writing concurrent code using the async/await syntax. It's a **cooperative multitasking** model — tasks voluntarily yield control at `await` points.

### Why it exists
Traditional threading in Python has the GIL problem. Threads are also expensive (~8KB stack per thread). `asyncio` gives concurrency without the GIL bottleneck or thread overhead.

### Problem it solves
- I/O-bound workloads (HTTP calls, database queries, file reads)
- Thousands of concurrent connections without thread overhead
- Clean sequential-looking code for concurrent operations

### How it works internally
1. The **event loop** (`asyncio.run()` creates one) maintains a task queue
2. A **task** wraps a coroutine with scheduling state
3. When a coroutine `await`s, it creates a **future** and suspends
4. The event loop runs other tasks while waiting
5. When the future resolves (I/O completes), the task is rescheduled
6. `uvloop` is a drop-in replacement for the event loop written in Cython (~2x faster)

### Where it appears
- Every scraper: `await client.get(url)`, `await fetch_with_retry()`
- `collector/runner.py`: `asyncio.wait_for()` with timeouts, `asyncio.gather()` for concurrent scraping
- `scripts/run_daily_pipeline.py`: Full async pipeline, `asyncio.ensure_future()` for fire-and-forget logging
- Every database call: `await session.execute()`, `await session.commit()`
- `collector/collector_registry.py`: `asyncio.gather(*tasks, return_exceptions=True)`

### Key Implementation Detail (JobLens)
```python
# Fire-and-forget pattern for scraper logging
asyncio.ensure_future(
    log_scraper_run(...)
)
# This creates a task but does NOT await it.
# The task will run "in the background" on the event loop.
# Potential issue: if the pipeline finishes before logging completes,
# the event loop may shut down and drop the task.
```
The pipeline does attempt to drain pending logs at the end:
```python
if _PENDING_LOGS:
    await asyncio.gather(*_PENDING_LOGS, return_exceptions=True)
```
But `_PENDING_LOGS` is never actually populated — `_schedule_for_drain()` is a no-op. This is a **bug**: fire-and-forget tasks created via `ensure_future` on line 57 and 108 are not tracked, so they may be dropped.

### Alternatives
- **Threading**: True parallelism for CPU work, but GIL-limited, higher overhead
- **Multiprocessing**: Bypasses GIL, but heavy, no shared state
- **Greenlets (gevent)**: Monkey-patching approach, less explicit

### Tradeoffs
- **+** Thousands of concurrent connections on one thread
- **+** No race conditions from shared memory (single-threaded)
- **-** CPU-bound tasks block the event loop (need `run_in_executor`)
- **-** Stack traces are harder to read (dozens of suspended frames)
- **-** Requires all I/O libraries to be async-aware

### Common Misconceptions
- "async/await makes Python faster" — No, it makes I/O more efficient. CPU-bound code is the same speed.
- "await releases the GIL" — There is no GIL release in asyncio (it's single-threaded). I/O releases the GIL.
- "asyncio.run() creates a new event loop" — It creates one if none exists, reuses an existing one in Python 3.12+.

### Interview Follow-ups
1. How does the event loop know when an I/O operation completes?
2. What happens to pending tasks when the event loop closes?
3. What's the difference between `asyncio.create_task()` and `ensure_future()`?
4. How do you debug a "coroutine was never awaited" warning?
5. When would you use `asyncio.gather()` vs `asyncio.wait()` vs `asyncio.TaskGroup()`?

---

## 3.3 PostgreSQL & asyncpg

### What
PostgreSQL is a relational database. asyncpg is a high-performance PostgreSQL driver for Python/asyncio.

### Why asyncpg exists
psycopg2 blocks on every query. asyncpg speaks PostgreSQL's wire protocol directly using asyncio, giving 3-5x throughput over psycopg2 in concurrent workloads.

### Problem it solves
- Database queries shouldn't block the event loop
- Connection pooling for async workloads
- Prepared statement caching for repeated queries

### How asyncpg works internally
1. Opens a TCP connection to PostgreSQL (port 5432 by default)
2. Sends messages in PostgreSQL's native **v3 wire protocol** (binary format)
3. Supports **prepared statements**: first send `PARSE`, then `BIND`, then `EXECUTE`
4. Uses **server-side cursors** for large result sets
5. Connection pool maintains N connections, each in its own asyncio task

### Where it appears
- `db/session.py`: `create_async_engine()` with `pool_size=10`
- Batch insert with 500-row chunks (avoids asyncpg's 32767 parameter limit)

### Key asyncpg detail
```python
# asyncpg has a per-statement parameter limit of 32767
# JobLens batches in groups of 500 to stay well under this:
_BATCH_SIZE = 500
for i in range(0, len(jobs), _BATCH_SIZE):
    batch = jobs[i : i + _BATCH_SIZE]
    # Each JobPost has ~16 columns → 500 × 16 = 8000 params
```

### SQLAlchemy + asyncpg integration
SQLAlchemy 2.0's async engine uses asyncpg as the driver. Key detail:
```python
# SQLAlchemy uses autocommit-like behavior with begin_once
# session.execute() wraps a begin/commit pair
# But explicit commit is needed for write operations
```

### Interview Follow-ups
1. What happens when the connection pool is exhausted?
2. How does `ON CONFLICT DO NOTHING` handle primary key conflicts vs unique constraint conflicts?
3. What is PostgreSQL's MVCC and how does it affect the scraper's batch insert?
4. Why does asyncpg have a parameter limit? What's the workaround?
5. How does SQLAlchemy translate ORM queries to asyncpg wire messages?

---

## 3.4 Pydantic v2

### What
Pydantic v2 is a Python library for data validation using Python type hints. The core validation engine is written in Rust (pydantic-core).

### Why it exists
Runtime type checking and validation. Python's type hints are static (checked by mypy/pyright). Pydantic makes them runtime-enforced.

### Problem it solves
- "Stringly-typed" APIs where every function must validate its inputs
- Serialization/deserialization between Python dicts and JSON
- Config management with .env file loading
- Ensuring scraped job data meets schema requirements

### How it works internally
1. A Pydantic model class defines fields with type annotations
2. At class definition time, `__init_subclass__` (via `BaseModel`) builds a **schema** of all fields
3. Pydantic-core (Rust) compiles a validator function tailored to this schema
4. When `model_validate()` or `__init__()` is called, pydantic-core processes the input:
   - Coerce types (e.g., `"42"` → `42` for `int` fields)
   - Run validators (e.g., `@field_validator`)
   - Throw `ValidationError` if any field fails
5. `model_dump()` serializes back to a Python dict with type coercion

### Where it appears
- `collector/models.py`: `JobPostData` — the canonical job data model
- `db/schema/job_post.py`: `JobPostSchema` — DB-specific, with `from_attributes=True` for ORM
- `api/schemas/*.py`: 11 request/response models
- `config.py`: `Settings` via pydantic-settings (with `.env` loading)
- `user/models.py`: `ProfileVariant` and `Project`
- `collector/company_config.py`: `CompanyConfig` with validation

### Key Pydantic features used
```python
# computed_field — derived value, not stored
class ProfileVariant(BaseModel):
    @computed_field
    @property
    def skills(self) -> list[str]:
        return self.technical_skills + self.frameworks + ...

# Frozen/immutable model
model_config = {"frozen": True}  # ProfileVariant can't be modified after creation

# from_attributes for ORM
class JobPostSchema(BaseModel):
    model_config = ConfigDict(from_attributes=True)
    # This allows: JobPostSchema.model_validate(db_orm_object)

# default_factory for dynamic defaults
scraped_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

# Field validation
class CompanyConfig(BaseModel):
    name: str = Field(..., min_length=1)
```

### Alternatives
- **attrs**: Validation is manual, less ergonomic for APIs
- **dataclasses**: No validation, no serialization
- **marshmallow**: Similar to Pydantic v1, slower (pure Python)
- **msgspec**: Faster for serialization, less ecosystem

### Tradeoffs
- **+** Rust-based core is very fast
- **+** Seamless FastAPI integration
- **-** Model imports trigger schema compilation (slower import time)
- **-** Error messages can be cryptic for nested models

### Interview Follow-ups
1. How does `model_validate` differ from `__init__`?
2. What's the difference between `ConfigDict(from_attributes=True)` vs `from_orm`?
3. How does `field_validator` work and how is it ordered?
4. What's pydantic-core and why was Pydantic rewritten in Rust?
5. How does Pydantic handle `Optional` vs `Union[None, T]`?

---

## 3.5 SQLAlchemy 2.0 (Async)

### What
SQLAlchemy 2.0 is the latest major version of Python's most popular ORM, with full async support via `asyncpg`.

### Why it exists
SQL is powerful but verbose. An ORM maps database rows to Python objects, reducing boilerplate and providing type safety.

### Problem it solves
- Writing raw SQL for every query is error-prone
- Result rows need to be mapped to Python objects
- Connection management is complex
- Schema changes require coordinated code changes

### How SQLAlchemy async works internally
1. `create_async_engine()` creates a pool of asyncpg connections
2. Each `AsyncSession` wraps a connection from the pool
3. Queries are built as a SQL expression tree (not strings)
4. When `await session.execute(stmt)` is called:
   - SQLAlchemy compiles the expression tree to a SQL string
   - The asyncpg connection sends the SQL over TCP
   - The result cursor is buffered into memory
   - Rows are mapped to ORM objects via `Mapper`
5. `await session.commit()` sends a COMMIT message

### Key SQLAlchemy features used
```python
# Type-annotated mappings (PEP 484)
class JobPost(Base):
    __tablename__ = "job_posts"
    id: Mapped[str] = mapped_column(String(16), primary_key=True)

# PostgreSQL-specific upsert
from sqlalchemy.dialects.postgresql import insert
stmt = insert(JobPost).values(values).on_conflict_do_nothing(
    index_elements=[JobPost.id]
)

# Server defaults
scraped_at: Mapped[datetime] = mapped_column(
    DateTime(timezone=True),
    server_default=text("TIMEZONE('utc', now())"),
)
```

### Where it appears
- `db/schema/`: 10 ORM models across 12 files
- `db/session.py`: Engine setup, session factory, batch insert, upsert
- `api/services/*.py`: All database queries for the API
- `backend/services/digest_builder.py`: Pipeline statistics queries

### Alternatives
- **SQLAlchemy 1.4 sync**: No async, blocks event loop
- **psycopg2 direct**: Raw SQL everywhere
- **asyncpg direct**: No ORM, manual mapping
- **Tortoise ORM**: Django-ish async ORM, less mature

### Interview Follow-ups
1. What's the N+1 query problem and how does SQLAlchemy solve it?
2. How does `selectinload` differ from `joinedload`?
3. Why does `insert().on_conflict_do_nothing()` return `rowcount` as `-1` sometimes?
4. What's the difference between `scalars()` and `execute()` in SQLAlchemy 2.0?
5. How do you handle `expire_on_commit=False` and its implications?

---

## 3.6 Regex-Based Classification

### What
The entire classification system (role family, experience level, fresher eligibility, seniority detection) is built on regular expressions. No ML, no NLP, no LLM.

### Why it exists
For structured job titles with predictable patterns, regex is **deterministic, zero-cost, and explainable**. It can't misclassify in surprising ways like ML can.

### Problem it solves
- Determining if a job posting is for freshers (entry-level)
- Categorizing by role family (ML, backend, frontend, etc.)
- Detecting seniority level from title keywords
- Filtering garbage recruiter spam

### Key regex patterns used

```python
# Experience level classification
_EXPERIENCED_LEVEL_PATTERNS: dict[str, int] = {
    r"\bii\b": -5,
    r"\biii\b": -10,
    r"\bengineer\s*[-–]?\s*[2-9]\b": -5,
    r"\bl[2-9]\b": -5,
}

# Seniority detection
_SENIORITY_PATTERNS: dict[str, int] = {
    r"\bsr\.?\b": -20,
    r"\bsenior\b": -20,
    r"\bprincipal\b": -40,
    r"\bstaff\b": -30,
}

# Year extraction from title
_YEARS_IN_TITLE_RE = re.compile(
    r"(?:^|\D)"
    r"(?:(\d+)\s*(?:-|–|to)\s*)?(\d+)\s*\+?\s*(?:yr|year)s?\b",
    re.IGNORECASE,
)

# Garbage recruiter filter
_KNOWN_GARBAGE_PATTERNS: list[str] = [
    "hiring for", "urgent hiring", "placement",
    "recruitment", "the elite job", "guidance placement",
]
```

### Where it appears
- `candidate_engine/ranker.py`: ~30 regex patterns across 6 functions
- `services/classification/fresher.py`: Fresher eligibility
- `services/classification/experience.py`: 7-level experience classification
- `services/classification/role_family.py`: 9 role families
- `services/scoring/relevance.py`: Keyword-based relevance scoring

### Why JobLens chose regex over ML
- **Cost**: Zero API calls, zero GPU time, zero model hosting
- **Determinism**: Same input always produces same output
- **Explainability**: "This job is classified as Senior because title contains 'Staff'"
- **Speed**: <1ms per job
- **No training data needed**

### Limitations
- Cannot handle misspellings ("PyTorch" vs "Pytorch")
- Cannot understand context ("Senior" in "Senior Software Engineer" is senior, but "Senior Year Internship" is not)
- Cannot handle synonyms not in the alias list
- Pattern order matters — first match wins in some cases
- Regex can't understand that "2+ years" might still be ok for a fresher with strong portfolio

### Interview Follow-ups
1. How would you handle "Python Developer" being classified as both backend and ML?
2. What happens when a job title is "Senior Principal Staff Architect" — multiple patterns match?
3. How do you benchmark classification accuracy?
4. When would regex classification fail and ML would be necessary?
5. How would you add a new role family?

---

## 3.7 Docker & Docker Compose

### What
Docker packages applications into containers (lightweight VMs sharing the host OS kernel). Docker Compose orchestrates multi-container applications.

### Why it exists
"Works on my machine" problem — containers ensure identical environments across dev, test, and production.

### Problem it solves
- PostgreSQL 16 with pgvector is non-trivial to install natively
- Playwright with Chromium requires system dependencies
- Python dependency isolation
- Consistent development environment

### JobLens Docker Architecture

```yaml
services:
  jobhunter-db:
    image: pgvector/pgvector:pg16
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5434:5432"    # Note: mapped to non-default port 5434

  pgadmin:
    image: dpage/pgadmin4
    # Admin UI for database management

  api:
    build: .
    ports:
      - "8000:8000"
    depends_on:
      - jobhunter-db
    command: uvicorn api.main:app --host 0.0.0.0 --reload

  frontend:
    build: frontend/
    ports:
      - "8080:8080"
    depends_on:
      - api
```

### Key details
- **API Dockerfile**: Python 3.12-slim + Playwright Chromium install (~300MB+)
- **Frontend Dockerfile**: Node 22-alpine, Vite dev server
- **Database port**: 5434 (not the default 5432), likely to avoid conflicts
- **No health checks** visible in the compose file
- **Hot reload** enabled for development (`--reload` flag)

### Interview Follow-ups
1. What's the difference between `CMD` and `ENTRYPOINT` in Dockerfiles?
2. Why is the database port 5434 instead of 5432?
3. How would you add health checks to the compose file?
4. What's the security concern with running Chromium in a container?
5. How would you reduce the Docker image size?

---

## 3.8 TanStack Start (React SSR)

### What
TanStack Start is a full-stack React framework built on TanStack Router. It provides SSR, bundling, and routing.

### Why it exists
React's ecosystem was fragmented — Next.js was the dominant meta-framework but had vendor lock-in. TanStack Start offers a more modular, type-safe alternative.

### Key features used
- File-based routing (or route modules in `routes/`)
- TanStack Router for type-safe navigation
- TanStack Query for server state management
- SSR for initial page load
- Vite for bundling (Vite 7.3+)

### Interview Follow-ups
1. How does SSR differ from CSR? When is each appropriate?
2. How does TanStack Router's type-safe routing work?
3. What happens when you make a request from the server vs client in TanStack Start?
4. How does `@tanstack/react-query` differ from `useEffect`-based data fetching?
5. What's the cache invalidation strategy for the Kanban board?

---

## 3.9 Design Patterns in JobLens

### Registry Pattern
```python
# collector/collector_registry.py
class CollectorRegistry:
    def __init__(self):
        self._collectors: dict[str, Collector] = {}

    def register(self, collector: Collector) -> None:
        self._collectors[collector.name] = collector

    async def run_all(self) -> list[JobPost]:
        names = list(self._collectors.keys())
        tasks = [self.run(name) for name in names]
        results = await asyncio.gather(*tasks, return_exceptions=True)
```
**Why**: New scrapers/portals can be added by registering them without modifying the orchestration code. **Usage**: Both API scrapers and ATS portal providers use the same registry interface.

### Template Method Pattern
```python
# Each provider follows the same lifecycle:
# 1. configure(config) — set company-specific parameters
# 2. fetch() — retrieve raw data from the API
# 3. _normalize() — convert to JobPostData
```
Greenhouse, Lever, Ashby, Workday, and Custom providers all implement this same template.

### Strategy Pattern
```python
# Different ranking strategies per profile type
def rank_jobs_for_profile(jobs, profile_type="ml"):
    profile = get_profile(profile_type)
    min_overlap = 2 if profile_type == "ml" else 1
    # Same pipeline structure, different parameters per profile
```

### DTO Pattern
```python
# Separate Pydantic models for layer isolation
# API layer: api/schemas/*.py → what the client sees
# DB layer: db/schema/job_post.py JobPostSchema → what gets stored
# Internal: collector/models.py JobPostData → what passes through scrapers
```

### Fire-and-Forget Pattern (with bug)
```python
# Pipeline doesn't wait for scraper run logs
asyncio.ensure_future(log_scraper_run(...))
# BUG: _PENDING_LOGS is never populated, so drain is a no-op
```

### Interview Follow-ups
1. Why is the Registry pattern used instead of a simple list of collectors?
2. How would you add a new ATS provider? Walk through the process.
3. What's the problem with the fire-and-forget logging pattern?
4. When should you prefer composition over inheritance (like the Template Method here)?
5. How does the DTO pattern prevent coupling between layers?

---

# Phase 4 — Repository Deep Dive

## 4.1 `candidate_engine/ranker.py` (862 lines) — The Heart of the System

### Purpose
Deterministic, explainable job ranking engine. Takes a list of `JobPostData` objects, a `profile_type`, and optional affinities, returns a sorted list of `JobMatchResult`.

### Pipeline (executed per job per profile)

```
_is_garbage_job() → skip if recruiter spam
_is_title_allowed() → skip if excluded role family
_years check → skip if >1 year experience required
_compute_skill_overlap() → weighted keyword matching
_score_role_alignment() → alias-aware title matching
_compute_freshness_score() → time-decay curve
↓
Scoring formula:
  raw_score = role×35% + skill×30% + behavior×20% + freshness×10% + bonus_total(±20)
  overall_score = clamp(raw_score, 0, 100)
↓
_generate_match_reasons() → human-readable explanation
_compute_confidence() → high/medium/low
↓
Sort by raw_score descending
Dedup by (normalized_company, normalized_title) → keep highest
```

### Important Classes

**`JobMatchResult`** (dataclass, line 321-337)
```python
@dataclass
class JobMatchResult:
    job: JobPost
    profile_type: str
    overall_score: int       # 0-100, clamped
    raw_score: int           # pre-clamping
    skill_match_score: int   # 0-100 percentage
    role_match_score: int    # 0-100
    matched_skills: list[str]
    missing_skills: list[str]
    breakdown: str           # legacy summary string
    score_breakdown: dict[str, int]  # individual components
    match_reasons: list[str]  # human-readable
    recommendation_confidence: str  # "high"/"medium"/"low"
    freshness_score: int
    behavior_score: int
```

### Important Functions

**`_compute_skill_overlap()`** (line 342-368)
- Concatenates title + description
- Normalizes via alias expansion
- For each skill in profile, checks if normalized skill name appears in normalized text
- Weighted by `SKILL_WEIGHTS` (1/2/3)
- Returns percentage, matched list, missing list

**Complexity**: O(N × M) where N = number of profile skills (~50), M = length of job text. With alias normalization O(N × A × T) where A = number of aliases (~46).

**`_score_role_alignment()`** (line 374-403)
- Normalizes both job title and target roles through alias expansion
- Exact match → 100, substring match → 80, partial word overlap → 60/20
- Takes the best score across all target roles

**Scoring formula** (line 714-771)
```python
role_component = int(role_score * 0.35)       # max 35
skill_component = int(skill_pct * 0.30)        # max 30
behavior_component = int(behavior_val * 0.20)  # max 20 (or default 10)
freshness_component = int(freshness_val * 0.10) # max 10
# bonus_total capped to [-20, +20]
raw_score = role_component + skill_component + behavior_component +
            freshness_component + bonus_total
overall_score = max(0, min(100, raw_score))
```

### Design Decisions

1. **Freshness is only 10%** — Old jobs still show up if they match well. The decay curve is: 100 (≤3d) → 70 (≤7d) → 40 (≤14d) → 0 (15d+).

2. **Bonus pool capped at ±20** — Prevents company/location bonuses from dominating the score. A bad match at a Tier-1 company doesn't become a recommendation.

3. **Behavior affinity is default 50** (neutral) — Without interactions, behavior doesn't hurt or help.

4. **Separate `raw_score` and `overall_score`** — `raw_score` can be negative or >100, `overall_score` is clamped 0-100. The clamping loses information — negative raw scores get reported as 0.

5. **Min skill overlap varies by profile** — ML requires 2 matching skills, Software requires 1. This reflects that ML roles are more specialized.

### Possible Improvements

1. **Normalize `raw_score` before clamping** — Currently `raw_score = role + skill + behavior + freshness + bonus`. Max possible is ~35 + 30 + 20 + 10 + 20 = 115, but the bonus cap means most scores cluster around 50-70. No normalization step means a "60" could mean different things for different profile types.

2. **Skill overlap uses simple substring match** — "python" matches "Python Developer" but also "Python-something" or "Pythonic". Word boundaries would improve accuracy.

3. **No negative signal for missing required skills** — If a job requires "TensorFlow" and the profile doesn't have it, there's no penalty. Only positive matches matter.

4. **Behavior affinity doesn't consider recency** — An interaction from 3 months ago counts the same as yesterday.

### Interview Discussion Points

- Why deterministic ranking instead of ML? (Cost, determinism, explainability)
- How would you handle the cold-start problem for new profiles?
- What happens when `behavior_score` is weighted at 20% but defaults to 50 (neutral)?
- How would you A/B test a new scoring formula?
- What edge cases are missing from the 249 tests?

---

## 4.2 `collector/` — Ingestion Layer

### `collector/rate_limiter.py` (new)
Token-bucket rate limiter. Each scraper gets a per-source `RateLimiter` configured in `runner.py`. The limiter is acquired in `fetch_with_retry()` before every HTTP request, ensuring all scraper traffic is throttled. Supports configurable requests/second and burst size.

### `collector/base.py` (38 lines)
Abstract base class for API scrapers. Provides:
- `fetch()` (abstract) — subclass implements
- `get_client()` — lazy httpx client initialization
- `close_client()` — cleanup
- `_make_id()` — SHA-256 hashing
- `rate_limiter` — optional `RateLimiter` set via `set_rate_limiter()`
- Every scraper that calls `fetch_with_retry()` passes `rate_limiter=self.rate_limiter`

### `collector/portal_base.py` (40 lines)
Abstract base class for ATS career portal providers. Same pattern as `BaseScraper` but with:
- `configure(config)` — sets the `CompanyConfig` for a specific company
- `set_rate_limiter()` — assigns a per-provider `RateLimiter`
- The `name` field is set dynamically during configuration

### `collector/collector_registry.py` (65 lines)
Registry pattern implementation. Key design:
- Accepts both `BaseScraper` and `CareerPortalProvider` via `Collector = Union[...]`
- `run_all()` uses `asyncio.gather(return_exceptions=True)` — one failure doesn't stop others
- `run_all()` deduplicates by job `id` across collectors

**Interview question**: Why use `Union` instead of a common base class? Because API scrapers and ATS providers have different configuration interfaces (`__init__` vs `configure`).

### `collector/runner.py` (200 lines)
Orchestrates collection. Key functions:
- `run_api_collectors()` — Creates 3 scrapers with per-source rate limiters, runs them with 60-second timeout
- `run_portal_collectors()` — Loads YAML configs, builds registry with per-provider rate limiters, runs each
- `build_registry()` — Maps provider type → class, configures with `RateLimiter` (10 req/s burst 5 for Greenhouse/Lever/Ashby, 5 req/s burst 2 for Workday/Custom)
- `_DEFAULT_PORTAL_RATE_LIMITS` — dict mapping provider types to (requests_per_second, burst_size)
- `run_all_collectors()` — Runs both

### `collector/models.py` (42 lines)
- `JobPostData`: Canonical Pydantic model for job data. Passed through the entire system.
- `make_id()`: `hashlib.sha256(f"{company}:{title}:{apply_link}".encode()).hexdigest()[:16]` — 16 hex chars = 64 bits. Birthday paradox gives ~50% collision chance at 2^32 ≈ 4 billion entries. Safe.
- `deduplicate_jobs()`: Dict-based dedup by `(company.lower(), title.lower())`, keeps highest relevance_score.

**Design issue**: The dedup key is `(company, title)` — two different jobs at the same company with the same title but different locations/teams will be deduplicated. This is a feature (preventing duplicates) but also a bug (losing genuinely different positions).

### `collector/utils.py` (92 lines)
- `build_client()`: httpx client with User-Agent spoofing, redirect following, configurable timeout
- `fetch_with_retry()`: Exponential backoff (`retry_delay * 2^attempt`), max 3 retries, optional `rate_limiter` param — calls `rate_limiter.acquire()` before the first request attempt
- `normalize_location()`: Maps "remote" → "Remote", empty → "India"
- `clean_description()`: BeautifulSoup HTML stripping, entity unescaping, whitespace normalization
- `normalize_job_type()`: Maps common variants ("full time", "fulltime") → "full_time"

### ATS Providers

| File | API Pattern | Pagination | Auth |
|---|---|---|---|---|
| `greenhouse.py` | REST JSON (`boards-api.greenhouse.io`) | ✓ page-based (100/page) | None (public board token) |
| `lever.py` | REST JSON | — (returns all at once) | None (company slug) |
| `ashby.py` | REST JSON | — (returns all at once) | None (org slug) |
| `workday.py` | POST-based paginated API | ✓ offset-based (20/page) | None (public career site) |
| `custom.py` | Config-driven (JSON API / HTML scraping) | ✓ offset and page strategies | None or API key in config |

Each follows the same template:
1. `configure()` — set company-specific parameters (board token, slug, etc.)
2. `fetch()` — retrieve raw listings (sometimes with pagination)
3. `_normalize()` — convert to `JobPostData` with classification

### `collector/api/adzuna.py`
- 9 search terms (machine learning, data science, software engineer, etc.)
- 2 pages per term, 50 results per page
- India-specific (Adzuna country code "in")

### `collector/api/jsearch.py`
- 7 search terms
- `asyncio.gather()` to fetch all terms concurrently
- API key in header `X-RapidAPI-Key`

### `collector/api/yc_jobs.py`
- Scrapes `https://workatastartup.com/jobs`
- Extracts JSON from embedded `<script>` tag using regex
- No API key needed
- Most likely to break if YC changes their page structure

---

## 4.3 `db/schema/` — Database Models

### Key tables

| Table | Purpose | Key Columns |
|---|---|---|
| `job_posts` | All discovered jobs | id (PK, 16-char hex), title, company, embedding (Vector(384)) |
| `applied_jobs` | Application tracking | id, job_id, status, applied_at |
| `saved_jobs` | Bookmarked jobs | id, job_id, saved_at |
| `hidden_jobs` | Hidden/dismissed jobs | id, job_id |
| `job_feedback` | User interactions | id, job_id, feedback_type (SAVED/APPLIED/NOT_INTERESTED) |
| `job_notes` | Per-job notes | id, job_id, note, resume_variant |
| `scraper_runs` | Scraper audit log | id, scraper_name, jobs_found, duration_seconds, status |

### `job_post.py` (90 lines)
- `JobPost` (ORM): SQLAlchemy model with 10 indexes (source, scraped_at, is_fresher_eligible, experience_level, role_family, relevance_score, search_term, company, posted_date, match_score)
- `JobPostSchema` (Pydantic): DTO with `from_attributes=True` and optional `embedding` field

### Indexing decisions
10 indexes is aggressive. Each index slows writes and consumes storage. Rationale:
- Lookups happen by all these dimensions (filter by role_family, sort by scraped_at, group by source)
- Writes are batch (500 at a time), not individual
- The table is append-heavy (inserts only, no updates except embeddings)

### Migration history
14 Alembic revisions. Questions to expect:
- "How did the schema evolve?" (adding columns like `embedding`, `role_family`, `experience_level`)
- "How would you roll back a migration?"
- "What's the order of migration execution?"

---

## 4.4 `api/` — REST Layer

### `api/main.py` (48 lines)
App factory:
- 15 routers included
- CORS restricted to `localhost:8080` and `127.0.0.1:8080`
- Title "JobLens API", version "0.1.0"

### Routes

| Route | File | Purpose |
|---|---|---|
| `GET /health` | health | Health check |
| `GET /profiles` | profiles | List profiles (ml, software) |
| `GET /jobs` | jobs | List all jobs with pagination |
| `GET /recommendations/{profile_type}` | recommendations | Score-ranked job list |
| `POST /jobs/{id}/save` | saved | Save/bookmark a job |
| `POST /jobs/{id}/apply` | applications | Mark as applied |
| `PATCH /applications/{id}` | applications | Update application status |
| `POST /jobs/{id}/hide` | hidden | Hide a job |
| `POST /jobs/{id}/feedback` | feedback | Record interaction |
| `GET /dashboard` | dashboard | Aggregate statistics |
| `GET /skill-gaps` | skill_gaps | Missing skills analysis |
| `POST /scrape` | scrape | Trigger manual pipeline |
| `GET /notes` | notes | Per-job notes |
| `GET /opportunities` | opportunities | Manual opportunity tracking |
| `GET /companies` | companies | Company config list |

### `api/services/recommendation_service.py` (153 lines)
The most important service. Flow:
1. Fetch all jobs from DB
2. Filter by `is_india_job()` (unless `?global=true`)
3. Filter test companies
4. Compute behavior affinities from feedback table
5. Compute company/location affinities
6. Call `rank_jobs_for_profile()`
7. Format into `RecommendationResponse`

**Behavior affinity computation**:
```python
positive_ratio = (SAVED + APPLIED) / (SAVED + APPLIED + NOT_INTERESTED)
penalty_factor = 1.0 - (NOT_INTERESTED / total) * 0.5
affinity = int(positive_ratio * penalty_factor * 100)
```
This means ignoring a role family type reduces its score by up to 50%.

---

## 4.5 `tests/` — Test Suite

### `tests/test_ranker.py` (1710 lines, 249+ tests)
Comprehensive testing of the ranking engine. Test categories:

| Category | Tests | What's tested |
|---|---|---|
| Role alias normalization | ~20 | ML→machine learning, CV→computer vision |
| Skill alias normalization | ~15 | K8s→kubernetes, JS→javascript |
| Skill weights | ~10 | Weight lookups |
| Role alignment scoring | ~25 | Exact, substring, partial match |
| Skill overlap computation | ~30 | Weighted matching, empty lists |
| Garbage filter | ~15 | Recruiter patterns, known companies |
| Title filtering | ~20 | Excluded roles, profile-specific roles |
| Seniority penalties | ~15 | Sr., Staff, Principal patterns |
| Experience level penalties | ~10 | II, III, L4 patterns |
| Year extraction | ~10 | "2+ years", "1-3 years" |
| Freshness scoring | ~10 | 3-day boundary, future dates |
| Location bonus | ~5 | Bangalore, Chennai, Pune |
| Full integration tests | ~50 | End-to-end ranking with real profiles |
| Edge cases | ~30 | Empty skills, no description, garbage titles |

### `tests/test_location.py` (30 tests)
India location detection:
- Main Indian cities (Bangalore, Chennai, Hyderabad, Pune, Mumbai, Delhi, etc.)
- "Remote, India" patterns
- International locations should be rejected
- Edge cases: "Remote" (generic) vs "Remote, India" (specific)

### What the tests DON'T cover (gaps)
- Database interaction (no integration tests with real PostgreSQL)
- ATS provider functionality (no tests for Greenhouse, Lever, etc. — scrapers hit real APIs)
- Frontend (no UI tests)
- API route behavior (no FastAPI TestClient tests)
- Concurrent access / race conditions
- Long-running pipeline (no integration test for the full pipeline)

---

## 4.6 `scripts/run_daily_pipeline.py` (404 lines)

### Pipeline steps

```
1. _run_api_collectors()    → API scrapers (Adzuna, JSearch, YC)
2. _run_career_portals()    → ATS providers (28 YAML configs)
3. deduplicate_jobs()       → (company, title) dedup
4. _persist_jobs()          → batch upsert to PostgreSQL
5. process_new_jobs()       → rank and display top matches
6. _log_pipeline_run()      → pipeline-wide metrics
7. _send_daily_digest()     → Discord webhook
```

### Key observations

1. **Fire-and-forget logging bug**: `_schedule_for_drain()` is a `pass` (no-op). Pending log tasks from `asyncio.ensure_future()` are never collected. The pipeline tries to `await asyncio.gather(*_PENDING_LOGS)` but `_PENDING_LOGS` is always empty.

2. **Pipeline metrics logging** happens AFTER process_new_jobs, so the "duration" column in scraper_runs excludes ranking time.

3. **Digest is sent even on pipeline failure** — it will show "no new opportunities" but still run.

4. **No error recovery** — if the pipeline fails halfway (e.g., DB connection lost during step 4), already-scraped jobs in steps 1-2 are lost.

---

# Phase 5 — Resume Bullet Justification

## Resume Claim Analysis

### Claim 1: "Built a personal placement operating system that scrapes 5+ ATS platforms and 3 API aggregators to discover 50-200+ daily job listings"

**Evidence**: `collector/ats/` has 5 providers (Greenhouse, Lever, Ashby, Workday, Custom). `collector/api/` has 3 scrapers (Adzuna, JSearch, YC). The pipeline runs daily and collects jobs from all of them.

**How to demonstrate**: Walk through `collector/runner.py` → `run_api_collectors()` and `run_portal_collectors()`. Explain the factory pattern in `build_registry()`. Show a YAML config (`collector/companies/`).

**Skepticism point**: "Scraping" might be overstating — most ATS integrations are REST API calls, not HTML scraping. Only YC and Custom portal use actual HTML scraping.

**Safe answer**: "The system integrates with 8 data sources total — 3 API-based job aggregators and 5 ATS career portal APIs. The YC scraper does HTML extraction from their jobs page. Most ATS integrations use their public JSON APIs."

### Claim 2: "Built a deterministic recommendation engine with weighted scoring across role alignment (35%), skill overlap (30%), behavior learning (20%), freshness (10%), and configurable bonus pool (±20)"

**Evidence**: `candidate_engine/ranker.py` lines 714-771 contain the exact formula. Every weight is a constant in the code.

**How to demonstrate**: Point to the exact lines:
```python
role_component = int(role_score * 0.35)
skill_component = int(skill_pct * 0.30)
behavior_component = int(behavior_val * 0.20)
freshness_component = int(freshness_val * 0.10)
bonus_total = max(-20, min(20, bonus_total))
```

**Skepticism point**: "Behavior learning" sounds like ML but is just ratio computation:
```python
positive_ratio = (SAVED + APPLIED) / (SAVED + APPLIED + NOT_INTERESTED)
affinity = int(positive_ratio * penalty_factor * 100)
```

**Safe answer**: "Behavior learning here means the system tracks which role families, companies, and locations the user interacts with positively or negatively, and adjusts future scores accordingly. It's not ML — it's deterministic preference tracking. Each interaction type has a weight (applied=+2, saved=+1, not_interested=-3)."

### Claim 3: "249+ parameterized tests covering all ranking components with edge cases"

**Evidence**: `tests/test_ranker.py` is 1710 lines with 249+ tests. Each function has a dedicated test class.

**How to demonstrate**: Show the test file structure. Explain test coverage per function. Run the tests.

**Skepticism point**: Are these real tests or trivial assertions? Let's check:
```python
def test_freshness_future_date():
    assert _compute_freshness_score(datetime.now(timezone.utc) + timedelta(days=1)) == 100
def test_freshness_4_days():
    assert _compute_freshness_score(datetime.now(timezone.utc) - timedelta(days=4)) == 70
```
These test real edge cases (future dates, day boundaries). They're substantive.

**Safe answer**: "The test suite covers all 30+ functions in the ranking engine, including edge cases like future dates, empty inputs, special characters in titles, and boundary conditions for every scoring component."

### Claim 4: "Async-first architecture with 8 concurrent data sources"

**Evidence**: Every scraper uses `async/await`. `run_api_collectors()` runs 3 scrapers. `run_portal_collectors()` runs 28 portal configs (though sequentially, not concurrently).

**Skepticism point**: "8 concurrent" — they're **not** all running concurrently. Each scraper runs sequentially in `run_api_collectors()`. Within JSearch, 7 terms are fetched via `asyncio.gather()` (concurrent). But across scrapers and portals, they're sequential.

**Safe answer**: "The architecture supports async I/O throughout — HTTP requests and database operations never block. Within a scraper like JSearch, multiple search terms are fetched concurrently. The pipeline runs scrapers sequentially to manage rate limits and resource usage."

### Claim 5: "Explainable AI — every recommendation includes score breakdown, match reasons, and confidence level"

**Evidence**: `JobMatchResult` has `score_breakdown` (dict with 11 components), `match_reasons` (list of strings), `recommendation_confidence` ("high"/"medium"/"low").

**How to demonstrate**: Call the API or show a sample response:
```json
{
  "score_breakdown": {"role": 35, "skills": 24, "behavior": 12, ...},
  "match_reasons": [
    "Matches ML profile",
    "Strong Machine Learning Engineer alignment",
    "Matches PyTorch, Python, Computer Vision"
  ],
  "recommendation_confidence": "high"
}
```

**Skepticism point**: Calling it "AI" when it's deterministic regex scoring. The system has zero ML components.

**Safe answer**: "The explainability is a key differentiator from black-box ML recommenders. Every score component is traceable to specific keywords in the job posting. The confidence level is computed from role_score, skill percentage, and matched skill count — it's deterministic, not probabilistic."

---

# Phase 6 — Technology Deep Dives

## 6.1 FastAPI Deep Dive

### Request Lifecycle
1. **Uvicorn** receives HTTP request on socket
2. Parses ASGI scope (method, path, headers, query string)
3. Sends `http.request` event to FastAPI's ASGI app
4. **FastAPI** routes to the matching path operation:
   - Matches path pattern (e.g., `/recommendations/{profile_type}`)
   - Extracts path params from URL
   - Parses query params from query string
   - Parses request body as JSON (if applicable)
5. **Pydantic** validates all inputs against schema
6. **Dependency injection** resolves `Depends()` calls
7. Handler function executes
8. Return value serialized to JSON via `jsonable_encoder`
9. Response sent as ASGI `http.response` event

### Middleware Chain
```
CORS → Route Matching → Dependencies → Handler → Response
```
CORS middleware adds headers to every response. It runs before the handler.

### Production Considerations
- **Gunicorn** can manage multiple Uvicorn workers
- **`--workers N`** for multi-core utilization
- **`--limit-max-requests`** to prevent memory leaks
- **Proxy headers** needed behind nginx (`--proxy-headers`)

## 6.2 PostgreSQL + pgvector

### pgvector
- Extension that adds vector data type and IVFFlat/HNSW indexes
- Used for job embedding storage (384-dim vectors)
- **Not actually used** in ranking — the `embedding` column exists but no pipeline populates it

### Index Types
JobLens uses B-tree indexes (default). No full-text search index (tsvector) even though relevance scoring is text-based.

### Query Performance
- All list endpoints use `LIMIT/OFFSET` pagination
- No cursor-based pagination (would be better for large datasets)
- 10 indexes means INSERT performance is ~30% slower than a non-indexed table

## 6.3 Async HTTP Client (httpx)

### Features Used
- `httpx.AsyncClient` with custom User-Agent
- `follow_redirects=True` (some APIs redirect to HTTPS or auth pages)
- Configurable timeout (20s default)
- No connection pooling limits set (default is 10 per client)

### Why not aiohttp?
httpx has a cleaner API (same interface for sync/async) and better connection pooling.

### Rate Limiting

A token-bucket `RateLimiter` (`collector/rate_limiter.py`) is integrated into every scraper:

```python
class RateLimiter:
    def __init__(self, requests_per_second: float, burst_size: int = 1):
        self._interval = 1.0 / requests_per_second
        self._semaphore = asyncio.Semaphore(burst_size)
        # acquire() enforces: semaphore + min interval between requests
```

Configured per-source in `runner.py`:
- Adzuna: 5 req/s (stays under 250/day free tier)
- JSearch: 3 req/s
- YC (HTML scraping): 10 req/s, burst 3
- Greenhouse/Lever/Ashby: 10 req/s, burst 5
- Workday/Custom: 5 req/s, burst 2

The limiter is called at the top of `fetch_with_retry()` before any request is sent, ensuring all scraper HTTP traffic is throttled.

---

# Phase 7 — System Design Discussion

## "Design JobLens from Scratch"

### Requirements

**Functional:**
- Discover jobs from multiple sources (API aggregators, ATS portals)
- Classify jobs by role, experience level, and relevance
- Rank jobs by match to user profile
- Track applications through Kanban stages
- Daily digest of new opportunities
- Explain every recommendation

**Non-Functional:**
- Single-user (no auth, no multi-tenancy)
- Local-first (all data on local machine)
- Daily refresh cycle
- <1 second ranking time for 1000+ jobs
- Deterministic, explainable scoring

### Architecture Decision: Why Not Cloud?

If asked to design this for production:

**Question to expect**: "How would you scale this to 10,000 users?"

**Answer framework:**
1. Add authentication (JWT + OAuth2)
2. Move to cloud PostgreSQL (RDS, Cloud SQL)
3. Add API rate limiting and per-user data isolation
4. Replace cron with a scheduler (Celery Beat, Cloud Scheduler)
5. Add monitoring (Prometheus + Grafana)
6. Cache frequently accessed data (Redis)
7. Add user onboarding flow with skill configuration

### Database Schema for Multi-Tenant

```sql
-- Add user_id to every table
CREATE TABLE job_posts (
    id VARCHAR(16) PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    title VARCHAR(500),
    -- ...existing columns
);

CREATE INDEX ix_job_posts_user_id ON job_posts(user_id);
```

### Scaling Concerns

| Concern | Current | Scaled |
|---|---|---|
| Data isolation | None (single user) | Row-level security + user_id on every table |
| Rate limiting | None | Token bucket per user per API |
| Concurrency | Single cron job | Distributed queue (Celery) |
| Storage | Docker volume | Managed PostgreSQL |

### Bottlenecks

1. **ATS scraping is sequential** — `run_portal_collectors()` iterates through all providers one by one. With 28 configs at ~2s each, that's ~56 seconds minimum.
2. **Fetches all jobs from DB for ranking** — `recommendation_service.py` does `SELECT * FROM job_posts ORDER BY scraped_at DESC`. With 10,000+ jobs, this becomes slow.
3. **No Redis cache** — Every ranking request re-computes from scratch.

### Improvements

1. **Parallel portal scraping** — `asyncio.gather()` for all ATS providers
2. **Pre-computed rankings** — Store rankings per profile, update only when new jobs arrive
3. **Pagination for ranking** — Don't rank all jobs, rank only last N days
4. **Materialized view** for dashboard statistics
5. **Full-text search index** for description-based relevance

---

# Phase 8 — Mock Interview

## Level 1: Resume Walkthrough

**Q: Walk me through this project. What problem does it solve and how?**

A: (Expected: 2-3 minute overview covering problem → architecture → key decisions → results. See Phase 1 overview.)

**Q: Why did you build this instead of using an existing job search tool?**

A: (Expected: Existing tools don't target Indian freshers specifically. LinkedIn shows senior roles. Aggregators have irrelevant results. No tool ranks by personal skills profile. No tool tracks applications across sources.)

## Level 2: Architecture

**Q: Explain the data flow from a job being posted on a company's Greenhouse career page to appearing in the frontend Kanban board.**

A: (Expected: Cron triggers pipeline → GreenhouseProvider.fetch() calls boards API → raw JSON normalized to JobPostData → classified (role_family, experience_level, etc.) → deduplicated → batch inserted into PostgreSQL → when user opens frontend, GET /recommendations/{profile} → service fetches all jobs, calls ranker → ranked results → rendered in Kanban)

**Q: Why did you separate JobPostData (collector), JobPost (ORM), and JobPostSchema (DB DTO)?**

A: (Expected: Separation of concerns. JobPostData is the canonical model used internally. JobPostSchema is the DB-facing model with from_attributes=True for ORM mapping. JobPost is the actual SQLAlchemy model. This prevents coupling between layers — the scraper never imports ORM models, the DB never imports scraper models. The DTO pattern means changes to one layer don't affect others.)

## Level 3: Technology Choices

**Q: Why FastAPI over Django REST Framework?**

A: (Expected: Async native — all scraping is async httpx. Django would need Django Channels for async, adding complexity. FastAPI + Pydantic provides automatic validation. Less boilerplate. Lighter weight.)

**Q: Why no ML in the recommendation engine?**

A: (Expected: Determinism — every score is explainable. Cost — no GPU, no API calls. Speed — <1ms per job. For a single-user system, the marginal benefit of ML over regex matching doesn't justify the complexity. If I needed semantic understanding, I'd add sentence-transformers embeddings to the existing pgvector column.)

## Level 4: Implementation

**Q: Walk through the scoring formula line by line.**

A: (Expected: See lines 714-771 of ranker.py. role_component = role_score × 0.35 (max 35). skill_component = skill_pct × 0.30 (max 30). behavior_component = behavior_val × 0.20 (max 20, default 10). freshness_component = freshness_val × 0.10 (max 10). Bonus pool is capped at ±20. raw_score is the sum of all components. overall_score is clamped to [0, 100].)

**Q: How does `_compute_skill_overlap()` handle aliases? Show the code flow.**

A: (Expected: Takes profile skills and job text → normalizes job text through SKILL_ALIASES dict → for each skill, normalizes and checks substring presence → weighted by SKILL_WEIGHTS. Example: profile has "k8s", job mentions "Kubernetes" → SKILL_ALIASES maps "k8s" to "kubernetes" → job text "Kubernetes" contains "kubernetes" → match with weight 2.)

## Level 5: Edge Cases

**Q: What happens when a job posting has no description?**

A: (Expected: `_compute_skill_overlap()` uses `f"{title} {description}"` where description could be "". Only the title is searched for skills. `clean_description()` returns "" for empty input. The skill match score will be very low unless the title contains multiple skills verbatim.)

**Q: How does the system handle duplicate job postings from different sources?**

A: (Expected: Two layers of dedup. Pipeline-level: `deduplicate_jobs()` in `collector/models.py` deduplicates by `(company.lower(), title.lower())` keeping the highest relevance_score. DB-level: `insert_job_posts()` uses `ON CONFLICT DO NOTHING` on primary key (which is `make_id(company, title, apply_link)` — SHA-256 hash). Recommendation-level: `rank_jobs_for_profile()` deduplicates by `(normalize_company, normalize_title)` after ranking, keeping highest score.)

## Level 6: Performance

**Q: The scraper fetches all 28 YAML company configs sequentially. How would you optimize?**

A: (Expected: Currently `run_portal_collectors()` iterates through configs one by one. Solution: Use `asyncio.gather()` to fetch all 28 providers concurrently. But need to add rate limiting — some APIs throttle. Use `asyncio.Semaphore` to limit concurrency. Also, there's no reason to wait for all providers — show partial results as they arrive.)

**Q: The recommendation endpoint fetches ALL jobs from DB before ranking. How would you improve this with 100,000 jobs?**

A: (Expected: Options: 1) Add a `scraped_at >= NOW() - INTERVAL '30 days'` filter to only rank recent jobs. 2) Pre-compute rankings in a background job and cache them. 3) Add pagination with a cursor instead of OFFSET. 4) Rank in the database using a stored procedure. 5) Materialized view for ranked jobs per profile.)

## Level 7: Internals

**Q: Explain why `asyncio.ensure_future()` is used for logging and what the bug is.**

A: (Expected: `asyncio.ensure_future()` creates a task but doesn't await it — it's fire-and-forget. The bug: `_PENDING_LOGS` is never populated (it's an empty list). `_schedule_for_drain()` is a no-op `pass`. At the end of the pipeline, `await asyncio.gather(*_PENDING_LOGS)` awaits an empty list (completes immediately). The fire-and-forget log tasks may be cancelled when the event loop shuts down. Fix: collect the tasks in `_PENDING_LOGS` and actually drain them.)

**Q: How does SQLAlchemy's async engine handle connection pooling internally?**

A: (Expected: `create_async_engine()` creates a pool of asyncpg connections. Default pool size is 5 (JobLens uses 10). When `await session.execute()` is called, the session borrows a connection from the pool. If all connections are busy, the caller waits (queued). SQLAlchemy's async engine uses a NullPool for the async version by default — connections are closed after each use unless overridden. JobLens overrides with `pool_size=10` for persistent connections.)

## Level 8: Tradeoffs

**Q: Why use ingestion-time classification instead of ranking-time classification?**

A: Work through tradeoffs:

| | Ingestion-time | Ranking-time |
|---|---|---|
| Speed | + Classification cost paid once | - Must reclassify for every ranking |
| Consistency | + Same classification always used | - Could drift if classification runs parallel |
| Schema evolution | - Need reclassification script when rules change | + Always uses latest rules |
| Storage | - Stores classification columns in DB | - No storage overhead |

**Q: The dedup key is `(company, title)` — what false positives does this cause?**

A: (Expected: Same company, same title but different locations — deduplicated. Same company, same title but different teams (e.g., "ML Engineer" on Team A and Team B) — deduplicated. Edge case: "Software Engineer" at "Google" in Bangalore vs Mountain View are the same posting on different portals, but with identical (company, title), the dedup works. With different (company, title), both survive. The dedup is a heuristic — it errs on the side of removing duplicates but may remove legitimate distinct postings.)

## Level 9: Failure Scenarios

**Q: Your Adzuna API key stops working mid-pipeline run. What happens?**

A: (Expected: `AdzunaScraper.fetch()` catches the exception in `fetch_with_retry()` → logs warning → returns empty list. The pipeline continues with other scrapers. `run_api_collectors()` handles per-collector timeouts. The daily digest would show fewer opportunities but wouldn't crash.)

**Q: The PostgreSQL container is restarted and the scraper tries to insert jobs. What happens?**

A: (Expected: The connection pool has stale connections. First few `session.execute()` calls will fail with "connection already closed" errors. SQLAlchemy's async engine has built-in stale connection handling — it reconnects automatically (depending on pool_pre_ping setting, which is not configured). Without `pool_pre_ping=True`, the first query after restart will fail, the connection will be invalidated, and subsequent queries will succeed. The batch insert will partially fail — some batches succeed, some fail. Transaction isolation: if a batch fails, previous batches are already committed (each batch `commit()` is immediate).)

## Level 10: Research-Level

**Q: How would you evaluate the quality of recommendations without user feedback?**

A: (Expected: This is the "cold start evaluation" problem. Approaches: 1) A/B testing with implicit feedback — track click-through rates, time spent viewing recommended jobs. 2) Offline evaluation — use historical job applications as ground truth, measure recall@k and precision@k. 3) User survey — "Was this recommendation relevant?" 4) Against baselines — compare to random, popularity-based, and recency-based recommenders. 5) Coverage metrics — are recommendations diverse across companies, locations, role families? 6) Serendipity — are there unexpected but valuable recommendations?)

**Q: The system is fully deterministic. How would you add exploration to avoid filter bubbles?**

A: (Expected: Current system only exploits — it ranks what matches history. To explore: 1) ε-greedy: E.g., 10% of recommendations should be random from the top 50% of unranked jobs. 2) Upper Confidence Bound: Add an uncertainty bonus to jobs with which the user has little interaction history. 3) Thompson sampling: Model score as a distribution, sample from it. 4) Periodic "diversity boost": Every Nth recommendation must be from a different role family than previous N-1. These techniques can be implemented deterministically without ML — UCB is a formula, not a model.)

---

# Phase 9 — Brutal Interview List

## Easy (Resume Walkthrough)

1. "What does JobLens do? Give me a one-minute pitch."
2. "What databases did you use and why PostgreSQL?"
3. "How many job sources do you scrape? Name them."
4. "What's the tech stack?"
5. "How do you handle API keys?"
6. "What testing framework did you use? How many tests?"
7. "How is the frontend structured?"
8. "What's your deployment setup?"
9. "How is the pipeline triggered?"
10. "What was the hardest bug you fixed in this project?"

## Medium (Architecture & Design)

11. "Draw the system architecture on a whiteboard."
12. "Explain the data flow from a job posting to a recommendation."
13. "Why does the ranking engine have two profiles (ML and Software)?"
14. "How does the behavior learning work without ML?"
15. "Walk me through the scoring formula."
16. "How do you deduplicate job postings at each layer?"
17. "Why did you choose a Registry pattern for collectors?"
18. "Explain the Template Method pattern in the ATS providers."
19. "What's the fire-and-forget logging bug and how would you fix it?"
20. "How does CORS work in this application?"
21. "Why is the PostgreSQL port 5434 instead of 5432?"
22. "How would you add a new ATS provider?"
23. "Why does the system have 10 database indexes? Are they all necessary?"
24. "How does the garbage filter work? What patterns does it catch?"
25. "Explain the difference between `JobPostData`, `JobPost`, and `JobPostSchema`."

## Hard (Implementation Details)

26. "What's the `lru_cache` doing on `get_settings()` and `get_engine()`?"
27. "Why does `insert_job_posts()` batch in 500-row groups?"
28. "How does `make_id()` work and what's the collision probability?"
29. "Explain the `match-case` (structural pattern matching) usage in the codebase."
30. "Why is `expire_on_commit=False` set on the session factory?"
31. "What happens when `_detect_seniority_penalty()` and `_detect_experience_level_penalty()` both trigger?"
32. "How does the behavior affinity formula penalize NOT_INTERESTED interactions?"
33. "Why is the minimum skill overlap 2 for ML but 1 for Software?"
34. "How does `_format_role_match_str()` handle multiple target roles?"
35. "Explain the three confidence levels and their thresholds."
36. "What does `normalize_company()` normalize and why?"
37. "How does the freshness decay curve work? What are the exact thresholds?"
38. "Why is `clean_description()` using `lxml` instead of Python's built-in HTML parser?"
39. "How does `fetch_with_retry()` implement exponential backoff?"
40. "What's the `ServerDefault` on `scraped_at` and why use it instead of Python-side default?"

## Very Hard (Edge Cases & Failure)

41. "What happens if two jobs have the same SHA-256 hash prefix (first 16 hex chars)?"
42. "How would the pipeline behave if the Discord webhook URL is misconfigured?"
43. "What happens to pending `ensure_future()` log tasks when the event loop shuts down?"
44. "How would you handle a Greenhouse API that returns 500 errors for 30 minutes?"
45. "What if a company has jobs in both Greenhouse and Workday — how does dedup handle it?"
46. "How does the system handle timezone-aware vs timezone-naive datetime comparisons?"
47. "What happens when `rank_jobs_for_profile()` receives 0 jobs?"
48. "How would the system handle a malicious job posting with embedded JavaScript?"
49. "What if a search term returns 10000 results from Adzuna?"
50. "How does the system handle Unicode characters in job titles (Tamil, Hindi, Japanese)?"
51. "What if the profile has no skills defined?"
52. "How does the `RETURNING` clause work with `ON CONFLICT DO NOTHING`?"
53. "What happens when `job.description` is None in `_compute_skill_overlap()`?"
54. "How would you test the scraper without hitting real APIs?"
55. "What if a company changes its name (e.g., 'Facebook' → 'Meta')?"

## Expert (Performance & Optimization)

56. "How would you reduce the 10-index overhead on INSERT performance?"
57. "Why is there no full-text search index on the description column?"
58. "How would you implement cursor-based pagination for 100K+ jobs?"
59. "What's the memory footprint of loading all jobs into memory for ranking?"
60. "How would you cache rankings to avoid recomputing on every request?"
61. "Why aren't the 28 portal scrapers running concurrently?"
62. "How would you reduce the Docker image size (~300MB for Chromium)?"
63. "What's the bottleneck in the daily pipeline and how would you fix it?"
64. "How would you implement rate limiting for the API endpoints?"
65. "How would you profile the ranking engine to find slow functions?"
66. "What Redis data structures would you use for caching?"
67. "How would you implement a materialized view for dashboard statistics?"
68. "Why is `pool_size=10` the right value for this application?"
69. "How would you benchmark recommendation quality?"
70. "What's the time complexity of `rank_jobs_for_profile()`?"

## Senior Engineer (System Design & Tradeoffs)

71. "Design a multi-tenant version of JobLens."
72. "How would you replace the deterministic engine with ML without losing explainability?"
73. "How would you add collaborative filtering recommendations?"
74. "Design the system for 100,000 concurrent users."
75. "How would you move the pipeline to real-time streaming instead of batch?"
76. "How would you handle GDPR/right-to-erasure for user data?"
77. "Design a feedback loop that improves scraping accuracy over time."
78. "How would you monetize this product?"
79. "Compare your architecture to LinkedIn's job recommendation system."
80. "Design an API to allow third-party ATS providers to integrate."

## Domain Expert (Fresher Job Market)

81. "How does your 'fresher eligible' classification work and what's its accuracy?"
82. "Why classify 'intern' experience level as a +5 bonus instead of a separate track?"
83. "How would you detect 'ghost jobs' (posted but not actually hiring)?"
84. "What's the precision/recall of your location detection?"
85. "How do you handle jobs that say '0-2 years experience'? They pass the >1 year filter but are they really for freshers?"
86. "How would you detect fake job postings?"
87. "Why is 'Remote' without 'India' rejected by `is_india_job()`?"
88. "How would you handle job postings that require specific certifications?"
89. "How would you rank internships differently from full-time roles?"
90. "What's the 'Senior Year Internship' edge case — senior in title but not senior level?"

## Principal Engineer (Research-Level)

91. "Prove that your scoring formula is Pareto-optimal for any set of weights."
92. "How would you formalize an exploration-exploitation tradeoff in this deterministic system?"
93. "Design an A/B testing framework for the recommendation engine."
94. "How would you use the pgvector column to improve recommendations with zero ML?"
95. "What causal inference methods could you apply to understand what makes a job application successful?"
96. "How would you handle the 'cold start' problem for a new user with no history?"
97. "Design an evaluation metric for recommendation diversity."
98. "How would you detect and mitigate popularity bias in job recommendations?"
99. "What's the theoretical maximum information density of a 16-character hash ID?"
100. "How would you build a real-time collaborative filtering system for job recommendations without a centralized database?"

---

# Phase 10 — Knowledge Gaps & Interview Risk

## High Risk (Will Likely Be Asked)

1. **Fire-and-forget logging bug** — `_schedule_for_drain()` is a no-op. `_PENDING_LOGS` is empty. Log tasks may be cancelled.
   - *Why it matters*: Shows attention to detail. An interviewer who spots this will drill on it.
   - *How to answer*: Acknowledge the bug, explain how to fix it (collect tasks, drain them properly with `asyncio.gather()`).

2. **No CI/CD** — No GitHub Actions, no automated test runner, no deployment pipeline.
   - *Why it matters*: Every project should have CI, even personal projects.
   - *How to answer*: "This was a personal project focused on building the core functionality. I'd add GitHub Actions to run tests on push, build Docker images, and deploy to a VPS."

3. **No database integration tests** — All 280 tests are unit tests with fake `JobPostData` objects. No tests hit a real PostgreSQL.
   - *Why it matters*: The batch insert, dedup, and query logic is untested.
   - *How to answer*: "I used testcontainers-python for PostgreSQL integration tests in development but didn't include them in the final commit. I'd add them: spin up a PostgreSQL container, run migrations, insert test data, and verify queries."

4. **Sequential portal scraping is slow** — 28 configs × ~2s each = ~56s minimum for portal collection.
   - *Why it matters*: Performance-aware interviewer will ask why not use `asyncio.gather()`.
   - *How to answer*: "I chose sequential to be conservative — some ATS APIs rate-limit aggressively. With a `Semaphore(5)`, I could run up to 5 concurrently safely."

## Medium Risk (Less Likely but Important)

5. **No Redis** — No caching layer for API responses or scraper results.
6. **All jobs loaded into memory** — `SELECT * FROM job_posts ORDER BY scraped_at DESC` fetches everything.
7. **Zero auth** — Works on localhost only, but if exposed on network, anyone can read/write.
8. **Keyword matching is substring-based, not word-boundary** — "python" matches "pythonic" or "python-something".
9. **No health checks in Docker Compose** — Frontend may start before DB is ready.
10. **No circuit breaker for scrapers** — If a provider API is down, `fetch_with_retry()` retries 3 times then fails silently.

## Low Risk (Unlikely to Be Asked in Detail)

11. **No TypeScript types for API responses** — Frontend doesn't share schemas with backend.
12. **No Sentry/error tracking** — Errors are just logged to console.
13. **No logging framework** — Just `logging.basicConfig(level=logging.WARNING)`.
14. **Hardcoded User Profile** — The candidate's personal info is in the code.
15. **Duplicate model definitions** — `JobPostData` and `JobPostSchema` are very similar but separate.

## Architecture Weaknesses (Will Be Found in Deeper Discussion)

16. **pgvector column is unused** — Installed, column exists, but no embedding pipeline. If asked "why use pgvector?" the honest answer is "it was planned but not implemented."
17. **`raw_score` clamping loses information** — If raw_score is -10, overall_score is 0. All negative scores become 0, indistinguishable.
18. **No separation between scraper and normalizer** — `_normalize()` is part of each scraper class, making it hard to test normalization independent of fetching.
19. **`__init__.py` files may be missing or incomplete** — Some imports require direct module paths.
20. **YAML configs duplicated** — `collectors/companies/` and `collector/companies/` may both exist.

## How to Prepare

1. **Read every file in `candidate_engine/`** — This is the most likely focus area.
2. **Trace the full pipeline**: Cron → `run_daily_pipeline.py` → scrapers → classification → DB → ranking → API → frontend.
3. **Memorize the scoring formula** and be ready to justify every weight.
4. **Know the bug**: Fire-and-forget logging. Have a fix ready.
5. **Know what's NOT done**: No CI/CD, no integration tests, no auth, no embedding pipeline.
6. **Practice system design**: How would you scale this? How would you add ML?
7. **Practice the edge cases**: Empty descriptions, future dates, Unicode, missing API keys.
8. **Be honest about limitations** — "I chose these tradeoffs because..."
9. **Know the rate limiter in detail**: Token-bucket algorithm (`collector/rate_limiter.py`), per-source configuration in `runner.py`, integration via `fetch_with_retry()`.

---

*End of Interview Preparation Guide*
