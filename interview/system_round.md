# System Design Round — QE Engineer Interview Guide

## Overview

The system design round for QE focuses on understanding distributed systems tradeoffs, failure modes, and how to test them. You're expected to reason about consistency, availability, partition tolerance, and replication strategies — especially relevant for companies building distributed storage (Couchbase, Cassandra, DynamoDB, etc.).

## Top 20 Questions

### CAP Theorem & Consistency Models

1. **Explain the CAP theorem. Why can't you have all three?**
   - **Expected**: In a network partition (P), you must choose between consistency (C) and availability (A). CP systems block until partition heals; AP systems accept stale reads.
   - **Mistake**: Saying CAP is always "pick 2 of 3" — in practice, you choose C or A only when partitions occur.

2. **What is the difference between strong, eventual, and causal consistency?**
   - **Strong**: All reads see the latest write (linearizability).
   - **Eventual**: Given enough time without writes, all replicas converge.
   - **Causal**: Writes that are causally related are seen in order by all nodes.

3. **When would you choose eventual consistency over strong consistency?**
   - **Expected**: When latency and availability are more important than recency (e.g., social media feeds, DNS). Couchbase uses eventual consistency for cross-datacenter replication.

4. **What is read-your-writes consistency? How is it implemented?**
   - **Expected**: A client always sees its own writes. Implemented with session tokens, version vectors, or reading from primary replica.

### Consensus Algorithms

5. **Explain Paxos in simple terms.**
   - **Expected**: Nodes propose values; an acceptor majority agrees on a single value. Phases: prepare (promise), accept (learn). Core insight: only one value can achieve majority.

6. **How does Raft differ from Paxos?**
   - **Expected**: Raft has a leader (stronger), uses log replication, and divides consensus into leader election, log replication, and safety. Easier to understand and implement.

7. **What is a split-brain scenario? How do consensus algorithms prevent it?**
   - **Expected**: Two nodes each think they're the leader. Raft prevents with randomized election timeouts and requiring a majority of nodes to elect a leader.

8. **What is the role of a quorum in distributed systems?**
   - **Expected**: Minimum number of nodes that must agree for an operation to succeed. Read quorum (R) + Write quorum (W) > N (total replicas) ensures strong consistency.

### Partitioning (Sharding)

9. **What sharding strategies exist? Compare hash vs range partitioning.**
   - **Expected**: Hash partitioning distributes evenly but breaks range queries. Range partitioning supports range scans but can cause hot spots.

| Strategy | Pros | Cons |
|----------|------|------|
| Hash | Even distribution, simple | Poor range queries, resharding hard |
| Range | Efficient range scans | Hot spots, needs load monitoring |
| Consistent hashing | Minimal reshuffling on add/remove | Complexity, virtual nodes needed |

10. **How do you handle resharding without downtime?**
    - **Expected**: Double writes during migration, read from old until migration complete. Or use consistent hashing with virtual nodes. Couchbase uses vBucket mapping.

11. **What is a hot key / hotspot problem?**
    - **Expected**: A single shard receives disproportionate traffic (e.g., a celebrity's profile). Mitigations: cache, split hot key, client-side load shedding.

### Replication

12. **Compare synchronous vs asynchronous replication.**
    - **Synchronous**: Strong consistency, higher latency, lower availability.
    - **Asynchronous**: Lower latency, higher availability, potential data loss on failure.

13. **How does multi-leader (active-active) replication work?**
    - **Expected**: Multiple nodes accept writes, replicate to each other. Conflict resolution strategy needed (last-write-wins, CRDTs, custom merge).

14. **What failure scenarios must a replicated system handle?**
    - **Expected**: Network partition, node crash, disk failure, clock skew, replication lag, zombie replicas.

### Trade-offs and Testing Perspective

15. **How would you test a distributed system's consistency guarantees?**
    - **Expected**: Write under partition, read after heal. Jepsen-style fault injection: partition network, kill leaders, examine anomalies. Property-based testing with invariants: no lost writes, monotonic reads.

16. **What's the difference between read repair and hinted handoff?**
    - **Expected**: Read repair corrects stale data on read (lazy); hinted handoff stores writes for a down node on another node and replays when it's back (eager). Both improve eventual consistency.

17. **Design a test for a distributed queue.**
    - **Expected**: At-least-once delivery (expect duplicates, test idempotency), at-most-once delivery (test no duplicates), ordering guarantees, message TTL, consumer failure, rebalancing.

18. **What consistency guarantee does "QUORUM" read/write provide in Couchbase/Cassandra?**
    - **Expected**: With N=3, R=2 (read quorum), W=2 (write quorum), a read always sees the latest write because at least one node overlaps between read and write sets. R+W > N ensures strong consistency.

19. **How does Couchbase's cross-datacenter replication (XDCR) handle conflicts?**
    - **Expected**: Last-write-wins based on timestamp or vector clock. Applications can use custom conflict resolution via JavaScript functions. Eventual consistency across datacenters.

20. **What is gossip protocol used for?**
    - **Expected**: Cluster membership propagation, failure detection, metadata dissemination. Each node periodically exchanges state with a random peer. Example: Cassandra's SWIM-based gossip.

---

## How to Structure Your Answer

```
1. State the problem & clarify scope
2. Identify tradeoffs (CAP, latency vs consistency)
3. Propose a system (components, data flow)
4. Discuss failure modes
5. How you would test each failure mode
```

## Common Mistakes

| Mistake | Better Approach |
|---------|----------------|
| Diving into implementation details | First establish requirements and tradeoffs |
| Ignoring failure scenarios | Always discuss: "What happens when X fails?" |
| Not considering testing perspective | Say "I would test this by injecting faults in..." |
| Vague consistency claims | Be specific: "We use quorum reads with R=2, W=2" |
| Forgetting about clock skew | Assume clocks are not synchronized in distributed systems |

## Hints for Improvement

- **Read**: DDIA (Designing Data-Intensive Applications) by Kleppmann
- **Practice**: Explain CAP, Raft, and gossip to someone non-technical
- **Hands-on**: Run a Jepsen test locally with MongoDB/Cassandra
- **Couchbase-specific**: Understand vBucket map, DCP protocol, XDCR conflict resolution

## Difficulty Levels

| Topic | Difficulty |
|-------|-----------|
| CAP theorem basics | ★☆☆ |
| Consistency models | ★★☆ |
| Quorum math (R+W>N) | ★★☆ |
| Partitioning strategies | ★★☆ |
| Paxos / Raft | ★★★ |
| CRDTs and Vector Clocks | ★★★ |
| Jepsen-style fault injection | ★★★ |
