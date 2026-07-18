# DecisionDrift — Configuration Drift Detection Engine

## Overview

DecisionDrift continuously monitors infrastructure configuration state across servers, network devices, and cloud resources. It captures immutable snapshots on a schedule or in response to events, computes diffs against baselines, evaluates them against policy rules, and generates reports. Built for air-gapped and regulated environments where agent-to-server communication must be pull-based and auditable.

## Architecture

```
┌─────────────────────────────────────┐
│          Management Server           │
│  ┌──────────┐  ┌──────────────────┐ │
│  │ Scheduler│──▶ Snapshot Registry │ │
│  │ (cron +  │  │ (content-addr    │ │
│  │  events) │  │  storage)        │ │
│  └────┬─────┘  └────────┬─────────┘ │
│       │                  │          │
│       ▼                  ▼          │
│  ┌──────────┐  ┌──────────────────┐ │
│  │ Diff     │◀─│ Baseline Store   │ │
│  │ Engine   │  │ (versioned)      │ │
│  └────┬─────┘  └──────────────────┘ │
│       │                             │
│       ▼                             │
│  ┌──────────┐  ┌──────────────────┐ │
│  │ Policy   │──▶ Reporter         │ │
│  │ Engine   │  │ (HTML, JSON,     │ │
│  │ (OPA-like)│  │  Slack, PagerDuty)│ │
│  └──────────┘  └──────────────────┘ │
└──────────────────┬──────────────────┘
                   │ (pull, HTTPS)
                   ▼
┌─────────────────────────────────────┐
│          Agent (per node)            │
│  ┌──────────┐  ┌──────────────────┐ │
│  │ Collector│──▶ Snapshot Builder  │ │
│  │ (OS,     │  │ (hash + sign)    │ │
│  │  cloud)  │  └──────────────────┘ │
│  └──────────┘                       │
└─────────────────────────────────────┘
```

## Why This Design

Air-gapped environments prohibit agents from phoning home. The pull-based model means the management server connects to each agent on a schedule, fetches the current snapshot, and processes it centrally. This keeps the agent minimal — it only needs to collect state and produce a signed, content-addressed blob. The server handles all diffing, policy evaluation, and alerting. This also means agents never need to know about policies or baselines, reducing their attack surface and update frequency.

## Key Technical Decisions

- **Immutable snapshots with content-addressable storage**: Each snapshot is identified by a SHA-256 hash of its content. The snapshot registry stores blobs by hash and maintains a separate index mapping (node_id, timestamp) → hash. Immutability guarantees that historical comparisons are always faithful — no snapshot can be modified after the fact. Deduplication is a side benefit: identical configs share one blob, reducing storage by ~70% in practice.

- **Diff engine using Myers algorithm**: Standard Myers diff (as used in Git) operates on lines. Infrastructure configs are structured (YAML, JSON, HCL), so a line-level diff loses semantic context. DecisionDrift's diff engine parses structured configs into canonical ASTs, flattens them into key paths, and runs Myers on the path-value pairs. The result is a structured diff that shows exactly which key changed, from what, to what, making it human-readable and machine-actionable.

- **Policy-as-code with OPA-like rule engine**: Policies are written in a Rego-inspired DSL that runs on the diff output. A policy rule like `deny if diff.path matches "firewall.rules.*" and diff.operation == "delete"` blocks the change in enforcement mode or flags it in audit mode. The rule engine is sandboxed (no filesystem, no network) and evaluated per-diff with a configurable budget (max 500ms per rule set) to prevent runaway policies.

## Testing Strategy

- **Snapshot comparison tests**: Golden snapshot files (known-good config states) are stored in the test repository. Tests verify that the diff engine produces the expected structured diff between two snapshots. Tests cover: no-op (identical snapshots), single value change, list insertion/deletion, nested map changes, and large configs with thousands of keys.

- **Chaos testing with random config mutations**: A fuzzer generates a valid base config, then applies random mutations (add/remove/modify keys, corrupt values, reorder lists). The diff engine must correctly identify all mutations. The number of false positives and false negatives is tracked over thousands of runs. This caught edge cases where list reordering was reported as a change when it was semantically identical.

- **Performance benchmarks for large configs (10k+ resources)**: CI benchmarks measure: snapshot ingestion time, diff computation time, policy evaluation time, and storage size for configs with 100, 1,000, 10k, and 100k resources. Results are compared against the previous release. A 10% regression in p50 fails the build. Current performance: 10k-resource diff completes in ~400ms, policy evaluation in ~200ms.

## Failures & Lessons

- **Initial polling interval missed changes**: The scheduler collected snapshots every 60 minutes. Changes that reverted within that window (e.g., a firewall rule added then removed by an attacker) were invisible. Switched to a hybrid model: periodic polling plus event-driven triggers. Agents expose a minimal webhook receiver that the management server calls when a watchable event occurs (file change via inotify, cloud resource change via AWS EventBridge). The snapshot is taken immediately, and the diff is computed against the last known baseline.

- **JSON diff was unreadable**: First version produced standard JSON patch (RFC 6902). While machine-parseable, it was indecipherable in Slack notifications and HTML reports. Implemented a structured diff renderer that groups changes by resource type, shows a before/after side-by-side view for value changes, and collapses unchanged subtrees by default with an expand action. This dropped the average report length from 800 lines to 40 lines for a typical change window.

## Tradeoffs

- **Agent overhead vs polling frequency**: A lighter agent consumes fewer node resources but limits how often you can poll without impacting production workloads. The current agent uses ~15MB RAM and ~1% CPU during snapshot collection. DNS resolution, which is the heaviest single collector, runs with a configurable timeout and concurrency limit. Heavier collectors (full filesystem scan) are opt-in and gated by a node's resource class.

- **Storage size vs history depth**: Content-addressable storage deduplicates aggressively, but history depth still grows linearly with change frequency. A cluster with 5,000 nodes and daily changes generates ~18GB/year of unique snapshot blobs. Mitigations: configurable retention policies (keep all snapshots for 90 days, then daily for 1 year, then weekly), compression (Zstd, ratio ~4:1 on config blobs), and a tiered backend (hot: local SSD, cold: S3/Blob).

## Interview Questions

**Q: How do you detect configuration drift efficiently?**
A: The critical optimization is not diffing everything every time. We maintain a fingerprint (Bloom filter of key-value pairs) for each snapshot. Before running the full Myers diff, we check whether the fingerprints differ. If they match, we skip the diff entirely. For fingerprints that differ, we use the Bloom filter to identify candidate changed keys and only diff those subtrees. This reduces the average diff workload by ~85% in stable environments.

**Q: What consistency model does your snapshot system use?**
A: Snapshots are collected at a point in time but are not transactional across resources. On a single node, we snapshot in a consistent order (kernel params first, then OS config, then application config, then cloud metadata) so that dependencies are captured together. Cross-node consistency is achieved via snapshot coordination — the scheduler can request snapshots from a group of nodes within a time window. The consistency model is documented as "eventually consistent across nodes, per-node point-in-time consistent." Users who need cross-node atomic snapshots must use an external orchestrator (e.g., Kubernetes backup hooks).

**Q: How would you scale to 100k nodes?**
A: The management server is stateless for the diff engine — all state lives in the snapshot registry database. To scale, shard nodes by region or team into independent scheduler instances, each with its own registry. The pull-based model means agents don't need to know about sharding. For cross-shard queries (e.g., "find all nodes where SSH key changed"), we use a secondary index in a search cluster. Each scheduler publishes snapshot metadata (node_id, timestamp, hash, key change flags) to a shared topic. The search cluster consumes this topic and builds a global search index. Actual snapshot blobs remain in the local registry to keep the hot path local.

## Related Problems from This Repository

- **Config Parser**: DecisionDrift's structured config parser (YAML/JSON/HCL to canonical AST) was adapted from the Config Parser problem's solution. The recursive descent parser for HCL was directly reused.
- **LRU Cache**: The snapshot registry uses an LRU cache for hot snapshot blobs (most recent 100 per node). The implementation was ported from this repo's LRU Cache solution.
- **Database CRUD**: The snapshot metadata index is backed by PostgreSQL, and the CRUD patterns (insert snapshot, query by node/timestamp, paginate) follow the same design as the Database CRUD problem.
- **Expression Evaluator**: The policy-as-code rule engine's expression evaluator (comparisons, boolean logic, path matching) reuses the tree-walking interpreter from the Expression Evaluator problem.
