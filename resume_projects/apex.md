# ApexS — Complete Interview Preparation Guide

---

## Phase 1 — Architectural Overview

**Project Purpose**: ApexS is an Explainable Sprint Planning platform that uses Integer Linear Programming (ILP) to mathematically optimize sprint backlog selection. It replaces gut-feeling sprint planning with constrained optimization that maximizes business value while respecting capacity, risk, skill, and dependency constraints.

**Main Problem**: Engineering managers and scrum masters have no tool to mathematically evaluate which stories to include in a sprint. Manual planning leads to unbalanced sprints, over-commitment, and suboptimal value delivery.

**Architecture**: Three-tier web application with async job processing.

```
┌─────────────────────────────────────────────────────┐
│                  Frontend (React + TS)               │
│  Dashboard → Upload → Configure → Processing →       │
│  Plan → Explain → Approve → Reports                  │
└──────────────────────┬──────────────────────────────┘
                       │ HTTP (Axios, 2s polling)
                       ▼
┌─────────────────────────────────────────────────────┐
│           Backend (FastAPI + Async SQLAlchemy)        │
│                                                       │
│  API Layer: /api/v1/{auth,teams,datasets,sprints,     │
│                     stories,plans,reports,context}     │
│                                                       │
│  Services: OptimizationEngine  (PuLP MILP + Greedy)  │
│            ExplainabilityEngine (rule-based)         │
│            WeightLearningModel  (LogisticRegression) │
│            ContextExtractor    (heuristic stats)     │
│            Preprocessing       (normalization)       │
│                                                       │
│  Workers: Celery task + Thread fallback               │
└──────────┬──────────────┬──────────────┬────────────┘
           │              │              │
     PostgreSQL     Redis (broker)    MinIO (S3)
```

**Folder Structure**:
- `backend/app/api/v1/` — 8 router modules (auth, teams, datasets, sprints, stories, plans, reports, context)
- `backend/app/core/` — config, database, security, minio_client, auth_backend
- `backend/app/models/` — 8 SQLAlchemy models (Team, Sprint, UserStory, SprintPlan, Explanation, DatasetUpload, Context, User)
- `backend/app/schemas/` — Pydantic request/response models
- `backend/app/services/` — 5 service modules (optimization, explainability, weight_learning, context_extraction, preprocessing)
- `backend/app/workers/` — Celery app + planning task with thread fallback
- `backend/migrations/` — Alembic (1 initial migration, 7 tables)
- `frontend/src/pages/` — 9 pages (Dashboard, DatasetUpload, SprintConfiguration, OptimizationProcessing, GeneratedSprintPlan, ExplainabilityPanel, SprintPlanApproval, Reports, NotUsed)
- `frontend/src/hooks/` — 3 hooks (useDatasetUpload, useExplanations, usePlanStatus)
- `frontend/src/store/` — Zustand sprint store
- `tests/` — 19 unit tests + 8 integration tests
- `scripts/` — Experiment runners for academic paper
- `tawos/` — TAWOS dataset pipeline integration

**Key Design Decisions**:
1. PuLP MILP with CBC solver (not OR-Tools) for optimization
2. Sync SQLAlchemy (not async) despite resume claiming async
3. Rule-based explainability (not SHAP) for transparency
4. Celery + Redis with thread fallback for async processing
5. MinIO S3 with local filesystem fallback for storage
6. Frontend polls every 2 seconds for job status
7. Weight learning uses LogisticRegression with fallback chain (context → equal weights)

**Performance**: MILP has 15-second solver time limit, single thread. Scales to ~2000 stories. Greedy fallback is O(n²) in worst case.

**Security**: JWT auth (optional, disabled by default), bcrypt passwords, role-based access (scrum_master, product_owner).

**Testing**: pytest with SQLite test DB, Celery disabled in tests, 27 total tests.

**CI/CD**: GitHub Actions runs tests on push/PR with PostgreSQL + Redis service containers. Docker Compose for dev/prod deployment.

---

## Phase 2 — Knowledge Dependency Graph

```
Beginner
├── HTTP / REST / JSON APIs
│   └── FastAPI
│       ├── ASGI (uvicorn)
│       │   └── Python async vs sync
│       ├── Pydantic (validation, serialization)
│       └── OpenAPI / Swagger docs
├── SQL / Relational Databases
│   └── PostgreSQL
│       └── SQLAlchemy ORM
│           ├── Declarative models
│           ├── Sessions / Unit of Work
│           ├── Relationships (FKs, back_populates)
│           └── Alembic migrations
├── React (hooks, components, props, state)
│   ├── React Router (SPA routing)
│   ├── TanStack React Query (server state, polling)
│   ├── Zustand (client state management, persist)
│   └── Tailwind CSS / shadcn-ui
├── Git / GitHub / GitHub Actions (CI)
└── Docker / Docker Compose
    └── Container orchestration / multi-service

Intermediate
├── PuLP / MILP / ILP
│   ├── Linear Programming
│   │   ├── Objective functions
│   │   ├── Constraints (inequality, equality)
│   │   └── Decision variables (binary, integer, continuous)
│   ├── MILP vs LP vs Integer Programming
│   ├── Constraint satisfaction vs optimization
│   ├── CBC solver (COIN-OR branch-and-cut)
│   ├── Solver time limits / optimality gaps
│   └── Greedy / heuristic fallbacks
├── Celery (task queue)
│   ├── Broker (Redis / RabbitMQ)
│   ├── Result backend
│   ├── Task states (PENDING, STARTED, PROGRESS, SUCCESS, FAILURE)
│   └── Worker processes
├── JWT (JSON Web Tokens)
│   ├── Stateless authentication
│   ├── Token structure (header, payload, signature)
│   └── Refresh vs access tokens
├── scikit-learn
│   ├── LogisticRegression
│   ├── Train/test split
│   ├── Feature scaling (StandardScaler)
│   └── Metrics (accuracy, F1, ROC-AUC, MAE, R²)
├── pandas / NumPy
│   ├── DataFrame operations (groupby, merge, concat)
│   └── CSV parsing / data cleaning
├── MinIO / S3-compatible storage
│   └── boto3 SDK
├── Vite (build tool)
│   └── HMR, proxy, TypeScript
└── Recharts / D3 (data visualization)

Advanced
├── Combinatorial Optimization
│   ├── NP-hard problems
│   ├── Knapsack problem (the optimization model reduces to this)
│   ├── Branch-and-bound / Branch-and-cut
│   └── LP relaxation
├── Asynchronous Task Processing
│   ├── Thread vs Process vs Celery worker
│   ├── In-memory vs Redis-backed state
│   └── Fallback patterns (circuit breaker, degrade)
├── Explainability / Interpretable ML
│   ├── Rule-based explanations vs SHAP/LIME
│   ├── Score decomposition
│   └── Feature importance from coefficients
├── MLOps basics
│   ├── Online vs batch learning
│   ├── Concept drift
│   └── Model serving patterns
└── UML / Entity-Relationship modeling

Research-level
├── Sprint planning as a constraint satisfaction problem
├── Multi-objective optimization (Pareto frontier)
├── ML for effort estimation / prioritization
├── SHAP / LIME for model-agnostic explanations
└── Production ML pipeline reliability
```

---

## Phase 3 — Concept Curriculum

### 3.1 HTTP / REST / FastAPI

**What it is**: FastAPI is a modern Python web framework for building REST APIs with automatic OpenAPI documentation, Pydantic validation, and async support.

**Why it exists**: Traditional Python web frameworks (Flask, Django) required manual validation, serialization, and documentation writing. FastAPI generates all of this from Python type hints.

**How it works internally**: FastAPI runs on ASGI (Asynchronous Server Gateway Interface) via uvicorn. When a request arrives:
1. uvicorn parses the HTTP request into an ASGI scope
2. FastAPI matches the path/verb to a route handler
3. Pydantic validates request body/params from type hints
4. The handler executes (sync handlers run in a threadpool)
5. Return value is auto-serialized to JSON via `jsonable_encoder`

**Where in the repo**: `backend/app/main.py` creates the `FastAPI` app. All routers are in `backend/app/api/v1/`. Pydantic schemas in `backend/app/schemas/common.py`.

**Why chosen**: Performance (on par with Node.js/Go for JSON APIs), automatic docs, type safety, Pythonic.

**Alternatives**: Flask (simpler, less performant), Django (heavy, ORM-coupled), Express.js (non-Python).

**Tradeoffs**: FastAPI's async is best for I/O-bound work. Synchronous database calls still block. FastAPI's dependency injection system is less explicit than Django's.

**Common misconception**: "FastAPI automatically makes everything async." No — sync route handlers like `def get_plan()` run in a threadpool and can block.

**Interview follow-ups**:
- How does FastAPI handle sync vs async route handlers?
- How does dependency injection work in FastAPI?
- What happens if you define a sync route with async database calls?
- How does FastAPI generate OpenAPI docs?

### 3.2 PuLP / MILP

**What it is**: PuLP is a Python linear programming library that allows declarative formulation of optimization problems. It wraps solvers like CBC (COIN-OR Branch-and-Cut).

**Why it exists**: Writing optimization problems directly in solver-specific formats is error-prone. PuLP provides a Pythonic DSL (Domain Specific Language) where you define variables, objective, and constraints in pure Python.

**How it works internally** (in this repo):
1. Decision variables: `x = {s.story_id: LpVariable(f"x_{s.story_id}", cat=LpBinary)}` — one binary variable per story (1=selected, 0=rejected)
2. Objective: `prob += lpSum(x[s.story_id] * self._story_score(s, weights) for s in stories)` — maximize total weighted score
3. Capacity constraint: `prob += lpSum(x[s.story_id] * s.story_points for s in stories) <= capacity`
4. Dependency constraints: For each dependency (A→B), `x[B] <= x[A]` — if B selected, A must also be selected
5. PuLP generates a `.lp` or `.mps` file, passes it to CBC solver
6. CBC performs branch-and-cut: relaxes integer constraints to LP, solves, branches on fractional variables, adds cutting planes
7. Returns solution status (Optimal, Infeasible, Unbounded, TimeLimit)

**Where in the repo**: `backend/app/services/optimization_engine.py:289-319` (`_milp_select`). The 15-second time limit is at line 312.

**Why chosen**: PuLP is free, well-documented, Python-native, and handles small-to-medium problems well.

**Alternatives**: Google OR-Tools (CP-SAT solver, generally faster for combinatorial problems), scipy.optimize (no integer support), commercial solvers (Gurobi, CPLEX — expensive).

**Tradeoffs**: CBC solver is single-threaded by default. PuLP can be slow for >2000 stories. No optimality gap guarantee. PuLP v3+ has breaking changes.

**Common misconception**: "PuLP is a solver." No, PuLP is a modeling layer. CBC is the solver. PuLP can also use Gurobi, CPLEX, GLPK, etc.

**Interview follow-ups**:
- Write the MILP formulation for this problem on a whiteboard.
- How does branch-and-bound work?
- What is LP relaxation?
- When would MILP fail to find optimal solution?
- How would you handle >10000 stories?
- What if capacity constraint can never be satisfied?

### 3.3 Celery + Redis

**What it is**: Celery is a distributed task queue for asynchronous job execution. Redis serves as both message broker (delivering tasks to workers) and result backend (storing task results).

**Why it exists**: The ILP solver is CPU-bound and takes seconds to minutes. Blocking the HTTP response until optimization completes would cause timeouts. Celery offloads this to a worker process, allowing the API to return immediately with a `job_id` for polling.

**How it works internally**:
1. API endpoint calls `run_async_job()` which calls `celery_task.delay()` on the Celery task
2. Celery serializes the task arguments (typically JSON) and publishes to Redis
3. Celery worker (separate process) picks up the task from Redis
4. Worker executes `execute_planning_pipeline()` which progresses through 6 stages
5. Each stage calls `self.update_state()` with progress metadata
6. Frontend polls `GET /plans/status/{job_id}` every 2 seconds
7. Status endpoint reads from `AsyncResult` (Celery) or in-memory `_JOB_STORE` (thread fallback)
8. On completion, plan_id is stored and frontend redirects to plan page

**Where in the repo**: `backend/app/workers/celery_app.py` (config), `planning_task.py:335-356` (task definition), `planning_task.py:359-387` (fallback).

**Why chosen**: Industry standard, well-understood, Redis is already in the stack.

**Alternatives**: RQ (simpler, Redis-only, no monitoring UI), Dramatiq (newer, simpler), Arq (async-only), background threads (no durability).

**Tradeoffs**: Celery adds operational complexity (need a worker process, Redis). Task state is not persisted in DB (in-memory `_JOB_STORE` is lost on restart). Thread fallback means tasks disappear if backend process restarts.

**Common misconception**: "Celery makes tasks durable." Only if you configure a persistent result backend. This project's in-memory `_JOB_STORE` is lost on restart.

**Interview follow-ups**:
- What happens if the Celery worker dies mid-task?
- How does Celery handle task idempotency?
- What's the difference between Celery task states?
- How does the retry backoff work?
- What happens if all retries are exhausted?
- How does `task_acks_late` interact with retry?
- What is a dead letter queue?

### 3.4 Logistic Regression (Weight Learning)

**What it is**: A statistical model that predicts a binary outcome (sprint_completed = 0 or 1) from features [story_points, business_value, risk_score]. The model coefficients are interpreted as feature importance.

**Why it exists**: Teams don't know their optimal prioritization weights. The model learns from historical sprint data what factors (small stories, high-value stories, low-risk stories) correlated with completion.

**How it works internally**:
1. Features X = [story_points, business_value, risk_score], target y = sprint_completed (binary)
2. StandardScaler normalizes features (mean=0, std=1)
3. LogisticRegression with liblinear solver fits: P(y=1) = 1 / (1 + exp(-(w·x + b)))
4. Coefficients w = [w₁, w₂, w₃] represent log-odds per feature
5. Sign-mapping: story_points has negative coefficient (smaller stories → more likely completed), business_value positive, risk_score negative
6. `_coefficients_to_weights()` maps coefficients to [urgency, value, alignment] weights using sign alignment
7. Fallback chain: sklearn unavailable → context weights → equal weights (0.33, 0.34, 0.33)

**Where in the repo**: `backend/app/services/weight_learning.py`. `_coefficients_to_weights()` at line 89. `train_with_metrics()` at line 118.

**Why chosen**: Simple, interpretable, fast, works with small data. Logistic regression coefficients have a direct probabilistic interpretation.

**Alternatives**: Random Forest (better accuracy, less interpretable), Gradient Boosting (XGBoost/LightGBM), Neural Networks (overkill), Ridge Regression (continuous target).

**Tradeoffs**: Logistic regression assumes linear decision boundary. Features are just 3 — may miss complex patterns. Only works if there's variation in `sprint_completed`. Minimum 10 samples required (arbitrary threshold).

**Common misconception**: "The model learns prioritization weights directly." No — it learns a completion prediction model. Weights are derived from coefficients via sign-mapping, which is a heuristic.

**Interview follow-ups**:
- Why use StandardScaler?
- What does the sign-mapping accomplish?
- Why require at least 10 samples?
- What happens if all stories are completed (y all 1)?
- How would you validate that learned weights improve planning?

### 3.5 Explainability Engine (Rule-based)

**What it is**: Generates per-story natural language explanations for why each story was selected or rejected.

**Why it exists**: The MILP optimization is a black box. Scrum teams need to trust the plan's recommendations. Explainability builds trust and allows humans to override when appropriate.

**How it works internally**:
1. For **selected stories**: Calculate score = urgency_contribution + value_contribution + alignment_contribution. Build reason string from components plus heuristic rules (e.g., "high business value" if >= 7, "low risk" if < 0.3, "low effort" if <= 5 points).
2. For **rejected stories**: Priority-ordered rejection reason detection:
   - Risk > threshold → "Risk exceeds threshold"
   - Skill mismatch → "Required skill not available"
   - Status non-plannable → "Status not plannable"
   - Dependency not satisfied → "Dependencies not satisfied"
   - Would exceed capacity → "Would exceed sprint capacity"
   - Fallthrough → "Lower priority under current objective"

**Where in the repo**: `backend/app/services/explainability_engine.py`. Generated in `generate()` at line 52.

**Why chosen**: Rule-based is deterministic, auditable, and doesn't require additional ML dependencies. Every explanation maps to a specific constraint.

**Alternatives**: SHAP (model-agnostic, but requires tree/linear model), LIME (local approximations, but slow), Counterfactual explanations (complex).

**Tradeoffs**: Rule-based cannot explain interaction effects. The priority order of rejection reasons is arbitrary (first match wins). Confidence score is just the normalized objective score — not a statistical confidence.

**Common misconception**: "This is SHAP-style explainability." It explicitly is not — `"shap_enabled": False`. It's a score breakdown, not Shapley values.

**Interview follow-ups**:
- Why is rule-based explainability preferred over SHAP here?
- What's the difference between local and global explainability?
- How would you explain a story that was rejected for multiple reasons?

### 3.6 Sync SQLAlchemy in FastAPI

**What it is**: SQLAlchemy is the Python ORM. This project uses synchronous sessions (`SessionLocal()`) despite running inside an async web framework.

**Why it exists**: Simplicity. Sync SQLAlchemy is easier to debug, has more examples, and works with more libraries.

**How it works internally**: FastAPI detects sync route handlers (`def` not `async def`) and runs them in a threadpool executor. The threadpool has a default of 40 threads. Each DB call blocks its thread but not the event loop.

**Key detail**: `connect_args = {"check_same_thread": False}` for SQLite — required because FastAPI threadpool may use different threads for different requests.

**Where in the repo**: `backend/app/core/database.py`. All routers use `async def` handlers with `AsyncSession` from `Depends(get_async_db)`.

**Tradeoffs**: Async SQLAlchemy avoids threadpool saturation under load. Each request yields the event loop during I/O, allowing thousands of concurrent DB operations without consuming a thread per request. The downside is SQLAlchemy 2.0 async requires `select()` style queries instead of the 1.x `Query` API, which can be verbose.

**Resume claim**: ✅ **Implemented**. The `database.py` exports both `get_db` (sync, used by Celery worker) and `get_async_db` (async, used by all FastAPI routes). Routes use `create_async_engine` with `aiosqlite`/`asyncpg`, and all handlers are `async def` with `await db.execute(select(...))`.

**Interview follow-ups**:
- What's the benefit of async SQLAlchemy over sync in FastAPI?
- How did you handle the migration from sync to async?
- How does Celery still use sync DB while FastAPI uses async?
- What's the difference between `select()` and `Query` APIs?

### 3.7 Frontend Polling (TanStack Query)

**What it is**: A React hook that polls the backend status endpoint every 2 seconds until the async job completes.

**Why it exists**: WebSocket support would add complexity. Polling is simple, reliable, and the status endpoint is lightweight.

**How it works internally**: `usePlanStatus` uses TanStack Query's `refetchInterval` option. When status becomes "complete" or "failed", the component redirects to the next page. The polling stops automatically via TanStack Query's `enabled` option.

**Where in the repo**: `frontend/src/hooks/usePlanStatus.ts`. The 2-second interval: `refetchInterval: 2000`.

**Tradeoffs**: 2-second polling means up to 2 seconds of latency between job completion and UI update. Each poll is a full HTTP request (headers, auth check, DB query). WebSockets would be instant and consume less bandwidth.

**Alternatives**: WebSockets (real-time but requires connection management), Server-Sent Events (unidirectional, simpler), HTTP/2 Server Push (rarely used).

**Interview follow-ups**:
- How does TanStack Query's `refetchInterval` work?
- What happens if the browser tab is backgrounded?
- How would you implement exponential backoff?

---

## Phase 4 — Repository Deep Dive

### 4.1 `backend/app/services/optimization_engine.py` (520 lines)

**Purpose**: The core optimization logic. Contains the ILP/MILP solver and greedy fallback.

**Classes**:
- `OptimizationResult` (dataclass, line 24): Return type containing selected/rejected stories, total_value, total_risk, capacity_used, solver_status, runtime_ms, score_distribution, skill_coverage. 18 fields total.
- `OptimizationEngine` (line 44): Main solver with configurable constraint enforcement flags.

**Key methods**:
- `solve()` (line 321): Entry point. Preprocesses stories, filters feasible, decides MILP vs greedy, returns result.
- `_milp_select()` (line 289): Formulates and solves PuLP MILP problem. 15-second time limit.
- `_greedy_selection()` (line 255): Sorts by score descending, iteratively selects stories that fit capacity and dependencies. O(n²) worst case.
- `_filter_feasible_stories()` (line 200): Applies risk, skill, dependency, status filters in priority order.
- `_score_components()` (line 81): Normalizes [0,1] and computes weighted score.
- `_preprocess_stories()` (line 191): Normalizes status, skill, depends_on.

**Design decisions**:
- Filtering is sequential (risk → skill → dependency → status) — order matters. If a story violates multiple constraints, only the first violation is reported.
- Greedy fallback on MILP failure — guarantees a (suboptimal) solution always exists.
- `enforce_*` flags allow ablation studies for the academic paper.
- `solve_baseline()` (line 433) provides 4 baseline modes for experiments.

**Complexity**:
- MILP: Exponential worst-case (NP-hard), but typically O(n²) for this knapsack-like problem with 15s limit.
- Greedy: O(n²) where n = feasible story count.
- Filtering: O(n).

**Possible improvements**:
- Add optimality gap reporting (MILP often finds near-optimal quickly).
- Use OR-Tools CP-SAT for better performance.
- Add warm start with greedy solution.
- Parallel solver runs with different seeds.

**Interview discussion points**:
- The greedy fallback is deterministic (same seed → same result). Is determinism desirable?
- Score distribution includes min/max/mean but not standard deviation.
- Capacity constraint uses <= but sprint overfill might be acceptable.

### 4.2 `backend/app/services/explainability_engine.py` (125 lines)

**Purpose**: Generate per-story explanations for optimization results.

**Key method**: `generate(result, weights)` (line 52). Two loops:
1. Selected stories (line 57): Build reason from score components + heuristic thresholds
2. Rejected stories (line 92): Priority-ordered rejection reason detection

**Design decisions**:
- Selected story reasons are verbose and informational. Rejected story reasons are actionable (tells WHY not selected).
- Priority order of rejection reasons: risk → skill → status → dependency → capacity → priority. This is hardcoded and arbitrary.
- `confidence_score` is just the normalized objective score, reused as confidence. This is misleading — it's not a statistical confidence.

**Problems**:
- `rejection_reason` and `reason` are duplicated for rejected stories (both attributes contain the same string).
- The rejection reason logic doesn't account for the case where multiple constraints simultaneously caused rejection.
- No explanation for why Story A was selected over Story B (comparative explanation).

**Interview discussion points**:
- Confidence score: what does it actually represent?
- How would you explain a story rejected by both risk AND capacity?
- How would you implement comparative explanations?

### 4.3 `backend/app/services/weight_learning.py` (225 lines)

**Purpose**: Learn prioritization weights from historical sprint data.

**Key method**: `train_with_metrics(df, context)` (line 118). Fallback chain:
1. Empty dataset → context weights
2. sklearn unavailable → context weights
3. < 10 samples → context weights
4. Single class target → context weights
5. Normal training → LogisticRegression → sign-mapped weights

**Key insight**: The sign-mapping at line 89-112 converts ML coefficients to interpretable weights. For story_points: `sign_alignment = -1.0` because smaller stories should have higher urgency. This is a heuristic — not derived from data.

**Metrics tracked** (line 207-221): sample_count, train_count, test_count, MAE, R², accuracy, F1, ROC-AUC, feature_importance, feature_coefficients, model_type.

**Issues**:
- The 10-sample minimum is arbitrary. Logistic regression can work with as few as 5 samples per feature.
- `stratify=y` can fail (line 163) and falls back to non-stratified split, which could create test sets with no positive examples.
- Only 3 features. No interaction terms. No categorical encoding.
- MAE is computed on probability (continuous), not binary prediction — should compute on y_pred, not y_prob.
- R² on probabilities is unusual.

**Interview discussion points**:
- Why sign-mapping instead of using coefficients directly as weights?
- What does it mean when all coefficients are near zero?
- How would you add feature interactions?
- Is 80/20 split appropriate for small datasets?

### 4.4 `backend/app/services/context_extractor.py` (52 lines)

**Purpose**: Extract team context statistics from historical data.

**Key method**: `extract(df, team_capacity)` (line 23). Computes:
- `velocity`: mean story_points per sprint (grouped by sprint_id)
- `completion_rate`: mean of sprint_completed
- `skill_distribution`: normalized value counts of required_skill
- `avg_risk_tolerance`: mean risk_score of completed stories
- `value_completion_correlation`: Pearson correlation between business_value and sprint_completed

**Design decisions**:
- `urgency = velocity / capacity`: Normalizes velocity to [0, 1]. A team that delivers exactly capacity gets urgency=1.
- `value_weight = correlation + 0.5`: Converts correlation [-1, 1] to weight [0.3?, 1]. If correlation=0, value_weight=0.5.
- `alignment = 1 - avg_risk`: Low-risk teams get high alignment weight.

**Issues**:
- `score_base = max(team_capacity, 1)`: If capacity is 0, uses 1 to avoid division by zero.
- The correlation heuristic `correlation + 0.5` is arbitrary.
- `urgency = min(1.0, velocity / score_base)`: Clips to 1.0 — a team that over-delivers gets urgency=1.

**Interview discussion points**:
- How would you validate that these heuristics produce good weights?
- What if the team has no historical data?
- How does `value_completion_correlation` behave with noisy data?

### 4.5 `backend/app/workers/planning_task.py` (387 lines)

**Purpose**: Orchestrates the end-to-end planning pipeline. Contains Celery task definition, thread fallback, and job state management.

**Key components**:
- `_JOB_STORE` (line 32): In-memory dictionary of job states, protected by `Lock`. Celery tasks also write here for unified polling.
- `get_job_state()` (line 41): Reads from `_JOB_STORE` first, then falls back to Celery `AsyncResult`.
- `load_dataset()` (line 82): Supports both local CSV and S3:// paths.
- `load_team_historical_dataset()` (line 93): Loads all historical uploads for a team, deduplicates, concatenates.
- `upsert_sprint_stories_from_dataset()` (line 183): Bulk inserts/updates stories from dataset CSV.
- `save_plan_to_db()` (line 255): Creates SprintPlan + Explanation records.
- `execute_planning_pipeline()` (line 279): The actual pipeline with 6 stages and progress reporting.
- `run_async_job()` (line 359): Dispatcher — tries Celery, falls back to thread.

**Pipeline stages and progress**:
```
10% → loading dataset
20% → syncing stories to DB
35% → loading historical data
50% → context extraction
65% → weight learning
80% → optimization
90% → explainability
100% → save plan
```

**Design decisions**:
- Progress percentages are hardcoded heuristics. Not proportional to actual time.
- In-memory `_JOB_STORE` means jobs state is lost on backend restart.
- `run_async_job()` raises RuntimeError if Celery is configured but unavailable AND thread fallback is disabled.
- The Celery task is defined inside `if celery_app is not None` block (line 335) — this means the task isn't visible to Celery workers during import if the import fails.

**Issues**:
- Retry logic: `max_retries=3` with exponential backoff, jitter, `task_acks_late`. All `Exception` types trigger auto-retry.
- `bulk_insert_mappings` and `bulk_update_mappings` are SQLAlchemy 1.x API, deprecated in 2.0 in favor of `session.bulk_insert()`.
- Thread fallback means if the backend pod restarts mid-job, the job disappears.
- No timeout for the thread fallback — it runs indefinitely.

**Interview discussion points**:
- Why use both in-memory and Celery-based state storage?
- What happens if the Celery worker crashes mid-pipeline?
- How would you implement retry?
- How would you make job state persistent?

### 4.6 `backend/app/core/minio_client.py` (57 lines)

**Purpose**: Abstraction over MinIO S3-compatible storage with local filesystem fallback.

**Key functions**:
- `get_s3_client()` (line 12): Creates boto3 S3 client configured for MinIO with aggressive timeouts (connect 1s, read 2s, no retries).
- `save_bytes(path, data)` (line 34): Tries S3 put_object, falls back to local filesystem write.
- `read_bytes(path)` (line 47): Tries S3 get_object, falls back to local filesystem read.
- `ensure_bucket()` (line 23): Creates bucket if not exists.

**Design decisions**: Aggressive timeouts ensure quick fallback to local storage. No retries means transient MinIO failures cause immediate local fallback.

**Issues**: The fallback is silent — no warning logged. Local paths are relative to CWD (`./storage`), which varies between services. No cleanup mechanism for local storage.

### 4.7 API Routers (8 files in `backend/app/api/v1/`)

**`plans.py`** (138 lines): The most complex router. Endpoints:
- `POST /generate`: Triggers async planning. Requires sprint + dataset to exist. Returns `job_id`.
- `GET /status/{job_id}`: Polling endpoint. Returns status/progress/step/plan_id/error.
- `GET /{plan_id}`: Get plan details.
- `PUT /{plan_id}/approve`: Sets plan status to "approved".
- `POST /{plan_id}/export`: CSV or JSON export.
- `PUT /{plan_id}/modify`: Reruns planning with different parameters.
- `GET /{plan_id}/explain`: List explanations with optional `selected` filter, pagination.
- `GET /{plan_id}/explain/{story_id}`: Single story explanation.
- `GET /{plan_id}/stories`: List selected stories.

**Design decisions**: Role-based auth via `require_roles("scrum_master", "product_owner")`. Pagination for explanations (default 500, max 5000).

### 4.8 Frontend Pages

**App.tsx**: 9 routes: `/` (Dashboard), `/upload`, `/configure`, `/optimizing/:jobId`, `/plan/:planId`, `/explain/:planId`, `/approve/:planId`, `/reports/:teamId`, catch-all redirects to `/`.

**`usePlanStatus` hook**: Polls `GET /plans/status/{jobId}` every 2s using TanStack Query `refetchInterval`. Stops polling on complete/failed.

**Zustand store** (`sprintStore.ts`): Client state with `persist` middleware. Stores current sprint, team, plan state across page navigations.

---

## Phase 5 — Resume Bullet Justification

### Bullet 1: "Built an explainable sprint planning platform using ILP optimization"

**Code evidence**: `optimization_engine.py` — PuLP MILP formulation with capacity, risk, skill, dependency constraints. `explainability_engine.py` — per-story rule-based explanations.

**How to demonstrate**: "I implemented the MILP formulation at `optimization_engine.py:289-319` using PuLP. The objective maximizes weighted story scores subject to four constraint types. If the MILP solver times out or is unavailable, a greedy fallback at line 255 guarantees a valid plan."

**Skepticism point**: Interviewer may ask "Why PuLP over OR-Tools?" or "Write the MILP formulation." Must be able to derive the knapsack formulation from scratch.

**Evidence**: The solve_baseline() function at line 433 provides 4 ablation baselines (fixed_weight, context_only, greedy, random) used for the academic paper.

### Bullet 2: "ML-assisted weight learning using scikit-learn"

**Code evidence**: `weight_learning.py` — LogisticRegression with feature scaling, train/test split, 5 metrics.

**How to demonstrate**: "The `train_with_metrics()` method at line 118 trains a LogisticRegression on [story_points, business_value, risk_score] to predict sprint completion. Coefficients are sign-mapped to [urgency, value, alignment] weights. I implemented a three-level fallback: sklearn unavailable → context weights → equal weights."

**Skepticism point**: Interviewer may ask "Why logistic regression and not XGBoost?" or "How do you validate that learned weights are better than default?"

**Evidence**: Metrics tracked include MAE, R², accuracy, F1, ROC-AUC. `_context_fallback()` at line 42 handles 4 failure modes.

### Bullet 3: "Async task processing with Celery + Redis pipeline"

**✅ TRUE** — Retry logic with `max_retries=3`, exponential backoff, and jitter.

**Code evidence**: `planning_task.py:335-356` — Celery task definition. `planning_task.py:359-387` — dispatcher with thread fallback.

**How to demonstrate**: "The pipeline runs asynchronously through 6 stages: load → sync → extract → learn → optimize → explain. Progress is reported after each stage via callbacks. Celery is preferred for production, with a thread-based fallback for development. Failed tasks retry up to 3 times with exponential backoff."

**Skepticism point**: "What types of errors trigger a retry?" Must distinguish between retryable (transient) and non-retryable (bug) errors.

### Bullet 4: "Async SQLAlchemy with PostgreSQL"

**✅ TRUE** — SQLAlchemy uses `create_async_engine` with `AsyncSession` for all route handlers.

**Code evidence**: `database.py:10` uses `create_async_engine()`. All route handlers use `async def` with `await db.execute(select(...))`. Sync `get_db` still exists for Celery worker.

**How to demonstrate**: "I migrated from sync to async SQLAlchemy 2.0. The database module uses `create_async_engine` with `AsyncSession`. All route handlers are `async def` with `await db.execute()`. FastAPI's async support means the event loop is never blocked by DB operations. The Celery worker still uses the sync engine since Celery tasks are synchronous."

### Bullet 5: "React + TypeScript frontend with 9 pages"

**Code evidence**: `App.tsx` — 9 routes. `pages/` — 9 page components. `hooks/usePlanStatus.ts` — polling hook.

**How to demonstrate**: "The frontend covers the full sprint planning workflow: dashboard, dataset upload, sprint configuration, async optimization with polling, plan review, explainability panel, approval, and reports."

**Skepticism point**: "9 pages" sounds like a lot but many are thin wrappers. Some pages like `NotUsed` are placeholders.

### Bullet 6: "Kanban-style sprint plan approval"

**✅ TRUE** — Drag-and-drop kanban board with 5 columns powered by `@dnd-kit/core`.

**Code evidence**: `frontend/src/pages/SprintPlanApproval.tsx` — uses `DndContext`, `SortableContext`, `DragOverlay` with 5 columns (Backlog, Selected, In Progress, Review, Approved). Items are sorted by priority within each column. Approve button syncs status changes to backend.

**How to demonstrate**: "The approval page is a kanban board with swimlanes: Backlog → Selected → In Progress → Review → Approved. Users can drag stories between columns. The final approve action persists the plan state to the backend via PUT /stories/:id for each moved story."

### Bullet 7: "Docker Compose microservice architecture"

**Code evidence**: `docker-compose.yml` (6 services: db, redis, minio, backend, celery_worker, frontend). `docker-compose.prod.yml` (adds nginx).

**How to demonstrate**: "Services include PostgreSQL 15, Redis 7, MinIO S3-compatible storage, FastAPI backend, Celery worker, and React+Vite frontend. Production adds an Nginx reverse proxy. All use named volumes for persistence."

### Bullet 8: "GitHub Actions CI with PostgreSQL + Redis"

**Code evidence**: `.github/workflows/test.yml` — runs pytest with PG 15 and Redis 7 as service containers.

**How to demonstrate**: "CI runs 27+ tests (19 unit, 8 integration) on push/PR with real PostgreSQL and Redis service containers. Tests verify the full pipeline: dataset upload → sprint creation → plan generation → approval → export."

---

## Phase 6 — Technology Deep Dives

### 6.1 PuLP / MILP Deep Dive

**Fundamentals**: PuLP is an LP/MILP modeling library. You define variables with categories (LpBinary, LpInteger, LpContinuous), add constraints with `+=`, set objective with `prob +=`, and solve.

**In this repo**: `_milp_select()` at line 289:
- 15 variables: one LpBinary per story
- 1 capacity constraint + n dependency constraints
- Objective: maximize weighted sum
- 15-second time limit, single thread

**Limitations**:
- CBC is not the fastest solver. For large problems, Gurobi or OR-Tools CP-SAT are 10-100x faster.
- No parallel solving (threads=1).
- No optimality gap tracking.
- PuLP 3.x changed API — `LpVariable.dicts` pattern differs.

**Scaling**: For >2000 stories, MILP becomes slow. Options:
- Column generation / decomposition
- Heuristic presolve to reduce problem size
- Switch to greedy-only for large backlogs
- Use OR-Tools CP-SAT (handles 10k+ variables easily)

### 6.2 Celery Deep Dive

**Architecture**:
```
Producer (FastAPI) → Broker (Redis) → Worker (celery_worker)
                            ↓
                    Result Backend (Redis)
```

**Task states**: PENDING → STARTED → PROGRESS (custom) → SUCCESS/FAILURE.

**In this repo**: Custom state "PROGRESS" via `self.update_state(state="PROGRESS", meta=...)`. This is not a valid Celery state — it's stored in result backend metadata but Celery's task state machine doesn't recognize it.

**Thread fallback** (line 359-387):
```python
thread = Thread(target=execute_planning_pipeline, args=(..., job_id), daemon=True)
```
Daemon threads are killed when the main process exits. Any running job is lost on restart.

**Production concerns**:
- No task result TTL (results persist in Redis indefinitely)
- No rate limiting
- No concurrency configuration (default is worker_prefetch_multiplier=4)
- No monitoring (no Flower)

### 6.3 JWT Authentication

**How it works**:
1. User registers → password hashed with bcrypt → stored in DB
2. User logs in → server verifies password → returns JWT token
3. Client sends JWT in `Authorization: Bearer <token>` header
4. `get_current_user()` decodes JWT, verifies signature, loads user from DB

**In this repo**: Authentication is disabled by default (`enforce_auth: false`). When disabled, all users get `_AnonymousUser` with role "scrum_master". This is a development convenience but a security risk in production.

**Key components**:
- `security.py`: `create_access_token()`, `decode_access_token()`, `verify_password()`
- `auth_backend.py`: FastAPI Users JWT strategy
- `users_fastapi.py`: FastAPI Users UserManager

**Issues**: Token has no refresh mechanism. Token expiry is 30 minutes. No token revocation (no blacklist).

### 6.4 React Query Polling

**The `usePlanStatus` hook pattern**:
```typescript
const { data } = useQuery({
  queryKey: ['planStatus', jobId],
  queryFn: () => fetch(`/api/v1/plans/status/${jobId}`),
  refetchInterval: (query) => 
    query.state.data?.status === 'complete' ? false : 2000,
})
```

**Why 2 seconds**: Balance between responsiveness and server load. Each poll hits the backend, does a DB check (or Celery state check), and returns. Under heavy load, 2-second polling from multiple clients could exhaust DB connections.

---

## Phase 7 — System Design Discussion

### If asked: "Design this system from scratch"

**Requirements**:
- Scrum master uploads backlog CSV
- System generates optimal sprint plan respecting capacity, risk, skills, dependencies
- System explains each decision
- System learns from historical data
- Web UI for entire workflow
- Asynchronous processing (optimization can take minutes)

**Architecture** (as implemented):
```
Client → Load Balancer → Nginx → /api/ → FastAPI Backend
                                   / → React SPA (served via Vite preview)
```

**API Design** (RESTful):
```
POST   /api/v1/datasets/upload
POST   /api/v1/sprints/
POST   /api/v1/stories/
POST   /api/v1/plans/generate       → returns job_id
GET    /api/v1/plans/status/{id}    → poll for completion
GET    /api/v1/plans/{id}           → get plan
GET    /api/v1/plans/{id}/explain   → get explanations
PUT    /api/v1/plans/{id}/approve   → approve
POST   /api/v1/plans/{id}/export    → CSV/JSON export
```

**Database Schema** (7 tables):
```sql
scrum_teams (team_id, name, team_size, capacity, skills)
users (id, email, hashed_password, role, is_active)
contexts (id, team_id, urgency_weight, value_weight, alignment_weight, computed_at)
dataset_uploads (upload_id, team_id, filename, file_path, row_count, is_valid, uploaded_at)
sprints (sprint_id, team_id, goal, start_date, end_date, capacity, status)
user_stories (story_id, sprint_id, title, story_points, business_value, risk_score, ...)
sprint_plans (plan_id, sprint_id, selected_stories, total_value, total_risk, capacity_used, status)
explanations (explanation_id, plan_id, story_id, is_selected, reason, confidence_score, ...)
```

**Concurrency model**:
- API: FastAPI with sync SQLAlchemy in threadpool
- Async processing: Celery worker (separate process) with thread fallback
- Frontend polling: TanStack Query 2s interval

**Scaling**:
- Backend: Horizontal scaling behind Nginx (stateless API)
- Database: PostgreSQL read replicas for queries (plans, explanations), primary for writes
- Celery: Multiple workers, increase concurrency
- Redis: Can be clustered for high availability
- Optimization: For large backlogs, pre-filter aggressively, use greedy heuristic, consider OR-Tools

**Caching**:
- Plan results are stored in PostgreSQL (not cached in Redis)
- No in-memory caching for optimization results
- Frontend caches via TanStack Query (stale-while-revalidate)

**Bottlenecks**:
1. **MILP solver** (single-threaded, 15s limit): The main bottleneck. Backlogs >2000 stories may not converge.
2. **Synchronous DB**: Threadpool saturation under high concurrent load.
3. **In-memory job state**: Lost on restart. Long-running jobs are fragile.
4. **Dataset upload**: Large CSVs are loaded entirely into memory with pandas.

**Future improvements**:
- Async SQLAlchemy 2.0
- Retry logic with dead letter queue
- Persistent job state (in DB)
- WebSocket for real-time job status
- OR-Tools CP-SAT for faster optimization
- Better caching (Redis for recent plans)
- Rate limiting and auth enforcement by default
- File size limits and streaming CSV parsing

---

## Phase 8 — Mock Interview (Simulated)

**Assume the candidate claims: "I built ApexS end-to-end — backend, frontend, ML, deployment."**

### Level 1: Resume Walkthrough

**Q**: Walk me through this project. What problem does ApexS solve?

**Expected**: Concise elevator pitch: "ApexS replaces gut-feeling sprint planning with mathematical optimization. It uses ILP to select the optimal set of stories that maximizes business value subject to capacity, risk, skill, and dependency constraints."

**Q**: What was your specific role?

**Expected**: "I built the entire system — backend services, API layer, Celery integration, frontend pages, ML weight learning, Docker deployment, CI pipeline, and academic experiments."

### Level 2: Architecture

**Q**: Explain the architecture. How do the services communicate?

**Expected**: "FastAPI serves REST endpoints. The frontend polls the status endpoint every 2 seconds. Celery with Redis handles async optimization. MinIO stores dataset files with local fallback. PostgreSQL stores everything else."

**Q**: Why did you choose Celery over other async task frameworks?

**Expected**: "Celery is mature, well-supported, and Redis is already in the stack for caching. The thread fallback provides development convenience without requiring Redis."

### Level 3: Technology Choices

**Q**: Why PuLP and not OR-Tools?

**Expected**: "PuLP is more Pythonic — you declare variables, add constraints with `+=`, and solve. OR-Tools has better performance for combinatorial problems and scales better, but PuLP was sufficient for our backlog sizes (<2000 stories). The greedy fallback handles cases where PuLP fails."

**Q**: Why async SQLAlchemy in an async FastAPI app?

**Expected** (honest): "Async SQLAlchemy 2.0 with `create_async_engine` and `AsyncSession` ensures the event loop is never blocked by database operations. This is essential for high-concurrency workloads."

### Level 4: Implementation

**Q**: Write the MILP formulation on the whiteboard.

**Expected**:
```
Variables:
  x_i ∈ {0, 1}  for each story i

Objective:
  maximize Σ (score_i * x_i)

Constraints:
  Σ (story_points_i * x_i) ≤ capacity                    [capacity]
  x_i = 0 if risk_i > threshold                          [risk, handled pre-solve]
  x_i = 0 if skill_i ∉ available_skills                  [skill, handled pre-solve]
  x_i = 0 if status_i ∈ {done, closed}                   [status, handled pre-solve]
  x_j ≤ x_i if story_j depends on story_i                 [dependency]
```

**Q**: How does the weight learning fallback chain work?

**Expected**: "train_with_metrics() checks: (1) Is dataset empty? → context weights. (2) Is sklearn importable? → context weights. (3) Are there ≥10 samples? → context weights. (4) Does target have both classes? → context weights. (5) Train LogisticRegression, sign-map coefficients to weights."

### Level 5: Edge Cases

**Q**: What happens if a sprint has 0 capacity?

**Expected**: "`solve()` at line 333 returns empty result with `solver_status='empty'`. The capacity constraint becomes `lpSum <= 0`, which forces x_i = 0 for all stories. The greedy fallback also produces empty selection."

**Q**: What happens if all stories are high-risk?

**Expected**: "The risk filter at line 238 removes all stories. `_filter_feasible_stories` returns empty list. `solve()` at line 360 returns empty result with `solver_status='no-feasible-stories'` and a warning."

**Q**: What if a story depends on itself?

**Expected**: "`depends_on = parse_depends_on(story.depends_on)` parses the dependency list. If a story has `depends_on = ['US1']` and its own id is US1, then the MILP constraint at line 306 becomes `x[US1] <= x[US1]`, which is always satisfied. No infinite loop."

### Level 6: Performance

**Q**: How would you optimize for 10,000 stories?

**Expected**: "The MILP solver with 15-second time limit won't converge. I'd: (1) Pre-filter aggressively — remove stories that clearly won't fit. (2) Use the greedy heuristic, which is O(n²) and handles 10k stories in seconds. (3) Consider OR-Tools CP-SAT, which scales better for combinatorial problems. (4) Batch process or use column generation."

**Q**: How many concurrent users can this system handle?

**Expected**: "API is stateless behind Nginx, so horizontal scaling works. DB is the bottleneck — with async SQLAlchemy, we can handle hundreds of concurrent DB operations. Celery workers can scale independently for optimization tasks."

### Level 7: Internals

**Q**: How does PuLP's CBC solver work internally?

**Expected**: "CBC (COIN-OR Branch-and-Cut) starts with LP relaxation (ignoring integer constraints). If the LP solution has fractional variables, it branches — creates two subproblems with additional bounds (x_i = 0, x_i = 1). It prunes branches that can't beat the current best integer solution. Cutting planes tighten the LP relaxation. The process continues until all branches are pruned or time limit is reached."

**Q**: How does TanStack Query's refetchInterval work for polling?

**Expected**: "`refetchInterval` is a callback or number. As a function, it receives the query state and returns the interval in ms. Returning `false` stops polling. TanStack Query uses `setInterval` internally. When the browser tab is backgrounded, the interval may be throttled by the browser."

### Level 8: Tradeoffs

**Q**: Tradeoffs of rule-based explainability vs. SHAP?

**Expected**: "Rule-based is deterministic, auditable, and requires no additional ML. Every explanation maps to a specific constraint violation or score component. SHAP provides Shapley values — mathematically fair feature attribution — but requires a model, is computationally expensive, and the results can be hard to interpret for non-technical users. For sprint planning, 'Risk exceeds threshold' is more actionable than 'SHAP value = -0.15'."

**Q**: Tradeoffs of thread fallback vs. always requiring Celery?

**Expected**: "Thread fallback simplifies development — no need to run Redis and a Celery worker locally. But daemon threads are killed on process exit, so in-flight jobs are lost. Celery provides durability (tasks survive worker restarts via broker), monitoring (Flower), and horizontal scaling (add more workers)."

### Level 9: Failure Scenarios

**Q**: What happens if Redis goes down?

**Expected**: "If `settings.use_celery = true` and Redis is down, `run_async_job` at line 363 catches the connection error and falls back to thread if `allow_thread_fallback = true`. Otherwise, it raises 503. The health endpoint at `main.py:122-126` reports Redis status."

**Q**: What happens if the Celery worker crashes mid-pipeline?

**Expected**: "The task is automatically re-queued up to 3 times with exponential backoff (5s → 10s → 20s, max 300s, with jitter). `task_acks_late=True` means the task is re-delivered if the worker crashes mid-execution. If all retries fail, the task moves to FAILURE state and the frontend displays an error."

**Q**: What if two users trigger plan generation for the same sprint simultaneously?

**Expected**: "Each call creates a separate `job_id`. The Celery worker processes both. The first to complete creates a SprintPlan. The second creates another SprintPlan. Both are valid plans for the same sprint. There's no deduplication."

### Level 10: Research-Level

**Q**: How would you prove that the ILP-optimized plan is better than the team's manual plan?

**Expected**: "We'd need a controlled experiment: have the team create a manual sprint plan, then run the ILP on the same backlog. Compare: (1) Total objective score of selected stories. (2) Capacity utilization. (3) Risk distribution. (4) Skill alignment. (5) Dependency satisfaction. The paper at `files/APEX_S_IEEE_Paper.pdf` does ablation studies showing MILP outperforms greedy, random, and fixed-weight baselines on 4 datasets."

**Q**: How would you extend this to multi-sprint planning?

**Expected**: "Multi-sprint planning requires: (1) Ordering stories across sprints (dependencies that cross sprint boundaries). (2) Capacity planning over a horizon. (3) Learning velocity per sprint for capacity estimation. (4) Maybe a rolling horizon approach: plan N sprints, execute 1, re-plan. The MILP formulation would add a sprint index variable: `x_{i,s} ∈ {0,1}` with constraints ensuring each story assigned to at most one sprint and dependencies respect sprint ordering."

---

## Phase 9 — Brutal Interview Questions

### Beginner (1-20)

1. **What is the project and what problem does it solve?**
2. **What is a REST API and how does FastAPI implement one?**
3. **What is SQLAlchemy and why use an ORM?**
4. **What is React and how does it differ from vanilla JavaScript?**
5. **What is Docker Compose and why use it?**
6. **What is pytest and how is it configured here?**
7. **What is JWT and how does authentication work?**
8. **What is Celery and what problem does it solve?**
9. **What is the difference between MILP and LP?**
10. **What is the knapsack problem and how is it related to this project?**
11. **What is scikit-learn and what model does this project use?**
12. **What is pandas and why use it for datasets?**
13. **What is Vite and how does it differ from Create React App?**
14. **What is Tailwind CSS and why use it?**
15. **What is TanStack Query and what problem does it solve?**
16. **What is Zustand and how does it differ from Redux?**
17. **What is MinIO and how does it compare to AWS S3?**
18. **What is Alembic and why use it?**
19. **What is the difference between SQLite and PostgreSQL?**
20. **What is a GitHub Actions service container?**

### Intermediate (21-50)

21. **Why does FastAPI run sync route handlers in a threadpool?**
22. **How does the MILP formulation handle dependency constraints?**
23. **What is the time complexity of the greedy fallback algorithm?**
24. **How does the weight learning sign-mapping work?**
25. **Why does context extraction compute urgency as velocity/capacity?**
26. **How does the Celery task report progress?**
27. **What happens when the MILP solver times out at 15 seconds?**
28. **How are duplicate story IDs handled across dataset uploads?**
29. **Why does `save_bytes` fall back to local storage silently?**
30. **How does the frontend handle the case where polling never completes?**
31. **What is the purpose of `check_same_thread=False` in SQLite configuration?**
32. **How does the health endpoint determine if the system is "degraded"?**
33. **Why are required skills normalized to lowercase?**
34. **How does the explainability engine determine the priority of rejection reasons?**
35. **What happens if the explainability engine receives a nonexistent story_id?**
36. **Why does the weight learning model train on ALL historical uploads?**
37. **How does `parse_depends_on` handle various input formats?**
38. **Why is Celery's result backend configured to use Redis?**
39. **How does the frontend redirect from optimization to plan view?**
40. **What is the role of `_filter_feasible_stories` in the optimization pipeline?**
41. **Why does the optimization engine track `filtered_by_risk` but not stories that pass all filters?**
42. **How does the Zustand persist middleware store state across page reloads?**
43. **What happens if a CSV upload is missing the `story_points` column?**
44. **How does `GET /plans/status/{job_id}` distinguish between Celery and thread-based jobs?**
45. **Why is the default team seeded on startup?**
46. **What is the `future=True` parameter in `create_engine`?**
47. **How does the optimizer handle stories with negative business_value?**
48. **Why does the frontend create a new TestClient for each concurrent upload?**
49. **How would you add a new baseline mode?**
50. **Why does the weight learning model use `liblinear` solver?**

### Hard (51-80)

51. **How does CBC's branch-and-cut algorithm solve the MILP?**
52. **What is the dual simplex method and how does PuLP use it?**
53. **How does the Celery retry mechanism work with exponential backoff?**
54. **What is the cost of using async SQLAlchemy vs. threadpool for DB calls?**
55. **How would you make the thread fallback job survive a server restart?**
56. **What is the optimality gap and how would you report it?**
57. **Why does the objective function use normalized scores instead of raw business value?**
58. **How would you add story-point uncertainty (e.g., PERT distributions) to the optimization?**
59. **How does the explainability engine's confidence score differ from a proper confidence interval?**
60. **What is the difference between L1 and L2 regularization in LogisticRegression?**
61. **How would you implement a randomized controlled trial to validate the optimization?**
62. **What is the cold-start problem for weight learning and how is it addressed?**
63. **How would you detect and handle concept drift in sprint completion patterns?**
64. **How does the MinIO client's aggressive timeout (1s connect, 2s read, no retries) affect reliability?**
65. **What happens to in-flight Celery tasks during a rolling deployment?**
66. **How would you benchmark the optimization engine across different backlog sizes?**
67. **What is the memory footprint of parsing a 100MB CSV with pandas?**
68. **How does the dependency constraint `x_j <= x_i` handle cyclic dependencies?**
69. **What is the theoretical worst-case runtime of the greedy algorithm?**
70. **How does the system prevent one user from approving a plan while another modifies it?**
71. **What is the difference between `LpVariable.dicts()` and individual `LpVariable()` calls?**
72. **How would you add a constraint that at least 2 high-priority stories must be selected?**
73. **What is the significance of `solver_status` appearing as "greedy-fallback:milp-status-not-optimal"?**
74. **How does the output CSV export map business_value to "Priority"?**
75. **Why does the frontend create 9 separate page components instead of a single dynamic page?**
76. **How would you implement undo for plan approval?**
77. **What is the impact of using `db.bulk_insert_mappings()` vs. individual `session.add()`?**
78. **How does the `@app.on_event("startup")` handle alembic migrations?**
79. **What is the fractional knapsack problem and why is this 0/1 knapsack?**
80. **How would you model sprint goal alignment as a constraint?**

### Expert (81-100)

81. **How does CBC's presolve reduce problem size?**
82. **What is the structure of the branch-and-bound tree for this specific problem?**
83. **How would you convert this MILP to a constraint programming formulation?**
84. **What is the duality gap and how does it relate to the objective function?**
85. **How would you implement a shadow price analysis on the capacity constraint?**
86. **How does TanStack Query's stale-while-revalidate pattern differ from simple polling?**
87. **What is the React Query garbage collection behavior for stopped queries?**
88. **How would you implement progressive loading for explanations with >5000 stories?**
89. **How does FastAPI's dependency caching work for `Depends(get_db)` across multiple dependencies?**
90. **What is the `request.state` pattern and how would you use it for request-scoped caching?**
91. **How would the system perform under a thundering-herd of 100 simultaneous plan generation requests?**
92. **What is the impact of the GIL on the thread fallback?**
93. **How would you implement circuit breaker for failed MinIO connections?**
94. **What is the database isolation level and how does it affect concurrent sprint plan creation?**
95. **How would you implement a saga pattern for the 6-stage pipeline?**
96. **How does the CBC solver's `timeLimit` parameter interact with the branch-and-cut process?**
97. **What is the difference between CBC's "feasible" and "optimal" solutions?**
98. **How would you extend the model to account for story dependencies with lag?**
99. **What is the theoretical minimum number of constraints needed for this formulation?**
100. **Prove that the greedy algorithm is a 2-approximation for the knapsack problem. Why doesn't this apply here?**

---

## Phase 10 — Knowledge Gaps and Risk Assessment

### Critical Risks

| Risk | Issue | Code Location | Impact |
|------|-------|---------------|--------|
| **"Async SQLAlchemy" (FIXED)** | Now uses `create_async_engine` + `AsyncSession` | `database.py` | Fully async, event loop safe |
| **"Retryable pipeline" (FIXED)** | Celery task has `max_retries=3`, exponential backoff | `planning_task.py` | Failures auto-retry with backoff |
| **"Kanban approval" (FIXED)** | Drag-and-drop kanban board with swimlanes | `SprintPlanApproval.tsx` | True kanban experience |
| **Resume: "SHAP-style explainability"** | Rule-based, not SHAP | `explainability_engine.py:6` | Must clarify: "rule-based score breakdown, not SHAP" |

### Architecture Weaknesses

| Weakness | Detail | Interview Risk |
|----------|--------|----------------|
| **In-memory job state** | `_JOB_STORE` lost on restart | Medium — should persist to DB |
| **No auth by default** | `enforce_auth = False` | Medium — security concern |
| **Large CSV in memory** | `pd.read_csv()` loads entire CSV | Medium — OOM risk |
| **Single-threaded solver** | CBC with `threads=1` | Medium |
| **No connection pooling config** | Default SQLAlchemy pool may be insufficient | Medium |
| **Hardcoded progress percentages** | Not proportional to actual time | Low |

### Code Smells

| Smell | Location | Issue |
|-------|----------|-------|
| `try: from sklearn... except: LogisticRegression = None` | `weight_learning.py:12` | Import guard masks missing dependencies |
| `_is_missing(value): return bool(value != value)` | `preprocessing.py:11` | NaN check via self-inequality |
| 15-second hardcoded time limit | `optimization_engine.py:312` | Should be configurable |
| No structured logging | `main.py:28-29` | Basic logging, no JSON format |
| MinIO fallback is silent | `minio_client.py:40` | No warning when falling back to local |

### Scalability Concerns

| Concern | Detail |
|---------|--------|
| **Solver doesn't scale** | MILP with 15s limit, single thread, no warm start |
| **No caching** | Every plan request hits DB. No Redis cache |
| **Frontend polling** | Every 2 seconds per user |
| **No cursor-based pagination** | Stories and explanations use offset pagination |
| **No read replicas** | All queries hit primary DB |

### Recommended Study Plan

**Before interview, candidate must**:

1. **Master PuLP** — Write the MILP formulation from scratch without looking at code.
2. **Understand async SQLAlchemy** — Know `AsyncSession`, `select()`, `scalars()`, connection pooling.
3. **Know Celery retry** — `max_retries`, `default_retry_delay`, `task_reject_on_worker_lost`.
4. **Understand CBC** — Branch-and-bound, LP relaxation, cutting planes, presolve.
5. **Practice kanban implementation** — `@dnd-kit/core` for drag-and-drop, column state management.
6. **Read one optimization textbook chapter** — Knapsack problem, NP-hardness, approximation algorithms.
7. **Prepare failure scenarios** — What breaks if Redis dies? PostgreSQL dies? MinIO dies?
8. **Know the exact code paths** — `solve()` → `_filter_feasible_stories()` → `_milp_select()` or `_greedy_selection()`.
9. **Understand async dependency injection** — `async def get_db()` with `async with SessionLocal()`.
10. **Be ready to design improvements** — WebSocket for real-time status, persistent job state, OR-Tools migration.
