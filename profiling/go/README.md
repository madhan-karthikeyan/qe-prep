# Go Profiling Examples

This directory contains Go profiling examples using `runtime/pprof`.

## Running

```bash
# From the repo root
cd profiling/go
go run main.go
# or (simpler CPU-only example)
go run ./simple/
```

## Profiling Output

Both examples produce:
- `cpu.prof` — CPU profile
- `mem.prof` — heap profile (main.go only)

## Viewing Profiles

### Interactive web UI (recommended)

```bash
go tool pprof -http=:8080 cpu.prof
go tool pprof -http=:8080 mem.prof
```

This opens a browser with flame graphs, call graphs, top functions, etc.

### CLI mode

```bash
go tool pprof cpu.prof
(pprof) top20
(pprof) list main.runLoggerOps
(pprof) web
```

### Comparing profiles

```bash
go tool pprof -http=:8080 --base cpu.prof cpu2.prof
```

## Interpreting Go Profiles

### CPU profile (`cpu.prof`)

Shows where the program spends most of its CPU time. Key commands in the
interactive pprof shell:

| Command        | Description                        |
|----------------|------------------------------------|
| `top`          | Top functions by flat time         |
| `list <func>`  | Show source with per-line timing   |
| `web`          | Open call graph in browser         |
| `peek <func>`  | Show callers/callees of a function |
| `traces`       | Show execution traces              |

### Memory profile (`mem.prof`)

Shows heap allocations. Use `-alloc_space` or `-inuse_space` to focus on
allocation rate vs. live memory:

```bash
go tool pprof -alloc_space mem.prof
go tool pprof -inuse_space mem.prof
```

## Tips

- Run multiple times; profiles are statistical samples.
- Focus on `cum` (cumulative) time in CPU profiles — it includes time spent
  in callees.
- For allocations, `-alloc_objects` shows object count, `-alloc_space` shows
  bytes allocated.
- Use `-base` to diff two profiles and find regressions.
