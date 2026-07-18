# Python Profiling Examples

This directory contains Python profiling examples using `cProfile` and `pstats`.

## Prerequisites

```bash
pip install snakeviz  # optional, for visualisation
```

## Running

### Basic profiling

```bash
cd profiling/python
python profile_logger.py
python profile_lru.py
```

### Using cProfile directly

```bash
python -m cProfile -o logger.prof profile_logger.py
python -m cProfile -o lru.prof profile_lru.py
```

## Viewing Results

### CLI with pstats

The scripts already print sorted stats. For interactive exploration:

```bash
python -c "
import pstats
p = pstats.Stats('logger.prof')
p.sort_stats('cumtime').print_stats(30)
"
```

### Visual with snakeviz

```bash
snakeviz logger.prof
snakeviz lru.prof
```

This opens an interactive flame chart in your browser.

### Visual with gprof2dot

```bash
pip install gprof2dot
gprof2dot -f pstats logger.prof | dot -Tpng -o profile.png
```

## Interpreting Results

Key columns in `pstats` output:

| Column    | Meaning                                      |
|-----------|----------------------------------------------|
| ncalls    | Number of calls (total / primitive)          |
| tottime   | Total time spent in the function alone        |
| percall   | tottime / ncalls                             |
| cumtime   | Total time spent in the function + callees   |
| percall   | cumtime / ncalls (cumulative average)        |

- **cumtime** is usually the most useful sort key — it shows which functions
  consume the most time including their children.
- **ncalls** helps identify functions called excessively.
- Functions with high **tottime** but low **cumtime** are CPU-bound leaf calls.
