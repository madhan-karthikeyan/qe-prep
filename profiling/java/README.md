# Java Profiling Examples

This directory contains Java profiling examples using JFR (Java Flight Recorder).

## Prerequisites

- JDK 17+ (preferably JDK 21 as used by the project)
- Maven (to build the project modules)
- JDK Mission Control (`jmc`) — bundled with Oracle JDK or
  available at https://jdk.java.net/jmc/

## Building Dependencies

First compile the project modules that the profiling example depends on:

```bash
# From the repo root
cd java
mvn package -DskipTests
```

## Running

### Basic run (no profiling)

```bash
# Build classpath from the Maven target directories
CP=$(find ../java -name "target/classes" -type d | tr '\n' ':')
javac -cp "$CP" ProfilingExample.java
java -cp ".:$CP" ProfilingExample
```

### With JFR profiling

```bash
CP=$(find ../java -name "target/classes" -type d | tr '\n' ':')
javac -cp "$CP" ProfilingExample.java

java -XX:StartFlightRecording=filename=record.jfr,duration=30s \
     -cp ".:$CP" ProfilingExample
```

### With JFR and advanced options

```bash
java -XX:StartFlightRecording=filename=record.jfr,\
                                 dumponexit=true,\
                                 maxsize=100MB,\
                                 maxage=1h,\
                                 settings=profile \
     -cp ".:$CP" ProfilingExample
```

## Viewing Results

### JDK Mission Control

```bash
jmc record.jfr
```

In JMC:
1. Open the recording file
2. Use the **Flight Recorder** tab to browse:
   - **Heat Map** — thread activity over time
   - **Method Profiling** — hottest methods (CPU samples)
   - **Memory** — allocation hotspots, GC pauses
   - **Threads** — thread dumps, lock contention
   - **Exceptions** — thrown exceptions

### CLI with jfr tool

```bash
# Print method profiling statistics
jfr summary record.jfr

# Print events
jfr print --events CPULoad record.jfr
jfr print --events ExecutionSample record.jfr
jfr print --events AllocationRequiringGC record.jfr
```

## Interpreting JFR Data

| Tab              | What to look for                              |
|------------------|-----------------------------------------------|
| Method Profiling | Methods with high "Sample Count" — CPU hogs   |
| Memory           | High allocation rate, frequent GC             |
| Lock Instances   | Contended locks — threading bottlenecks       |
| File I/O         | Slow file reads/writes                        |
| Socket I/O       | Network latency                               |
| Exceptions       | High exception throw rates                    |
| Threads          | Blocked threads, context switching            |

## Tips

- Use `settings=profile` for more detailed sampling (higher overhead).
- For allocation profiling, use `-XX:+UnlockDiagnosticVMOptions -XX:+DebugNonSafepoints`.
- Compare multiple recordings to find regressions.
- Use JMC's "Automated Analysis" tab for a quick health check.
