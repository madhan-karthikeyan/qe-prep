# DecisionDrift — Complete Interview Preparation Guide

> For the engineer who built this project. Covers every technical question a senior interviewer at a top-tier company could ask.

---

## Phase 1: Architectural Overview

### Project Purpose
DecisionDrift is a **deterministic-first architecture decision governance tool**. It codifies team architecture decisions (documented as ADRs — Architecture Decision Records) into enforceable rules that run on every PR/commit **without requiring an LLM for the critical enforcement path**.

### Core Problem
Engineering teams document decisions in ADRs, but nobody checks them against every code change. Dependencies, imports, APIs, and file patterns drift from the original decisions over time. DecisionDrift closes that gap by turning decisions into automated rules, running them on every change, and catching violations before they ship.

### Architecture (8 Subsystems)

```
CLI (Click) ──┬── bootstrap ──→ Registry → Evidence Collection → V3 Model → ADR Suggestion
              ├── enforce    ──→ Rule Engine (5 scanners) → ReportEnvelope → Output
              ├── init       ──→ bootstrap + approve + hook + config + CI
              ├── guard      ──→ pre-commit hook management
              ├── review     ──→ Impact Analysis → Retrieval → LLM Classification
              ├── audit      ──→ Drift + Expiry + Coverage + Quality
              ├── impact     ──→ Diff parser → Symbol extraction → ImpactReport
              ├── doctor     ──→ Health diagnostics
              ├── ingest     ──→ LLM free-text → ADR
              └── adr        ──→ approve/reject/deprecate/supersede/edit/history/show/list
```

### Folder Structure

```
src/decisiondrift/
├── __init__.py
├── adr/                  # ADR loading, parsing, writing, dedup, rule generation, supersession
│   ├── parser.py         # Frontmatter extraction + JSON Schema validation
│   ├── loader.py         # Iterates ADR-*.md files
│   ├── writer.py         # Writes ADR markdown files
│   ├── rule_generator.py # DANGER: Converts ADR prohibitions → Rule objects
│   ├── supersession.py   # Resolves which ADRs are actively governing
│   ├── dedup.py          # Duplicate detection for bootstrap
│   └── id_allocator.py   # Next available ADR-NNNN
├── adr_manager/          # CLI-facing ADR lifecycle commands
├── bootstrap/            # V3 pipeline: registry, detectors, patterns, evidence, knowledge provider
│   ├── v3.py             # 1219-line core — evidence, modeling, governance, enforceability
│   ├── registry.py       # Layered YAML/HTTP/cache registry loader
│   ├── detectors.py      # Technology signature matching (60+ technologies)
│   ├── patterns.py       # Pattern-based technology detection
│   ├── structure_scan.py # File/directory structure analysis
│   ├── candidate_generator.py # Dedup + ADRSuggestion creation
│   ├── bootstrapper.py   # Top-level bootstrap orchestrator
│   ├── template_generator.py # ADR markdown rendering
│   ├── knowledge_provider.py # LLM integration for unknown techs
│   ├── synthesis.py      # LLM-based candidate synthesis
│   ├── suggester.py      # Confidence-based suggestion filtering
│   ├── cache.py          # Template caching
│   └── default_registry.yaml # 998-line bundled technology profiles
├── classification/       # LLM-based diff classification
│   ├── classifier.py     # Pairwise (ADR, symbol) LLM classification
│   ├── models.py         # ClassificationInput, ClassificationResult
│   └── prompts.py        # System + user prompt templates
├── cli.py                # 888 lines — all Click commands
├── config.py             # decisiondrift.yml loading + custom rules
├── github/               # GitHub Action integration
│   ├── action_entrypoint.py # Docker entrypoint for PR review
│   ├── client.py         # GitHub API wrapper
│   └── comment_manager.py # PR comment upsert
├── impact/               # Diff parsing + AST analysis
│   ├── diff_parser.py    # Unified diff → ChangedFile[]
│   ├── ast_python.py     # Python symbol extraction via ast module
│   ├── ast_treesitter.py # Multi-language AST via tree-sitter
│   ├── language_registry.py # 12 language definitions + extension maps
│   ├── models.py         # ChangedFile, ChangedSymbol, ImpactReport
│   ├── reference_scan.py # Search term generation from symbols
│   ├── service.py        # analyze_diff orchestrator
│   └── treesitter_queries/ # 12 per-language query files
├── ingest/               # Free-text → ADR via LLM
├── init/                 # Project initialization orchestrator
├── llm/                  # OpenAI-compatible LLM client
│   └── client.py         # complete() + complete_json() + _call()
├── models/               # All Pydantic schemas
│   └── schema.py         # DecisionRecord, ReportEnvelope, Finding, ReviewResult, ADR_SCHEMA
├── report/               # Output formatters
│   ├── formatter.py      # text/json/sarif/markdown/html
│   ├── compiler.py       # Human-readable report assembly
│   └── github_formatter.py # GitHub-specific comment formatting
├── retrieval/            # ADR retrieval backends
│   ├── backend.py        # Abstract RetrievalBackend
│   ├── keyword.py        # Weighted keyword scoring
│   ├── embedding.py      # FastEmbed + cosine similarity
│   └── models.py         # RetrievalResult
├── review/               # Semantic review pipeline
│   └── service.py        # Orchestrator: impact → retrieval → classification
├── rules/                # Deterministic rule engine
│   ├── models.py         # Rule, RuleSet, RuleType, Action, EnforcementFinding, EnforcementResult
│   ├── engine.py         # enforce(), enforce_from_adrs() — 3 modes: diff, file, repo
│   └── scanner.py        # Full-repo dependency/import scanning
└── utils/
    └── dependency_parser.py # Parsers for 13 dependency file formats
```

### External Services
- OpenAI/Groq/Ollama LLM API (optional, for `review` and `ingest` and bootstrap `--llm`)
- HTTP technology registries (optional, for shared decision catalogs)
- GitHub API (for PR comment posting, status checks, review submission)

### Key Dependencies
- **click** — CLI framework
- **pydantic** — All data models (v2)
- **tree-sitter** + **tree-sitter-languages** — Multi-language AST (optional `[ast]` extra)
- **fastembed** — Local embedding model (optional `[embeddings]` extra)
- **openai** — LLM API client
- **PyYAML** — Config + registry parsing
- **python-frontmatter** — ADR frontmatter extraction
- **jsonschema** — ADR validation
- **httpx** — HTTP registry fetching
- **python-dotenv** — Environment variable loading

### Testing
- **pytest** with unit/integration/snapshot markers
- **syrupy** for CLI output snapshots
- 363 tests passing, 26 skipped (tree-sitter optional dep)
- coverage via pytest-cov
- Test data in `tests/data/` and `tests/sample_repos/`

### CI/CD
- GitHub Actions (`ci.yml`)
- Docker-based GitHub Action (`action.yml` + `Dockerfile`)
- Published to PyPI

---

## Phase 2: Knowledge Dependency Graph

```
Level 1 (Fundamental)
├── Python: typing (generics, Literal, union), pathlib, dataclasses, re, json, os, sys
├── git: diff, log, pre-commit hooks
├── JSON/YAML: serialization, schema validation
├── Pydantic: BaseModel, field types, validation, serialization
└── Click: CLI groups, commands, arguments, options

Level 2 (Intermediate)
├── Tree-sitter: parser generation, AST queries, captures, tree navigation
├── AST: Python ast module, NodeVisitor, walk, Import, ImportFrom, Call, Attribute
├── Docker: images, Dockerfiles, entrypoints, CMD vs ENTRYPOINT
├── OpenAI API: chat completions, system prompts, response_format
├── Frontmatter: YAML frontmatter in Markdown files
├── Dotenv: environment variable loading
└── Packaging: pyproject.toml, setuptools, entry_points, optional-dependencies

Level 3 (Advanced)
├── ADR Methodology (MADR/Y-Statements): decision records, status lifecycle, supersession
├── Governance Engineering: policy-as-code, deterministic rule engines, semantic vs syntactic gates
├── Layered Registries: bundled → remote HTTP → global cache → project cache
├── Retrieval-Augmented Generation (RAG): keyword retrieval, embedding search, cosine similarity
└── SARIF v2.1.0: Static Analysis Results Interchange Format

Level 4 (Expert)
├── Evidence-Based Reasoning: evidence collection, aggregation, role inference, contradiction detection
├── Repository Modeling: role inference (API service, frontend, library, monorepo, framework)
├── Enforcement Analysis: enforceability classification (none/weak/moderate/strong)
├── Context Window Budgeting: max pairs, similarity threshold, wall clock optimization
└── CI/CD Integration Patterns: commit status, PR comments, SARIF upload, auto-review
```

---

## Phase 3: Concept Curriculum

### 3.1 Click — CLI Framework

**What it is:** A Python library for building command-line interfaces with minimal boilerplate. Uses decorators to define commands, groups, arguments, and options.

**Why it exists:** Python's `argparse` requires verbose setup. Click provides declarative command composition, automatic help generation, and nested command groups.

**How it works internally:** Click uses decorators that wrap functions, building a tree of `Command` and `Group` objects. When invoked, it parses `sys.argv`, matches against registered commands, and calls the appropriate function with keyword arguments from parsed options.

**Where it appears:** `src/decisiondrift/cli.py` — every function is a Click command decorated with `@cli.command()` or `@cli.group()` and `@click.option()`/`@click.argument()`.

**Why chosen:** Declarative, well-tested, supports nested groups (`adr` command has 9 subcommands), automatic version option, and integrates with `CliRunner` for testing.

**Alternatives:** Typer (built on Click but more modern), argparse (stdlib), rich-click (rich formatting).

**Tradeoffs:** Click's decorator-based API can make parameter lists unwieldy for commands with many options (e.g., `bootstrap` has 15+ options). Dependency on a third-party library.

**Common misconception:** That Click manages state between commands. It doesn't — each command is stateless. The `cli` group's `@click.group()` just defines a parent command.

**Interview follow-ups:**
- How would you add middleware (e.g., timing all commands)?
- What happens when you nest `@click.group()` inside another group?
- How does `CliRunner` work for testing?

### 3.2 Tree-sitter

**What it is:** An incremental parsing library that builds concrete syntax trees (CST) for source code. Unlike regex-based parsing, it produces accurate, language-aware trees for any language with a grammar.

**Why it exists:** Traditional approaches (regex, string matching) cannot reliably parse nested language constructs. Full parsers for every language are expensive to maintain. Tree-sitter provides a shared parsing infrastructure with language-specific grammars.

**How it works internally:**
1. Each language has a tree-sitter grammar (a `.js` or `.json` file defining language syntax rules).
2. The grammar is compiled into a C parser via `tree-sitter generate`.
3. At runtime, the parser processes source text byte-by-byte into a tree with start/end byte positions and (row, column) points.
4. Queries (similar to CSS selectors for AST nodes) extract specific patterns — imports, function calls, class definitions.
5. Incremental parsing: on re-parse, unchanged regions are reused.

**Where it appears:**
- `src/decisiondrift/impact/ast_treesitter.py` — imports, API calls, symbol extraction
- `src/decisiondrift/impact/treesitter_queries/` — per-language query files (12 languages)
- `src/decisiondrift/impact/language_registry.py` — maps extensions to tree-sitter grammars
- `src/decisiondrift/rules/engine.py` — delegates to tree-sitter for non-Python files
- `src/decisiondrift/rules/scanner.py` — full-repo import scanning via tree-sitter

**Why chosen:** The project needs to scan imports and API calls across 12 languages. Using tree-sitter avoids maintaining 12 separate parsers. It's optional (`pip install decisiondrift[ast]`) — Python-only scanning works without it.

**Alternatives:**
- **Regex**: fragile, breaks on edge cases
- **Language-specific parsers** (ast modules per language): maintenance burden
- **AST-grep**: pattern-based but less flexible

**Tradeoffs:** Tree-sitter grammars may not be perfectly complete for all language versions. The `tree-sitter-languages` Python package bundles pre-compiled grammars but may lag behind. Installation adds ~15MB to the package.

**Common misconception:** Tree-sitter produces an AST. It actually produces a Concrete Syntax Tree (CST) — it retains comments, whitespace, and all tokens. ASTs are typically derived from CSTs by semantic analysis.

**Interview follow-ups:**
- How does tree-sitter handle syntax errors in input?
- What's the difference between tree-sitter and a traditional parser like the one in `ast` module?
- How would you add support for a new language?

### 3.3 Architecture Decision Records (ADR)

**What it is:** A lightweight documentation practice where each significant architecture decision is captured in a short, structured document — typically a markdown file with YAML frontmatter.

**Why it exists:** Architecture knowledge is tacit. When team members leave, or decisions age, nobody remembers why something was done a certain way. ADRs create a persistent, reviewable trail of decisions.

**How it works internally in this project:**
- ADRs are Markdown files in `docs/adr/ADR-NNNN.md`
- YAML frontmatter contains: `id`, `title`, `status`, `severity`, `source`, `prohibitions`, `keywords`, `rationale`, `evidence`, `date`, `confidence`, `owner`, `review_after`, `expires_after`, `depends_on`, `superseded_by`, etc.
- Status lifecycle: proposed → accepted/rejected → deprecated/superseded
- Only `accepted` ADRs generate enforcement rules
- Parsing: `src/decisiondrift/adr/parser.py` uses `python-frontmatter` to extract metadata, then validates against JSON Schema (`ADR_SCHEMA` in `schema.py`)
- Loading: `src/decisiondrift/adr/loader.py` iterates `ADR-*.md` files
- Supersession resolution: `src/decisiondrift/adr/supersession.py` filters out superseded and dependency-invalidated ADRs

**Why chosen:** ADRs are widely adopted in the industry, language-agnostic, work with git, and are human-readable without tooling.

**Alternatives:** RFC documents, Google Docs, Notion/Confluence (not code-coupled), architectural decision capture tools.

**Tradeoffs:** ADRs require discipline to maintain. They can become stale if not reviewed. The frontmatter format can be verbose.

**Common misconception:** ADRs must document only "big" decisions. In practice, small decisions (use this library, avoid that pattern) are equally valuable.

**Interview follow-ups:**
- How would you handle ADR versioning when multiple ADRs conflict?
- What happens when someone manually edits an ADR file?
- How does the supersession chain work (A → B → C)?

### 3.4 Deterministic Rule Engine

**What it is:** A rule matching system that evaluates code against predefined patterns without any statistical/ML inference. Five rule types: dependency, import, API, path, config.

**Why it exists:** The core value proposition of DecisionDrift — enforcement without LLM cost or unreliability. Deterministic rules are fast, predictable, auditable, and have zero false positives (if the patterns are correct).

**How it works internally:**
1. Rules are generated from ADR prohibitions (in `src/decisiondrift/adr/rule_generator.py`):
   - Each prohibition → dependency rule + import rule (both `BLOCK` action)
   - Confidence defaults: manual=HIGH, bootstrap=MEDIUM, ingest=LOW
2. Custom rules can be added in `decisiondrift.yml` without an ADR
3. Enforcement flows (3 paths):
   - `_enforce_diff()`: Diff-based, only scans changed files
   - `_enforce_file()`: Single file mode (editor integration)
   - `_enforce_repo()`: Full repo scan (audit mode)
4. Each scanner extracts matches and hands them to the rule matcher
5. Action downgrade logic: if confidence < 0.50 → INFO; if < 0.80 and action is BLOCK → WARN

**Why chosen:** Reliability and speed. The `enforce` command is the primary CI gate — it must complete in seconds with zero false positives.

**Alternatives:** LLM-based enforcement (expensive, non-deterministic), Open Policy Agent (infra-focused, not code-aware), custom scripts per language.

**Tradeoffs:** Simple substring matching can produce false positives (e.g., prohibiting "flask" would also catch "flask-cors"). The match logic handles some cases but not all.

**Common misconception:** That the rule engine does deep semantic analysis. It does simple string/substring matching. A dependency named "fastapi" is detected by checking if the rule's `match` string appears in the dependency name.

**Interview follow-ups:**
- How would you implement a rule that only applies to specific directories?
- How do you avoid false positives with substring matching?
- What happens when a dependency file format isn't supported?

### 3.5 Layered Technology Registry

**What it is:** A hierarchical configuration system where technology definitions come from multiple sources in order of precedence.

**How it works:** `src/decisiondrift/bootstrap/registry.py:171-225`
1. Bundled `default_registry.yaml` (~100 technologies, 998 lines)
2. Remote HTTP registries (via `--registry-url` or config)
3. Global cache (`~/.config/decisiondrift/cache.yaml`)
4. Project cache (`.decisiondrift/cache.yaml`)

Each layer's `technologies`, `governance_templates`, `file_evidence`, `dir_evidence`, `language_evidence`, etc. are merged. Later layers override earlier ones.

**Why chosen:** Allows organizations to extend the registry with internal technologies without modifying the package. Enables sharing via HTTP.

**Tradeoffs:** No conflict detection. If two layers define different categories for the same technology, the last one wins silently. No version pinning for HTTP registries.

### 3.6 Bootstrap V3 Pipeline

**What it is:** The most complex subsystem (1219 lines in `v3.py`). Automatically discovers architecture decisions from repository structure.

**Pipeline stages:**
1. **Evidence Collection** (`collect_evidence`): Scans dependency files (requirements.txt, pyproject.toml, package.json, go.mod, Cargo.toml), import statements (Python AST + JS/Go regex), file names, directory names, language files
2. **Technology Candidate Building** (`build_technology_candidates`): Groups evidence by technology name using signature matching from registry. Calls LLM for unknown dependencies (optional).
3. **Repository Role Inference** (`infer_repository_role`): Heuristic-based — counts runtime frameworks, checks for package manifests, detects monorepo via scope paths
4. **Context Application** (`apply_repository_context`): Adjusts technology roles based on repo type (library repos don't produce framework ADRs)
5. **Governance Discovery** (`discover_governance_candidates`): Matches technologies to governance templates, runs enforceability analysis
6. **Suggestion Generation** (`generate_v3_suggestions`): Filters by confidence, deduplicates, creates ADR records with rules

**Key design decisions:**
- Repositories typed as "framework" or "library" suppress most governance candidates (they build the decision, they don't consume it)
- Evidence level aggregation requires 2+ runtime STRONG pieces to reach STRONG — conservative by design
- `_governance_suppression_reason` prevents data-access ADRs in library repos

### 3.7 LLM Integration

**What it does:** Three optional LLM paths:
1. **Bootstrap `--llm`**: Recognizes unknown technologies via `KnowledgeProvider`
2. **Review command**: Classifies code changes against ADRs via `Classifier`
3. **Ingest command**: Converts free-text notes into ADR candidates

**LLM Client** (`src/decisiondrift/llm/client.py`):
- Wraps OpenAI-compatible API
- `complete()` for text, `complete_json()` with `response_format={"type": "json_object"}`
- Supports custom base URLs (Ollama, Groq, OpenRouter)
- Silent fallback if no API key is configured

**Classification flow** (`src/decisiondrift/classification/classifier.py`):
- For each (ADR, symbol) pair (max 15), sends a prompt with the ADR title, rationale, exceptions, symbol name, file path, and diff hunk
- LLM returns JSON: classification (5 levels), evidence_strength (3 levels), reasoning, suggested_action
- On failure → `needs_human_review`

**Why chosen:** LLM provides semantic understanding that deterministic rules cannot (e.g., "this function call violates the intent of the ADR even though the specific library isn't prohibited").

**Tradeoffs:** Cost ($0.01-0.10 per pair), latency (2-5s per pair), non-determinism (same input may produce different results).

---

## Phase 4: Repository Deep Dive

### `src/decisiondrift/rules/engine.py` (484 lines)

**Purpose:** Core enforcement engine — the heart of the deterministic value prop.

**3 enforcement modes:**

| Mode | Entry Point | When Used | Files Scanned |
|------|-------------|-----------|---------------|
| Diff | `_enforce_diff()` | `enforce --from-git`, pre-commit hook | Changed files only |
| File | `_enforce_file()` | Editor integration (`--file`) | Single file |
| Repo | `_enforce_repo()` | `audit`, `enforce .` | All files |

**Key function: `enforce_from_adrs()`**
```python
def enforce_from_adrs(adrs, repo_path, diff_text, custom_rules, file_path):
    all_rules = []
    for adr in adrs:
        all_rules.extend(_rules_for_adr(adr))  # prohibitions → dependency + import rules
    if custom_rules:
        all_rules.extend(custom_rules.rules)
    return enforce(RuleSet(rules=all_rules), ...)
```

**Key function: `_enforce_diff()`**
1. Parse diff via `parse_diff()` → `list[ChangedFile]`
2. For each file type, check corresponding rules:
   - Dependency files (requirements.txt, etc.) → `_extract_deps_from_file()` → match against dependency rules
   - Source files → `_scan_imports_in_diff()` → match against import rules
   - All changed files → regex match against path rules
   - Source files → `_scan_api_calls()` → match against API rules
   - Config files → `_scan_config_pattern()` → match against config rules

**Key function: `_to_finding()`**
- Converts a `Rule` match into an `EnforcementFinding`
- Confidence downgrade: <0.50 → INFO, <0.80 + BLOCK/REQUIRE_APPROVAL → WARN
- Severity map: BLOCK=critical, REQUIRE_APPROVAL=high, WARN=medium, INFO=low

**Design decisions:**
- Separate scan+match loops for each rule type (not a single visitor pattern)
- Confidence-based action downgrade prevents bootstrap-generated rules from blocking CI
- Diff mode skips non-changed files entirely → fast for small PRs

**Complexity:** O(R × F) where R = rules per type, F = files. Each file is parsed once per rule type.

**Possible improvements:**
- Batch processing of files per rule type (currently re-parses same file across multiple rules)
- Caching parsed dependency trees across invocations
- Parallel file scanning

**Interview discussion points:**
- Why does `_scan_imports_in_diff` return `list[tuple[str, str]]` (import, file) instead of a structured object?
- What happens when a dependency file is modified but doesn't add the prohibited dependency? (It's still scanned — every change to requirements.txt triggers a full scan of all deps.)
- How would you implement a rule that checks for the *absence* of something (e.g., "must use Redis")?

### `src/decisiondrift/bootstrap/v3.py` (1219 lines)

**Purpose:** V3 bootstrap pipeline — the largest and most complex module.

**Pipeline stages (detailed):**

#### Stage 1: Evidence Collection (`collect_evidence`)
Three sub-functions:
- `_collect_dependency_evidence()` — scans requirements.txt, pyproject.toml, package.json, go.mod, Cargo.toml. Extracts dependency names and assigns roles (runtime/dev/test).
- `_collect_import_evidence()` — scans Python AST for imports, JS/TS for import statements and `require()`, Go for `import "..."`. Also detects FastAPI app entrypoints.
- `_collect_file_evidence()` — uses registry `file_evidence` (e.g., Dockerfile → Docker), `dir_evidence` (e.g., `routers/` → FastAPI), `language_evidence` (e.g., `*.py` → Python).

Each evidence item gets:
- `kind`: dependency/import/file/directory/entrypoint/language
- `value`: the matching string
- `source_path`: relative path
- `role`: runtime/dev/test/example/tooling/optional/unknown (inferred from path)
- `level`: strong/moderate/weak (based on role and evidence type)
- `scope_path`: first directory segment (for monorepo support)

#### Stage 2: Technology Candidate Building (`build_technology_candidates`)
- Groups evidence by technology name using `_technology_for_evidence()`
- Compares evidence values against `TECH_SIGNATURES` and registry `evidence_resolution`
- If LLM is available, unrecognized runtime dependencies are sent to `KnowledgeProvider.recognize_technology()`
- For each technology, computes:
  - **Category**: framework/database/orm/cache/queue/frontend/language/testing/etc.
  - **Evidence level**: WEAK/MODERATE/STRONG via `_aggregate_evidence_level()`
  - **Role**: primary/supporting/incidental/dev/test/example/tooling via `_candidate_role()`
  - **Contradictions**: mismatches between evidence and expected patterns
  - **Suppression reason**: why this technology should not generate an ADR

#### Stage 3: Repository Role Inference (`infer_repository_role`)
Heuristic chain:
1. Check `SELF_FRAMEWORK_REPOS` (repo name matches known framework)
2. Check for runtime frameworks (FastAPI/Flask/Django/Express/Gin → api_service)
3. Check for frontend + package.json without runtime framework → frontend_app
4. Check for build manifests without runtime framework → library
5. Multiple scope paths with multiple frameworks → monorepo
6. Default → unknown

#### Stage 4: Context Application (`apply_repository_context`)
Adjusts technology roles based on repository context:
- Library repos → framework evidence is "supporting" not "primary"
- API service repos → database is "primary", ORM is "supporting", container is "tooling"
- Framework repos → tech matching the repo product is suppressed (it's the output, not a decision)

#### Stage 5: Governance Discovery (`discover_governance_candidates`)
For each non-suppressed technology with primary/supporting role:
1. Look up governance template in registry (title, prohibitions, rationale)
2. If no template and LLM available, generate via `KnowledgeProvider.generate_template()`
3. Skip if enforceability analysis fails
4. Run `_governance_suppression_reason()` for repo-specific filtering

#### Stage 6: Enforceability Analysis (`analyze_enforceability`)
- Checks if the candidate has a template with prohibitions
- Tooling/supporting-only technologies → not enforceable
- Creates `RuleCandidate` objects: each prohibition → dependency match + import match (both `block`)
- Special case: SQLAlchemy → also warns on `sqlite3` imports (violation of ORM usage)
- Weak evidence → not enforceable (but rules are still generated for reference)

#### Stage 7: Suggestion Generation (`generate_v3_suggestions`)
- Filters by minimum confidence level
- Deduplicates against existing ADR titles (Jaccard similarity on title + keywords)
- Creates `DecisionRecord` with generated rationale, evidence, prohibitions
- Creates `ADRSuggestion` wrapping tech + ADR + rules

**Key design decisions:**
- Exclude dirs are hardcoded (no config) — `.git`, `node_modules`, `__pycache__`, etc.
- Role inference from path segments is heuristic and limited (first segment only)
- Evidence level uses a conservative aggregation (needs 2+ runtime STRONG items for STRONG)
- Suppression is aggressive — many valid candidates may be suppressed

**Complexity:** O(F + I + D + T × E) where F = files, I = imports, D = dependencies, T = technologies, E = evidence items per technology.

**Possible improvements:**
- Make exclude dirs configurable
- Add multi-segment scope paths for deep monorepos
- Cache AST parsing for import evidence collection
- Make registry signatures editable without modifying the bundled YAML

**Interview discussion points:**
- The `_role_from_path()` function uses `set & set` operations for efficiency — why is this correct?
- `Scope_path` is only the first directory segment — what happens with deeply nested services?
- Evidence levels rely on role inference which relies on path patterns — circular dependency?
- How would you support evidence from `docker-compose.yml` or `Kubernetes manifests`?

### `src/decisiondrift/impact/ast_treesitter.py` (208 lines)

**Purpose:** Bridge between tree-sitter and DecisionDrift's impact analysis.

**3 extraction functions:**

1. `extract_symbols_treesitter()`:
   - Loads language grammar + parser via `tree_sitter_languages`
   - Runs `symbols` query from `treesitter_queries/`
   - Returns `ChangedSymbol` with name, type (function/method/class), file, line range

2. `extract_imports_treesitter()`:
   - Runs `imports` query
   - Language-specific segment extraction:
     - JS/TS: first path segment (`express` from `'express'`, not `express/lib/router`)
     - Go: first path segment after stripping quotes
     - Java: first dotted segment
     - Rust: first `::` segment
     - Others: full import string
   - Filters relative imports in JS/TS

3. `extract_api_calls_treesitter()`:
   - Runs `api_calls` query
   - Returns raw call text

**Design decisions:**
- Lazy loading via `HAS_TREESITTER` global — checked before every function
- Per-language import segmentation is a long if/elif chain — not extensible without modifying the function
- Empty results on any error (file read failure, parse failure, query failure) — silent

### `src/decisiondrift/adr/rule_generator.py` (58 lines)

**Purpose:** Converts ADR prohibitions → Rule objects.

**Key pattern:** Each prohibition generates TWO rules:
```python
Rule(type=DEPENDENCY, match=prohibition, action=BLOCK)  # Catch in requirements.txt
Rule(type=IMPORT, match=prohibition, action=BLOCK)       # Catch the import itself
```

**Why both?** Defense in depth:
- Dependency rule catches cases where the dependency is declared in a manifest file
- Import rule catches cases where the code directly imports the module (even if someone vendors it without declaring it)

**Confidence assignment:**
- manual source → HIGH (0.9)
- bootstrap source → MEDIUM (0.6)
- ingest source → LOW (0.3)

**Interview question:** What if a prohibition is a file path pattern? (Nothing — both rules are always dependency and import. Path rules would need to be manually added in decisiondrift.yml.)

### `src/decisiondrift/review/service.py` (140 lines)

**Purpose:** Orchestrates the LLM-powered semantic review pipeline.

**Flow:**
```
diff_text → parse_diff → ChangedFile[]
         → analyze_diff → Symbol[]
         → generate_search_terms(symbols) → terms[]
         → keyword_backend.query(terms, ADRs) → scored ADRs
         → (if no results above threshold) embedding_backend.query(terms, ADRs)
         → for each (ADR, symbol) pair → ClassificationInput
         → classifier.classify() → Findings
         → ReviewResult
```

**Key design decisions:**
- **Two-tier retrieval**: keyword first (fast, zero deps), embedding fallback (slower, needs `fastembed`)
- **Budgeting**: `max_pairs_per_pr` (default 15) prevents unbounded LLM costs
- **Threshold**: `similarity_threshold` (default 0.5) — pairs below this don't get classified
- **Graceful degradation**: No LLM → no LLM findings, only rule engine findings

**Interview discussion points:**
- Keyword scoring is additive and normalized — does this favor ADRs with more text?
- The `_extract_hunks()` function groups diff lines by file — but uses `+++ b/` lines which are reliable?
- Why are LLM findings structured differently from rule engine findings? (They use `Finding` vs `EnforcementFinding` models.)
- How would you handle a PR with 50 changed files and 300 symbols?

### `src/decisiondrift/github/action_entrypoint.py` (233 lines)

**Purpose:** Docker entrypoint for the GitHub Action.

**Flow:**
1. Read `GITHUB_EVENT_PATH` → parse event JSON
2. Extract PR details (owner, repo, number, head SHA)
3. Fetch PR diff via `GitHubClient`
4. Load ADRs, resolve active ones
5. Run deterministic enforcement (`enforce_from_adrs`)
6. If LLM key present, run semantic review (`run_review`)
7. Combine findings, post PR comment
8. Set commit status (success/failure)
9. Optionally generate SARIF output file
10. Optionally submit formal PR review (approve/request_changes)

**Key design decisions:**
- Runs as Docker (`action.yml` specifies `using: docker` and `Dockerfile`)
- Environment variable mapping: `INPUT_*` → `config` dict
- Supports 3 review modes: comment (default), request-changes, auto
- Always sets commit status even without findings
- SARIF output is optional but supports GitHub code scanning upload

**Interview discussion points:**
- Why Docker instead of a JavaScript action? (Dependency on Python + tree-sitter native libs)
- The `GITHUB_EVENT_PATH` reading assumes the event is a pull_request — what about push events?
- Environment variable handling has no validation — what happens with garbage input?

### `src/decisiondrift/config.py` (88 lines)

**Purpose:** Loads `decisiondrift.yml` config + environment variable overrides.

**Key functions:**
- `find_config()`: Searches for `decisiondrift.yml` or `.decisiondrift.yml` in current directory
- `load_config()`: Loads YAML, sets defaults, merges env vars for LLM settings
- `load_custom_rules()`: Parses `rules:` section into `RuleSet` objects

**Config merge order (for LLM):**
1. Config file value
2. Environment variable (`DECISIONDRIFT_LLM_API_KEY`, `DECISIONDRIFT_LLM_MODEL`, `DECISIONDRIFT_LLM_BASE_URL`)
3. Default values

**Custom rules format:**
```yaml
rules:
  - match: deprecated-library
    type: dependency
    action: block
    description: "Block deprecated library"
```

**Interview question:** How would you add support for array-type environment variables for `--registry-url`? (Currently only config file supports it.)

### `src/decisiondrift/report/formatter.py` (251 lines)

**Purpose:** Unified output formatting for all commands.

**5 formats:**

| Format | Use Case | Key Details |
|--------|----------|-------------|
| `text` | CLI output | Human-readable, colored action prefixes |
| `json` | CI/tooling | `ReportEnvelope.model_dump_json(indent=2)` |
| `sarif` | GitHub code scanning | SARIF v2.1.0, maps actions to error/warning/note |
| `markdown` | PR comments | Structured with summary + findings |
| `html` | CI artifacts | Self-contained HTML with inline CSS |

**Design decision:** Single `ReportEnvelope` schema enables consistent tooling across all commands.

**SARIF limitations:** No line numbers (`region` in `physicalLocation`). Only `artifactLocation.uri` is set. This reduces usability in code scanning annotations.

**HTML output:** Inline CSS, responsive table, color-coded badges per action level.

---

## Phase 5: Resume Bullet Justification

| Resume Claim | Supporting Code | How to Demonstrate | Skepticism Points |
|---|---|---|---|
| "Built a deterministic rule engine that enforces architecture decisions across 12 languages" | `rules/engine.py`, `rules/scanner.py`, `impact/ast_treesitter.py`, `impact/language_registry.py` | Walk through: ADR → rule generation → diff parsing → per-language scanning → finding generation. Show `_enforce_diff()` → `_scan_imports_in_diff()` → tree-sitter dispatch. | "12 languages" requires tree-sitter installed. Without `[ast]` extra, only Python works. Per-language import parsing is simplistic (`startswith` / `split("/")[0]`). Silent failure on parse errors. |
| "Implemented a V3 bootstrap pipeline that automatically discovers technology decisions from repository structure" | `bootstrap/v3.py` (1219 lines) | Walk evidence collection → technology candidates → role inference → governance discovery → enforceability analysis → suggestion generation. Show how `_aggregate_evidence_level()` works. | Heuristics tuned for popular ecosystems. Novel stacks produce poor results. Suppression logic is complex and may miss valid candidates. `scope_path` is only first directory segment. |
| "Designed a layered technology registry system with YAML, HTTP, and cache layers" | `bootstrap/registry.py` | Show registry merging: bundled → HTTP → global cache → project cache. Walk through `load_registry()` with the layered merge. | No conflict detection. No HTTP caching headers/ETag support. Error handling for HTTP registries is minimal (prints warning, returns empty dict). |
| "Built a GitHub Actions integration with SARIF output and automatic PR reviews" | `github/action_entrypoint.py`, `action.yml`, `report/formatter.py` | Walk Docker entrypoint → PR diff fetch → enforce + LLM → comment posting → status setting → SARIF generation → review submission. | Docker adds ~1GB pull. Environment variable handling uses `os.environ.get()` without input validation. SARIF has no line numbers. |
| "Created an AST-based multi-language impact analysis system using Python AST and Tree-sitter" | `impact/ast_python.py`, `impact/ast_treesitter.py`, `impact/treesitter_queries/*` | Show symbol extraction for Python (`ast.walk()`) and JS/Go/Rust/etc (tree-sitter queries). Explain per-language import segmentation. | Tree-sitter queries are basic — may miss destructured imports, dynamic imports, re-exports. No incremental parsing in CI context. |
| "Developed a semantic review pipeline combining keyword and embedding retrieval with LLM classification" | `review/service.py`, `retrieval/keyword.py`, `retrieval/embedding.py`, `classification/classifier.py` | Walk through: impact → search terms → keyword → embedding fallback → pairing → LLM classification. Show the budget mechanism. | Embedding model is single fixed default. LLM classification is pairwise (slow). No batching. Keyword scoring favors verbose ADRs. |
| "Implemented complete ADR lifecycle management with 9 CLI subcommands" | `adr/` modules, `adr_manager/commands.py` | Demo: list, show, approve, reject, deprecate, archive, supersede, edit, history. Show supersession resolution. | ADR editing opens $EDITOR — no validation on save. History relies on `git log` for a specific file. No conflict detection between ADRs. |

### High-Risk Claims

| Claim | Risk Level | Why |
|---|---|---|
| "Multi-language support across 12 languages" | **HIGH** | Tree-sitter parsing is best-effort and may fail silently. Queries are basic. Without `[ast]` extra, only Python works. |
| "Deterministic enforcement" | **MEDIUM** | True for the core `enforce` path, but LLM-based `review` command is not deterministic. Reviewers may conflate the two. |
| "Zero LLM cost for the critical path" | **MEDIUM** | True for `enforce` and `bootstrap` (without `--llm`), but `review` and `ingest` require LLM. |
| "363 tests passing" | **LOW** | 26 tests are skipped (tree-sitter optional dep). Some tests may test trivial getter/setter logic. |
| "Automatically discovers decisions" | **MEDIUM** | Discovery is heuristic. Many valid decisions won't be discovered. Many discovered candidates will be suppressed. |

---

## Phase 6: Technology Deep Dives

### Python AST Module

**Fundamentals:** The `ast` module parses Python source into a tree of AST nodes. `ast.walk()` does depth-first traversal. Key node types used in this project:
- `ast.Import` — `import foo`
- `ast.ImportFrom` — `from foo import bar`
- `ast.Call` — function calls
- `ast.Attribute` — attribute access (for method calls like `obj.method()`)
- `ast.Name` — bare names (for simple function calls like `print()`)

**Usage in project:**
- `rules/engine.py:406-421`: Scans `.py` files for `ast.Import` and `ast.ImportFrom` nodes. Extracts top-level module name (`alias.name.split(".")[0]`).
- `rules/engine.py:442-456`: Scans `ast.Call` nodes. For method calls (`ast.Attribute`), reconstructs the dotted call path (e.g., `client.users.create`). For bare calls, extracts the name.
- `bootstrap/v3.py:664-703`: Import evidence collection via AST walk. Also detects FastAPI app entrypoints by looking for `FastAPI()` constructor calls.
- `impact/ast_python.py`: Symbol extraction via AST.

**Limitation:** `ast.walk()` visits every node in the tree, including nested scopes (function bodies, class bodies, etc.). This means imports inside `if` blocks or functions are reported the same as top-level imports — they're still "in the code" even if they may never execute.

**Improvement:** Use `ast.NodeVisitor` with explicit `visit_*()` methods for better control over traversal, rather than `walk()` with isinstance checks.

**Interview follow-ups:**
- How would you handle `try/except ImportError` patterns? (The import is still in the AST — should it be reported?)
- How does `ast.parse()` handle Python 3.12+ features like PEP 695 type parameter syntax?

### Tree-sitter Queries

**What they are:** Query strings in tree-sitter's S-expression-based query language. Similar to CSS selectors for AST nodes.

**Example from `treesitter_queries/javascript.py`:**
```python
QUERIES = {
    "imports": """
        (import_statement
          source: (string (string_fragment) @import))
        (call_expression
          function: (identifier) @require (#eq? @require "require")
          arguments: (arguments (string (string_fragment) @require)))
    """,
    "api_calls": """
        (call_expression
          function: (member_expression
            property: (property_identifier) @call))
    """,
    "symbols": """
        (function_declaration name: (identifier) @function)
        (method_definition name: (property_identifier) @method)
        (class_declaration name: (identifier) @class)
    """,
}
```

**What catches happen:** `query.captures(root_node)` returns `list[tuple[Node, str]]` where the string is the capture name (e.g., `@import`, `@call`).

**Per-language query files:**
- `c.py`, `cpp.py`, `csharp.py`, `go.py`, `java.py`, `javascript.py`, `kotlin.py`, `php.py`, `ruby.py`, `rust.py`, `swift.py`
- Each has exports for `imports`, `api_calls`, `symbols` queries

**Limitation:** Queries are hand-written and may miss language-specific patterns (e.g., Rust `use crate::...` vs `use std::...`, Python dynamic imports, Go `import _ "package"`).

### FastEmbed + Embedding Retrieval

**Fundamentals:** FastEmbed loads ONNX-optimized embedding models locally (no API calls). Uses `BAAI/bge-small-en-v1.5` by default (384-dimensional embeddings).

**Implementation details:**
- `retrieval/embedding.py` — lazy model loading (`_model = None` initially)
- Embeddings are cached per ADR ID in `self._adr_embeddings: dict[str, np.ndarray]`
- Cosine similarity: `dot(a,b) / (norm(a) * norm(b))`
- Query text is concatenation of search terms with spaces
- ADR text is concatenation of title + rationale + keywords

**Scaling:** Each embedding is ~384 floats = ~1.5KB. With 1000 ADRs that's ~1.5MB in memory. Well within acceptable range.

**Performance:** ONNX inference on CPU is ~10-50ms per query depending on model size. `bge-small-en-v1.5` is one of the fastest options.

**Interview point:** The embedding backend is a "fallback" behind keyword retrieval. Why?
- Keyword is faster (zero model inference, just string matching)
- Keyword works without the `[embeddings]` dependency
- Keyword is deterministic and auditable
- Embedding adds recall for semantically similar but lexically different terms

**Limitation:** `bge-small-en-v1.5` is English-only. For multilingual ADRs, a different model would be needed. Model is hardcoded in `config.py` default.

### SARIF v2.1.0

**What it is:** Static Analysis Results Interchange Format — a JSON-based standard for static analysis tool output.

**Current implementation** (`report/formatter.py:24-75`):
```json
{
  "$schema": "https://raw.githubusercontent.com/oasis-tcs/...",
  "version": "2.1.0",
  "runs": [{
    "tool": { "driver": { "name": "DecisionDrift", ... } },
    "results": [
      {
        "ruleId": "ADR-0001-dep-fastapi",
        "level": "error",
        "message": { "text": "..." },
        "locations": [{
          "physicalLocation": {
            "artifactLocation": { "uri": "requirements.txt" }
          }
        }]
      }
    ]
  }]
}
```

**What's missing:** Line numbers (`region.startLine`, `region.startColumn`). The current output only sets `artifactLocation.uri`, which means GitHub code scanning annotations won't highlight specific lines.

**Action → Level mapping:**
- `block` → `"error"`
- `require_approval` → `"warning"`
- `warn` → `"warning"`
- `info` → `"note"`

---

## Phase 7: System Design Discussion

**"Design DecisionDrift from scratch"**

### Requirements

**Functional:**
1. Scan a repository and detect technologies/frameworks in use
2. Generate architecture decision record (ADR) candidates from detected technologies
3. Allow users to approve/reject/modify ADRs
4. Enforce approved ADRs against code changes (deterministically, without LLM)
5. Support multiple programming languages for import and API analysis
6. Manage ADR lifecycle (accept, reject, deprecate, supersede)
7. Output results in multiple formats (text, JSON, SARIF, HTML)
8. Integrate with GitHub Actions via PR comments and commit status

**Non-functional:**
1. Enforcement must be deterministic (same input → same output)
2. Enforcement must complete in <5s for typical PRs
3. No LLM required for the critical enforcement path
4. Must handle 500+ ADRs without degradation
5. Must support 12+ programming languages
6. Bootstrap should complete in <60s for typical repos

### Architecture (as-built)

**Client-side tool, not a server.** Key implication: no database (filesystem-based ADR storage), no centralized auth, no hosting.

**Why not server-based?**
- Governance enforcement must work in CI where there's no server access
- ADRs are part of the codebase — they should be versioned alongside code
- A server adds latency, cost, and a failure domain for CI gates
- ADRs are semi-private team knowledge — not suitable for a shared server

**Tradeoff accepted:** No collaborative ADR editing, no centralized ADR database, no real-time enforcement dashboards, no cross-repo analytics.

### Data Flow Diagram

```
Repository (filesystem)
    │
    ▼
[bootstrap]
    │
    ├── collect_evidence() → Evidence[]
    ├── build_technology_candidates() → TechnologyCandidate[]
    ├── infer_repository_role() → "api_service"
    ├── discover_governance_candidates() → GovernanceDecisionCandidate[]
    └── generate_v3_suggestions() → ADRSuggestion[]
    │
    ▼
ADR Files (docs/adr/ADR-NNNN.md)
    │
    ▼
[enforce]
    │
    ├── load_adrs() → DecisionRecord[]
    ├── resolve_active() → DecisionRecord[]
    ├── rule_generator._rules_for_adr() → Rule[]
    ├── enforce() → EnforcementResult
    │   ├── Dependencies: match_dependency_rules()
    │   ├── Imports: scan_imports() + match_import_rules()
    │   ├── APIs: _scan_api_calls()
    │   ├── Paths: regex match
    │   └── Configs: key-value scan
    └── format_output() → text/json/sarif/html
```

### API Design (Internal interfaces)

```python
# Bootstrap
def collect_evidence(repo: Path) -> list[Evidence]: ...
def build_technology_candidates(repo: Path, evidence: list[Evidence]) -> list[TechnologyCandidate]: ...
def infer_repository_role(repo: Path, evidence: list[Evidence], technologies: list[TechnologyCandidate]) -> str: ...
def discover_governance_candidates(model: RepositoryModel) -> list[GovernanceDecisionCandidate]: ...
def generate_v3_suggestions(model, existing_titles, next_id) -> list[ADRSuggestion]: ...

# ADR
def load_adrs(adr_dir: str, status_filter: set[str] | None = None) -> list[DecisionRecord]: ...
def resolve_active(records: list[DecisionRecord]) -> list[DecisionRecord]: ...
def parse_adr_file(path: Path) -> DecisionRecord | None: ...
def write_adr(path: Path, metadata: dict, body: str) -> None: ...

# Rules
def enforce(rules: RuleSet, repo_path, diff_text=None) -> EnforcementResult: ...
def enforce_from_adrs(adrs, repo_path, diff_text, custom_rules, file_path) -> EnforcementResult: ...

# Impact
def analyze_diff(diff_text: str, repo_path: str | Path = ".") -> ImpactReport: ...
def parse_diff(diff_text: str) -> list[ChangedFile]: ...

# Retrieval
class RetrievalBackend(ABC):
    def query(search_terms, decisions, top_k=5) -> list[RetrievalResult]: ...

# Classification
class Classifier:
    def classify(inputs: list[ClassificationInput]) -> list[ClassificationResult]: ...

# Review
def run_review(diff_text: str, repo_path, adr_dir, config) -> ReviewResult: ...
```

### Database (Filesystem)

Three main stores:
1. `docs/adr/` — ADR markdown files (the "database")
2. `decisiondrift.yml` — configuration + custom rules (near the root)
3. `.decisiondrift/cache.yaml` — project-level registry cache

No actual database is used. All data is file-based. This is by design.

### Concurrency

**Current state:** Everything is single-threaded.

**How to improve:**
- Enforcement file scanning could be parallelized per file group (split by extension)
- LLM classification could batch multiple (ADR, symbol) pairs into one call
- Diff parsing could be done incrementally instead of re-parsing every time

### Scaling

**With 500+ ADRs:**
- ADR loading: O(n) file reads, ~500ms for 500 files
- Rule generation: O(p) where p = total prohibitions (p = ~1000 for 500 ADRs with 2 each)
- Enforcement: O(R × F) — worst case 1000 rules × 1000 files = 1M checks (still <1s)
- Bootstrap: O(E) where E = evidence items — this is the bottleneck for large repos

**With 50+ languages:**
- Each language needs: tree-sitter grammar, queries for imports/symbols/api_calls, dependency file parsers
- Current language registry is a Python dict — adding a language means code changes
- A plugin system would be needed for truly extensible language support

### Caching

**Current caching:**
- Registry loaded once into module-level global (`_REGISTRY` in v3.py)
- Embedding model loaded once (lazy) and cached per session
- Embeddings cached per ADR ID in memory

**What's missing:**
- Tree-sitter ASTs could be cached between enforcement runs (modified time check)
- Dependency files could be cached with checksum validation
- No disk-based caching for any of these

### Bottlenecks

1. **Tree-sitter parsing:** Each file is parsed from scratch on every enforcement. For a 1000-file monorepo, this could be 1000+ parse operations.
2. **LLM classification:** Sequential per-(ADR, symbol) call. 15 pairs = ~30s+ with GPT-4o. Cost = ~$0.15-0.50.
3. **Full-repo scan:** `_enforce_repo()` visits every file including binary files and configs — O(total files) path regex + O(language files) import/API scan.
4. **Bootstrap evidence:** Collecting imports from every Python/JS/Go file in the repo is O(source files).

### Deployment

**PyPI:**
```bash
pip install decisiondrift
pip install decisiondrift[ast]      # Multi-language support
pip install decisiondrift[embeddings]  # Semantic retrieval
```

**Docker (GitHub Action):**
```bash
docker pull ghcr.io/madhan-karthikeyan/decisiondrift:latest
```

**Pre-commit hook:**
```bash
decisiondrift guard --install
```

### Monitoring

**No monitoring built in.** The tool is designed for CI, not for production monitoring. Metrics that could be added:
- Number of enforcements run
- Violation rate by ADR
- Bootstrap coverage improvement over time
- Average enforcement time
- LLM cost per review

### Future Improvements

1. **Server-based ADR registry** — team-wide centralized decisions
2. **Real-time IDE integration** — beyond basic save-triggered diagnostics
3. **Architectural diff visualization** — see what changed architecturally between commits
4. **Vulnerability database integration** — detect if a dependency is both prohibited AND has a CVE
5. **Plugin system for languages** — external packages can add language support
6. **Conflict detection between ADRs** — detect when two ADRs contradict each other
7. **ADR templates from code** — generate ADRs from import graphs
8. **Open Policy Agent integration** — for infrastructure-level decisions
9. **PR review dashboard** — track governance compliance over time

---

## Phase 8: Mock Interview

*Instructions for the candidate: Read one question at a time. Answer out loud. Then read the critique. Continue drilling until the topic is exhausted.*

### Level 1 — Resume Walkthrough

> **Q1:** Walk me through the architecture of DecisionDrift. Where do I start?

**What a good answer covers:**
- 5 main entry points (enforce, bootstrap, review, audit, init)
- The rule engine as the core (deterministic, no LLM)
- ADR lifecycle (files in docs/adr/)
- The bootstrap pipeline (evidence → technology → governance → suggestion)
- Output formatting abstraction (ReportEnvelope → 5 formatters)
- CI integration (Docker Action + pre-commit hook)
- Multi-language via optional tree-sitter

**Common mistakes:**
- Starting with bootstrap rather than enforce (enforce is the core value prop)
- Not mentioning that LLM is optional
- Confusing the `review` and `enforce` commands
- Not distinguishing the 3 enforcement modes (diff/file/repo)

**Follow-up:** Which command would you optimize if you had to make one of them 10x faster?

> **Q2:** What problem does this project solve that can't be solved with grep in a CI script?

**What a good answer covers:**
- ADR lifecycle management (approve/reject/supersede/expire)
- Multi-language import analysis (grep for imports is fragile)
- Structured output (SARIF, JSON for tooling)
- Automated ADR generation from repo structure (bootstrap)
- Semantic review via LLM (beyond pattern matching)
- Pre-commit hook integration
- GitHub Action with PR comments + commit status
- Repository role inference and context-aware decisions

**Common mistakes:**
- Claiming the rule engine is more powerful than regex (it's simple substring matching)
- Not acknowledging that a simple bash script could do 80% of what `enforce` does
- Over-promising the LLM capabilities

### Level 2 — Architecture

> **Q3:** How does the enforcement engine handle diff-based vs full-repo scanning? What are the tradeoffs?

**What a good answer covers:**
- Diff mode: parses git diff → `parse_diff()` → `ChangedFile[]` → only scans changed files
- Full-repo mode: uses `repo.rglob("*")` for all files
- Tradeoffs: diff is faster but misses pre-existing violations; full-repo catches drift but is slower
- Audit mode uses full-repo; pre-commit hook uses diff

**Critique if missing:**
- The diff mode silently skips binary files and excluded dirs
- Full-repo mode re-parses every dependency file on each run (no caching)
- Diff mode may miss violations if a prohibited dependency was already in the repo before the ADR was created

### Level 3 — Technology Choices

> **Q4:** Why Click instead of Typer or argparse?

**What a good answer covers:**
- Click was chosen before Typer existed (or before it was mature)
- Click's group nesting is more explicit than argparse
- Click has `CliRunner` for testing (used in `guard` command)
- Community standard at the time of creation

**Critique if missing:**
- Tradeoff: Click's decorator approach makes it harder to compose commands programmatically
- Current code would benefit from Typer's type hints for auto-completion
- 15+ option decorators on `bootstrap` command is unwieldy

### Level 4 — Implementation

> **Q5:** Walk me through how `_to_finding()` works, including the confidence downgrade logic.

**What a good answer covers:**
- Takes a `Rule` and `match_value`, returns `EnforcementFinding`
- Confidence numeric: HIGH=0.9, MEDIUM=0.6, LOW=0.3
- If confidence < 0.50 → action becomes INFO regardless of original
- If confidence < 0.80 and action was BLOCK/REQUIRE_APPROVAL → action becomes WARN
- Severity mapping: BLOCK→critical, REQUIRE_APPROVAL→high, WARN→medium, INFO→low

**Critique if missing:**
- Why these thresholds? (0.50 for the LOW→INFO floor, 0.80 for bootstrap MEDIUM→WARN)
- What happens with custom rules? (Custom rules don't set confidence, default is HIGH)
- This means bootstrap-generated rules (MEDIUM=0.6) will never produce BLOCK findings

**Follow-up:** Is there a race condition where an ADR with LOW confidence still produces a BLOCK finding?

### Level 5 — Edge Cases

> **Q6:** What happens when someone submits a PR that adds a dependency declared as a Git submodule or vendored directly into the repo without a package manager?

**What a good answer covers:**
- Dependency scanner only checks known manifest files (requirements.txt, package.json, etc.)
- A vendored library won't appear in any manifest → dependency scanner won't catch it
- Import scanner would catch it if the code imports from it (both Python AST and tree-sitter)
- Path scanner could catch it with a rule like `vendor/*` or `lib/*`
- Gap: no automatic detection of vendored dependencies
- Recommendation: a "vendored code" rule requiring ADR approval for any new directory in vendor/

**Critique if missing:**
- Not acknowledging the gap is the dangerous answer
- The import scanner only catches the *first* import — if the vendored lib doesn't use the same module path, it's invisible

### Level 6 — Performance

> **Q7:** Estimate the worst-case runtime of `_enforce_repo()` with 5000 files, 200 ADRs, 400 prohibitions, and tree-sitter installed.

**What a good answer covers:**
- O(R × F) with 400 rules spread across 5 types
- File discovery via `rglob` is O(F) ~5000 files (fast, <100ms)
- Dependency scanning: 5000 file checks with dict lookups, only ~50 will be manifest files (fast)
- Import scanning: Python files via ast.parse (~100ms each), other languages via tree-sitter (~200ms each)
- If 500 Python files: 500 × 100ms = 50s (bottleneck!)
- If 500 JS/TS/Go/etc files: add tree-sitter parse time
- API scanning: same set of files, re-parsed
- Path scanning: 5000 files × 400 rules = 2M regex checks (fast)
- Config scanning: subset of files
- Estimated total: 2-5 minutes for worst case

**Critique if missing:**
- Not mentioning that Python AST parsing is the bottleneck
- Not proposing caching (parse files once, cache by mtime)
- Not mentioning that `EXCLUDED_DIRS` helps but doesn't help with large Python monorepos

### Level 7 — Internals

> **Q8:** How does `_match_dep()` in `detectors.py` work? What are all the matching strategies?

**What a good answer covers:**
```python
def _match_dep(dl, dn):
    if dl == dn: return True                    # Exact match
    if dl.startswith(dn + "."): return True     # Submodule (sqlalchemy.orm)
    if dl.startswith(dn + "-"): return True     # Extras (psycopg2-binary)
    if dl.startswith(dn + "_"): return True     # Underscore variant
    if dl.endswith("/" + dn): return True       # Go module (gin-gonic/gin)
    bracket = dl.find("[")
    if bracket > 0 and dl[:bracket] == dn:      # Extras syntax (psycopg[binary])
        return True
    return False
```

**Critique if missing:**
- Not mentioning the false positive risk (`flask-cors` matches `flask` via `startswith`)
- Not mentioning Go module normalization (paths are not lowercased)
- The function doesn't handle scoped npm packages (`@angular/core` — the `@` prefix is preserved)

### Level 8 — Tradeoffs

> **Q9:** The keyword retrieval backend normalizes scores by dividing by `len(search_terms)`. Why is this problematic? How would you fix it?

**What a good answer covers:**
- Problem: symbols with more search terms get lower normalized scores even if they match perfectly
- Problem: a symbol that generates 20 search terms (camelCase parts + path parts + underscore parts) needs 20× the match evidence to score the same as a symbol with 2 search terms
- Fix options: softmax over terms, only count unique matched terms, use max instead of sum, or use a learned weighting

**Critique if missing:**
- Not identifying that the real issue is in `generate_search_terms()` producing too many low-value terms
- Better fix: filter search terms by information content before scoring
- Alternative: use embedding only (no keyword scoring)

### Level 9 — Failure Scenarios

> **Q10:** What happens when tree-sitter silently fails on a file? Can a developer bypass governance by introducing syntax errors?

**What a good answer covers:**
- Yes, this is a real vulnerability in the current design
- `extract_imports_treesitter()` catches all exceptions and returns an empty list
- `extract_api_calls_treesitter()` does the same
- A developer could add `// SyntaxError: ` or an equivalent to a JS file and tree-sitter would fail to parse it
- The file would produce zero imports and zero function calls, so no rules would trigger
- This is especially dangerous in CI where the diff shows the file was modified but the scanner reports no findings
- Fix: if a file fails to parse, it should be reported as an "unparseable" finding or at least logged

**Critique if missing:**
- Not acknowledging the severity of this issue
- Proposing a remediation: add a `--fail-on-parse-error` flag
- The same issue exists in `ast_python.py`'s `extract_symbols()` (but at least SyntaxError is caught and logged)

### Level 10 — Research-Level

> **Q11:** The V3 bootstrap's `_aggregate_evidence_level()` requires 2+ runtime STRONG evidence pieces for STRONG level. Give concrete examples where this heuristic fails.

**What a good answer covers:**
- Fails for monorepo microservices: each service has 1 runtime STRONG dep, aggregated as MODERATE even though the org as a whole has 10 services using it
- Fails for languages: Go is detected by `go.mod` (file evidence, STRONG) and `import` (runtime, MODERATE). Language evidence is treated differently — the language evidence `*.go` glob is STRONG but might not count toward the 2-runtime requirement depending on role inference
- Fails for Next.js: `package.json` with `next` as dependency gives MODERATE, but there's no separate runtime evidence. It would stay MODERATE even though it's clearly the main framework
- Fails for libraries like SQLAlchemy: a single strong dependency + strong import = 2 pieces, but the import is usually in a single file. If the import is in a test file, its role is TEST, not RUNTIME

**Critique if missing:**
- Not suggesting an alternative aggregation strategy
- The `has_runtime_strong` check requires `role == RUNTIME` AND `level == STRONG` — this is a double filter that may exclude evidence from other roles
- Better approach: Bayesian confidence aggregation with prior probabilities per technology

---

## Phase 9: Brutal Interview Questions

### Easy (1-10)

1. **What problem does DecisionDrift solve that can't be solved with grep in a CI script?**
   - *Tests:* Understanding of core value beyond pattern matching
   - *Excellent:* ADR lifecycle, multi-language, structured output, semantic review, pre-commit hook, CI integration
   - *Common mistake:* Claiming the rule engine is more sophisticated than it is

2. **Why does each prohibition generate both a DEPENDENCY and an IMPORT rule?**
   - *Tests:* Defense-in-depth design thinking
   - *Excellent:* Dependencies catch manifest declarations, imports catch direct usage even without manifest
   - *Follow-up:* What kind of violation would each catch that the other misses?

3. **How are ADRs stored? Why this format?**
   - *Tests:* Understanding of filesystem-based storage vs database
   - *Excellent:* Markdown + YAML frontmatter in `docs/adr/ADR-NNNN.md`. Versioned with git, human-readable, no database dependency.

4. **What's the difference between `enforce` and `review`?**
   - *Tests:* Understanding deterministic vs LLM paths
   - *Excellent:* `enforce` is deterministic (no LLM), 5 rule types, runs in CI. `review` uses LLM, retrieves relevant ADRs, classifies changes semantically, needs API key.
   - *Common mistake:* Saying review replaces enforce

5. **How does the pre-commit hook work?**
   - *Tests:* Understanding git hooks
   - *Excellent:* Bash wrapper that runs `decisiondrift enforce --staged`. Installed via `guard --install`. Blocks commit on violations.

6. **What are the 5 rule types?**
   - *Tests:* Understanding scope of enforcement
   - *Excellent:* DEPENDENCY (manifest files), IMPORT (source code), API (function calls), PATH (file paths), CONFIG (config file values)

7. **What does the `doctor` command check?**
   - *Tests:* Understanding health diagnostics
   - *Excellent:* CLI version, config file, registry, tree-sitter, embeddings, ADR directory, LLM connectivity. Returns structured report.

8. **How does the audit command detect drift?**
   - *Tests:* Understanding reuse of rule engine
   - *Excellent:* Runs `enforce_from_adrs()` in full-repo mode against accepted ADRs. Any findings = drift. Also checks expiry, staleness, coverage.

9. **What's the format of ADR-0001.md?**
   - *Tests:* Understanding ADR schema
   - *Excellent:* YAML frontmatter with `id`, `title`, `status`, `severity`, `prohibitions`, etc. Markdown body with context, rationale, evidence.

10. **How does the `init` command differ from `bootstrap`?**
    - *Tests:* Understanding orchestration vs single step
    - *Excellent:* `init` runs bootstrap AND interactively approves/rejects AND installs pre-commit hook AND generates `decisiondrift.yml` AND optionally generates CI workflow. `bootstrap` only generates candidate ADRs.

### Medium (11-25)

11. **How does keyword retrieval score ADRs? What are the multipliers and why?**
    - *Tests:* Understanding retrieval weighting
    - *Excellent:* Title=3× (most authoritative), Keywords=3×, Evidence paths=2×, Rationale=1×. Exceptions=-1× (penalty). Normalized by term count.

12. **In the embedding backend, what similarity metric is used and how is it computed?**
    - *Tests:* Understanding vector similarity
    - *Excellent:* Cosine similarity. `dot(a,b) / (norm(a) * norm(b))`. Implemented via numpy.

13. **Why does `_enforce_diff` use different scanning logic than `_enforce_repo`?**
    - *Tests:* Understanding optimization for context
    - *Excellent:* Diff mode only scans changed files (faster). Repo mode scans everything (thorough). Diff uses `parse_diff()` output, repo uses `rglob()`. Both share scanner functions.

14. **How does the repository role inference work?**
    - *Tests:* Understanding heuristic classification
    - *Excellent:* Checks runtime frameworks (FastAPI/Flask/Express → api_service), frontend-only → frontend_app, build manifest → library, multiple scopes + frameworks → monorepo. Default → unknown.

15. **What is the evidence aggregation logic that determines WEAK/MODERATE/STRONG?**
    - *Tests:* Understanding evidence levels
    - *Excellent:* `_aggregate_evidence_level()`: 2+ runtime STRONG → STRONG. Any runtime STRONG or 2+ runtime MODERATE → MODERATE. Any runtime → MODERATE. Otherwise WEAK.

16. **How does the enforceability analysis decide a candidate is enforceable?**
    - *Tests:* Understanding governance gating
    - *Excellent:* Must have a governance template (from registry or LLM). Must not be tooling or supporting-only. Must have prohibitions that map to deterministic rules. Evidence must be >= MODERATE.

17. **Why does the `guard` command use `CliRunner` internally instead of calling the CLI function directly?**
    - *Tests:* Understanding CLI testing patterns
    - *Excellent:* Using `CliRunner` simulates a real CLI invocation including Click context/error handling. Direct function calls bypass Click's parameter processing and exit behavior.

18. **How does the GitHub action handle SARIF output?**
    - *Tests:* Understanding SARIF integration
    - *Excellent:* Transforms `ReportEnvelope` findings into SARIF v2.1.0 format via `_format_sarif()`. Maps action→level (block=error, warn=warning, info=note). Writes to file specified by `INPUT_SARIF_OUTPUT_PATH`.

19. **What happens when an ADR expires?**
    - *Tests:* Understanding ADR lifecycle
    - *Excellent:* The `audit` command flags expired ADRs (past `expires_after` date). They remain active — `audit` just reports them. No automatic deactivation. Recommendation: use `adr supersede` to replace.

20. **How does the supersession resolution work?**
    - *Tests:* Understanding dependency chain resolution
    - *Excellent:* `resolve_active()` filters accepted ADRs. Removes any that have `superseded_by` set. Removes any whose `depends_on` is not accepted. Returns clean list.

21. **What's the purpose of `scope_path` in evidence collection?**
    - *Tests:* Understanding monorepo support
    - *Excellent:* First directory segment of evidence path. Used to group evidence by sub-project within a monorepo. Enables per-service governance candidates.

22. **How does the V3 bootstrap handle duplicate ADRs?**
    - *Tests:* Understanding dedup logic
    - *Excellent:* Uses `_is_duplicate_title()` with Jaccard word similarity + keyword overlap. Weighted: 60% title similarity, 20% keyword matches, 20% keyword set overlap. Threshold: 0.5.

23. **What happens when no technology registry matches a dependency?**
    - *Tests:* Understanding fallback behavior
    - *Excellent:* Without `--llm`, the unmatched dependency is ignored (no technology candidate created). With `--llm`, it's sent to `KnowledgeProvider.recognize_technology()` which asks the LLM.

24. **How does the LLM client handle API failures?**
    - *Tests:* Understanding resilience
    - *Excellent:* `complete_json()` raises `LLMResponseError`. `Classifier.classify_one()` catches it and returns `needs_human_review` finding. No retry by default (configurable via `max_retries`).

25. **What is the `EvidenceRole` system and why does it exist?**
    - *Tests:* Understanding context-aware evidence
    - *Excellent:* 7 roles: RUNTIME, DEV, TEST, EXAMPLE, TOOLING, OPTIONAL, UNKNOWN. Inferred from path segments (test/ → TEST, docs/ → EXAMPLE). Controls whether evidence contributes to governance (only RUNTIME does). Prevents false positives from test/dev dependencies.

### Hard (26-40)

26. **The keyword retrieval backend normalizes scores by dividing by `len(search_terms)`. Why does this cause unfair scoring between different symbols? Design a fairer approach.**
    - *Tests:* Deep algorithm analysis
    - *Excellent:* More terms = lower normalized score for same matches. Fairer: max-score (best match wins), threshold-based (binary above/below), or learned weights per term type.

27. **The `_enforce_diff` function re-parses dependency files on every invocation. How would you cache parsed dependencies across multiple invocations in a CI context?**
    - *Tests:* Caching + CI optimization
    - *Excellent:* Use file checksum (SHA256) as cache key. Store in `.decisiondrift/cache/`. Invalidate only when file content changes. Use `stat().st_mtime` for fast comparison.

28. **The repository role inference in `infer_repository_role()` uses hardcoded framework names. Design an extensible approach that doesn't require code changes for new frameworks.**
    - *Tests:* Extensibility + registry design
    - *Excellent:* Add `role_signals` to technology registry entries (e.g., `FastAPI: { role_signal: "api_service" }`). Let `infer_repository_role()` query the registry instead of hardcoded names. New framework = new registry entry.

29. **The current tree-sitter import extraction doesn't handle dynamic imports (`const x = import('module')` in JS or `importlib.import_module()` in Python). How would you add support?**
    - *Tests:* AST analysis + tree-sitter queries
    - *Excellent:* Add a new query pattern for dynamic imports. In JS: `(call_expression (import) ...)`. In Python: detect `Call` nodes with `func=Attribute(attr='import_module')`. Parse first argument as module name.

30. **The LLM classification is pairwise: one (ADR, symbol) per call. Design a batched approach and explain the tradeoffs.**
    - *Tests:* LLM optimization
    - *Excellent:* Batch multiple pairs in one prompt. Tradeoffs: larger context window (cost), potential confusion between pairs, less granular error handling, but 5-10x faster. Approach: limit batch to 3-5 pairs, use structured output format.

31. **How would you implement incremental enforcement where only files changed since the last clean scan are checked?**
    - *Tests:* Incremental processing
    - *Excellent:* Track per-file mtime + checksum. Store "last clean" state per file. Only scan files with changed content. Fall back to full scan when ADRs change.

32. **The evidence collection only scans `.py`, `.js`, `.ts`, and `.go` files for imports. Design a general approach for any language supported by tree-sitter.**
    - *Tests:* Generic import extraction
    - *Excellent:* Loop over all languages in `LANGUAGE_REGISTRY`, find files by extension, run tree-sitter import query. Cache results by file. Add timeout per file.

33. **How does the confidence downgrade logic in `_to_finding()` interact with severity? Can a BLOCK rule produce a "low" severity finding?**
    - *Tests:* Understanding the mapping chain
    - *Excellent:* No — confidence downgrade affects `action` (BLOCK→WARN), and severity is derived from the *downgraded* action. BLOCK→WARN maps to "medium" severity. LOW confidence (0.3) downgrades BLOCK→INFO which maps to "low" severity. So yes, indirectly.

34. **The `load_registry()` function merges layers but doesn't detect conflicts. How would you add conflict detection and resolution?**
    - *Tests:* Configuration merging
    - *Excellent:* Before merge, check for key conflicts. Log warnings. Add resolution strategies: "original-wins" (bundled takes priority), "override-wins" (last takes priority), or "fail-on-conflict". Make strategy configurable.

35. **What's the worst-case performance of `_enforce_repo` with 500 ADRs and 10,000 files?**
    - *Tests:* Performance estimation
    - *Excellent:* File discovery ~100ms. Dependency scanning ~50ms (most files are not dep files). Import scanning: if 1000 Python files at ~100ms each = 100s bottleneck. Tree-sitter languages at ~200ms each additional. Total: 2-3 minutes.

36. **The init command runs bootstrap, then approve/reject, then hook install, then config write. How would you make this idempotent?**
    - *Tests:* Idempotent design
    - *Excellent:* Check each step before running: skip bootstrap if ADRs exist, skip hook if installed, skip config if exists, skip CI if exists. Add `--force` to override. Ensure re-running produces same result.

37. **How would you add support for YAML-based architecture decisions that are not stored as ADR files?**
    - *Tests:* Alternative storage backends
    - *Excellent:* Create an abstract `ADRStore` interface. Implement `FileStore` (current), `YAMLStore` (single YAML with all ADRs), `DatabaseStore` (SQLite). Load via config. Each store returns `DecisionRecord` list.

38. **The current ADR schema is validated with JSON Schema. How would you support schema migration (adding/removing fields)?**
    - *Tests:* Schema evolution
    - *Excellent:* Version the ADR schema in frontmatter (`schema_version: 2`). Write migration scripts (v1→v2). Use JSON Schema `$schema` field. Loader detects version and applies migration before validation.

39. **How would you secure the LLM API key in the GitHub Action? The current approach passes it as an input.**
    - *Tests:* CI security
    - *Excellent:* Current approach exposes key in workflow YAML. Better: GitHub Secrets for `DECISIONDRIFT_LLM_API_KEY` environment variable. Never pass as input. Mask in logs. Use OIDC for keyless auth when possible.

40. **How does the project handle circular supersession (A supersedes B, B supersedes A)?**
    - *Tests:* Cycle detection
    - *Excellent:* It doesn't — `resolve_active()` has no cycle detection. A cycle would cause both to be filtered out (since superseded_by is checked). But it could infinite loop if supersession chain is traversed recursively. Fix: add visited-set detection.

### Expert (41-50)

41. **The V3 bootstrap's `_aggregate_evidence_level` requires 2+ runtime pieces at STRONG level to achieve STRONG. Is this heuristic sound? Provide concrete examples where it fails.**
    - *Tests:* Evidence-based reasoning
    - *Excellent:* Fails for: (a) single-file Python library with one `import requests` — MODERATE even though usage is STRONG. (b) Next.js with only `package.json` dependency — MODERATE despite being primary framework. (c) Monorepo with the same technology in every service but only one evidence per service.

42. **The `_suppression_reason` logic suppresses candidates when contradictions exist. Design a system that instead downgrades confidence proportionally to contradiction severity.**
    - *Tests:* Confidence-based filtering
    - *Excellent:* Weight contradictions by severity. Each contradiction reduces a multiplier (e.g., "no runtime evidence" = 0.5×, "only favicon evidence" = 0.8×). Final confidence = base × multiplier. Present ADR with adjusted confidence instead of suppressing entirely.

43. **The evidence collection scans all files in the repo, including vendored dependencies in `node_modules` (excluded) but not in `vendor/` dirs for Go. Propose a robust vendor detection system.**
    - *Tests:* Repository scanning
    - *Excellent:* Detect `vendor/`, `third_party/`, `deps/`, `_deps/` directories. Check for known vendoring tools (go mod vendor, gitsubmodule, yarn PnP). Add user-configurable exclude patterns in `decisiondrift.yml`. Exclude by default unless configured otherwise.

44. **How would you architect a real-time IDE extension that provides decisions as the user types (not just on save)?**
    - *Tests:* Real-time analysis
    - *Excellent:* Use LSP (Language Server Protocol) or tree-sitter's incremental parsing. Watch ADR files for changes. Cache parsed decision rules in extension. Debounce analysis (500ms). Use `--file` mode for single-file analysis. Show violations as editor diagnostics. Background cache refresh.

45. **The project claims "363 tests passing, 26 skipped". The skipped tests are tree-sitter tests. How would you make tree-sitter tests reliable across different CI environments?**
    - *Tests:* Testing infrastructure
    - *Excellent:* Use pytest markers for tree-sitter tests. In CI, install `[ast]` extra. For offline environments, vendor tree-sitter grammars. Run tree-sitter tests in a separate CI job that ensures the dependency is available. Use `pytest.skipif` with `HAS_TREESITTER` flag.

46. **The keyword retrieval scores title matches at 3×, keywords at 3×, evidence paths at 2×, rationale at 1×. These weights were presumably chosen experimentally. Design an approach to learn optimal weights from labeled data.**
    - *Tests:* Learning + optimization
    - *Excellent:* Collect labeled data: for each (search_term, ADR) pair, label as relevant/irrelevant. Use grid search or Bayesian optimization over weight space. Optimize for NDCG@5 or MAP (Mean Average Precision). Evaluate via cross-validation.

47. **How does the project handle ADRs that prohibit a technology that's already in use? The enforce command would immediately fail. Design a graceful onboarding period.**
    - *Tests:* Governance rollout
    - *Excellent:* Add ADR field `grace_period: "30 days"` or `grace_until: "2026-09-01"`. In `_to_finding()`, check grace period. If within grace, downgrade BLOCK→WARN. After grace, enforce fully. Generate migration plan during bootstrap.

48. **The embedding backend uses `BAAI/bge-small-en-v1.5` which is English-only. How would you support Chinese ADRs?**
    - *Tests:* Multilingual embedding
    - *Excellent:* Detect ADR language (via `langdetect` or heuristic). Use multilingual model like `intfloat/multilingual-e5-small` or `BAAI/bge-base-zh-v1.5`. Make model configurable in `decisiondrift.yml` with per-language selection.

49. **The current output formatter produces 5 formats. Design a plugin system for custom formatters.**
    - *Tests:* Plugin architecture
    - *Excellent:* Define `Formatter` protocol/ABC with `format(envelope: ReportEnvelope) -> str`. Use entry_points or `decisiondrift.formatters` namespace. Load via `importlib.metadata.entry_points()`. User configures `format: my-custom-plugin`.

50. **How would you add a `decisiondrift diff <adr-id>` command that shows only the diff lines relevant to a specific ADR?**
    - *Tests:* Cross-cutting feature design
    - *Excellent:* Generate rules for that ADR. Run enforcement in diff mode against the working tree. Filter findings to only those matching the ADR's rule IDs. Display the corresponding diff hunks (from git diff) with ADR annotations. Highlight the violating lines with the ADR's rationale.

---

## Phase 10: Knowledge Gaps & Weaknesses

### Ranked by Interview Risk

| Risk | Issue | Code Evidence | Likely Attack Question | How to Defend |
|---|---|---|---|---|
| **CRITICAL** | Silently failing tree-sitter | `ast_treesitter.py:30-63,90-120,165-208` — all return empty lists on any error | "What happens if I submit a PR with syntax errors? Can I bypass governance?" | "Yes, this is a known gap. Tree-sitter exceptions are silently caught. A `--fail-on-parse-error` flag would address this. In practice, code reviews catch syntax errors, but we should add this safeguard." |
| **CRITICAL** | Substring matching false positives | `detectors.py:128-148`: `startswith(dn)` matches too broadly | "If I prohibit 'flask', what happens with 'flask-cors' or 'flask-security'?" | "flask-cors would match because `startswith('flask')` is true. We accept this tradeoff because (a) false positives catch related packages that should be reviewed, (b) the ADR can add exceptions, and (c) exact name matching would miss fork/variant packages." |
| **HIGH** | No ADR versioning/conflict detection | `adr/supersession.py`: simple linear supersession | "What if two ADRs prohibit the same thing with different severities?" | "The rule engine would generate two rules with different severities. Both would trigger. There is no conflict resolution. A future improvement would be to detect conflicts during `adr approve` and warn." |
| **HIGH** | Single-threaded enforcement | `rules/engine.py`: sequential file scanning | "How do you handle a 50,000 file monorepo?" | "Currently single-threaded. For large repos, we'd (a) parallelize file scanning per rule type, (b) cache parsed ASTs by mtime, (c) add incremental enforcement — only scan files changed since last run." |
| **HIGH** | Embedding model hardcoded | `config.py:40`: `BAAI/bge-small-en-v1.5` | "What if I want to use a multilingual or code-specific embedding model?" | "The model name is configurable in `decisiondrift.yml`. Currently only English. For multilingual needs, you'd configure `intfloat/multilingual-e5-small`. The `EmbeddingBackend` accepts a `model_name` parameter." |
| **MEDIUM** | No line numbers in SARIF | `report/formatter.py:70-73` — only `artifactLocation.uri`, no `region` | "How does this integrate with GitHub code annotations?" | "This is a limitation. The current SARIF output only identifies the file, not the specific line. Line-level annotation would require `physicalLocation.region.startLine`. This would be a straightforward enhancement." |
| **MEDIUM** | HTTP registry no caching | `registry.py:81-89`: no ETag/If-Modified-Since | "How many times does this fetch the same URL in CI?" | "Every invocation of `load_registry()`. This is a performance issue for CI. We should add ETag-based caching with a `~/.cache/decisiondrift/registries/` directory." |
| **MEDIUM** | ADR quality scores are field-based | `cli.py:656-675`: checks field presence, not quality | "My ADR has all fields filled but the rationale is 'use flask'. Score is 100 but it's useless." | "Correct — the current quality score is structural, not semantic. A useful ADR quality metric would also check: rationale length > 20 words, prohibitions are specific, evidence is present, exceptions are documented." |
| **MEDIUM** | No diff parsing for package lock files | `rules/engine.py:94` — only checks requirements.txt and pyproject.toml | "What if I change a dependency in package-lock.json only?" | "If the dependency change isn't reflected in `requirements.txt` or `pyproject.toml`, it won't be caught by the dependency scanner. The import scanner might still catch it if the code imports the package. Best practice: always declare dependencies in source manifest files." |
| **MEDIUM** | Evidence role inference is path-segment based | `v3.py:1083-1101`: maps path segments to roles | "What if my test directory is named `specs/` instead of `tests/`?" | "Then evidence in `specs/` would get UNKNOWN role instead of TEST. The role inference is heuristic and tuned for common naming conventions. For non-standard layouts, evidence roles may be misclassified." |
| **LOW** | No input validation in GitHub action | `github/action_entrypoint.py:84-95` — direct `os.environ.get()` | "What if someone sets INPUT_MAX_PAIRS to -1?" | "`int(os.environ.get(...))` would raise ValueError. We should use `try/except` or a validation function. In practice, GitHub Actions sets default values so this is rare." |
| **LOW** | Global module state in v3.py | `v3.py:29-30`: `_REGISTRY = None` | "How does this behave with multiple registries across different repos in the same process?" | "The global is reset when `_get_registry()` is called with different URLs. This is safe in CLI mode (single process, single repo). In server mode it would need to be per-request." |
| **LOW** | Duplicate dependency parsing logic | `detectors.py` has its own scanners, `utils/dependency_parser.py` has another set, `rules/scanner.py` has references | "Are the dependency parsers in sync across these files?" | "There may be drift. The `detectors.py` scanners are the original implementation. `utils/dependency_parser.py` was added later for shared use. An audit would ensure consistency." |

### Code Smells

1. **Global state pattern in v3.py**: `_REGISTRY`, `TECH_SIGNATURES`, `DECISION_TEMPLATES`, `SELF_FRAMEWORK_REPOS`, `TEST_DEPENDENCIES`, `TOOLING_DECISIONS`, `SUPPORTING_ONLY_DECISIONS` are all global mutable state. Populated once and never reset. Not thread-safe.

2. **Duplicate dependency parsing logic**: Three separate implementations of dependency file parsers across `detectors.py`, `rules/scanner.py` + `rules/engine.py`, and `utils/dependency_parser.py`. Risk of inconsistencies.

3. **Exception swallowing**: `ast_treesitter.py`, `rules/engine.py`, `scanner.py`, `v3.py` — most exceptions are silently caught and return empty lists or `None`. Debugging is extremely difficult.

4. **Long modules**: `cli.py` at 888 lines, `v3.py` at 1219 lines. These should be split into smaller, focused modules.

5. **Mixed responsibilities in v3.py**: Contains evidence collection, technology candidate building, role inference, governance discovery, enforceability analysis, suggestion generation, file iteration helpers, and evidence level logic. At least 5 separate concerns in one file.

6. **Hardcoded magic values**: String literals like `"api_service"`, `"frontend_app"`, `"framework"` appear throughout v3.py without constants. Framework names are hardcoded in `infer_repository_role()` instead of being registry-driven.

### Architecture Weaknesses

1. **No ADR conflict resolution**: Two ADRs can prohibit the same thing with different severities and both will generate rules. No detection or resolution.

2. **No database**: All state is filesystem-based. This works for CLI but prevents collaborative features, dashboard, or cross-repo analytics.

3. **Monolithic bootstrap**: V3 pipeline is 1219 lines and tightly coupled. Evidence collection, technology detection, role inference, and governance discovery cannot be tested independently without running the full pipeline.

4. **No caching layer**: Tree-sitter ASTs, parsed dependency files, and registry responses are never cached to disk. Every invocation re-parses everything.

5. **Limited extensibility for languages**: Adding a new language requires modifying `language_registry.py`, creating a query file, and possibly updating `extract_imports_treesitter()` if the import segment logic needs customization.

### Scalability Concerns

1. **Full-repo scan performance**: `_enforce_repo()` visits every file. For a 50K-file monorepo, this could take minutes.
2. **Memory usage**: No streaming — all evidence is collected into memory before processing.
3. **No incremental mode**: Every enforcement is a full scan. No way to say "only check files changed since last clean run."
4. **Tree-sitter per-file overhead**: Each file parse creates a new tree-sitter parser instance (not reused across files).

### Security Concerns

1. **LLM API key exposure**: Passed as GitHub Action input, visible in workflow logs unless masked.
2. **No input sanitization in action entrypoint**: Environment variables directly parsed without validation.
3. **Silent tree-sitter failures**: Allows governance bypass via parse errors.
4. **No sandboxing for file scanning**: In theory, symlink traversal could escape the repo directory (though `is_symlink()` check mitigates this).

### Maintainability Issues

1. **Global state**: Module-level mutable state in v3.py makes testing in isolation difficult.
2. **Duplicate code**: Three versions of dependency parsers.
3. **Long functions**: `cli.py` has functions over 100 lines. `v3.py` has many 50+ line functions.
4. **Mixed abstractions**: v3.py mixes high-level (repository modeling) with low-level (file iteration, regex matching).
5. **No interface contracts**: Core subsystems (bootstrap, enforcement, retrieval, classification) lack formal interface definitions beyond informal function signatures.

---

*This guide was generated by an automated repository analysis. All code references are to the actual DecisionDrift repository. Distinguish facts from assumptions: the code is the source of truth.*
