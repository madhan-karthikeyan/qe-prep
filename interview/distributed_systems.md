# Distributed Systems — 20 Interview Questions

## Overview

For QE roles at companies building distributed storage, databases, or infrastructure, these questions test your understanding of distributed systems fundamentals and how to validate them.

## 20 Questions with Expected Answers

### 1. What's the difference between consistency and availability in the CAP theorem?

**Expected**: Consistency means every read receives the most recent write (or an error). Availability means every request receives a (non-error) response — without guarantee it contains the most recent write. During a network partition, you must choose: block until consistent (CP) or return possibly stale data (AP).

**Mistake**: Saying "you pick 2 of 3" — CAP only applies during partitions. When the network is healthy, you can have both.

### 2. Explain partition tolerance. Is it optional?

**Expected**: Partition tolerance means the system continues to operate despite dropped or delayed messages between nodes. In distributed systems, partitions are inevitable — so you must tolerate them. You never choose "no partition tolerance" in practice.

### 3. How does a gossip protocol work?

**Expected**: Each node periodically picks a random peer and exchanges state information (membership, metadata, failure detection). Convergence is O(log N) rounds. Example: SWIM protocol used in Cassandra.

**Testing angle**: Test convergence time after a node join/failure. Inject message loss and verify eventual convergence.

### 4. What is hinted handoff?

**Expected**: When a write target node is down, another node temporarily stores the write (with a hint about the intended recipient). When the target recovers, the hint is replayed. Improves write availability.

**Testing angle**: Kill a node, write data, restart node, verify data arrives.

### 5. Explain vector clocks and how they detect conflicts.

**Expected**: Each node maintains a vector of (node → version). On write, increment own counter. On read, return vector. Concurrent writes produce vectors that are incomparable (neither dominates the other), indicating a conflict.

**Testing angle**: Generate concurrent writes across nodes and verify conflict detection fires. Verify conflict resolution (LWW, CRDT merge, custom).

### 6. What are CRDTs? Give an example.

**Expected**: Conflict-free Replicated Data Types converge automatically without coordination. Examples: G-Counter (grow-only counter uses max), OR-Set (add/remove set with observed-remove semantics). Operations commute.

**Testing angle**: Test convergence after arbitrary reordering of operations. Verify idempotency.

### 7. How does consistent hashing work? Why use virtual nodes?

**Expected**: Nodes and keys map to a hash ring. Each key is assigned to the next clockwise node. Virtual nodes distribute each physical node across multiple ring positions, improving load balance and reducing reshuffle on node changes.

**Testing angle**: Add/remove nodes and measure % of keys relocated. Verify load distribution.

### 8. Compare quorum-based and leader-based replication.

**Expected**: 
- **Leader-based**: Single leader accepts writes, replicates to followers. Simple but single point of failure.
- **Quorum-based**: Any replica can accept writes, but requires R+W > N for consistency. Higher availability, more complex conflict resolution.

### 9. How do read repair and anti-entropy differ?

**Expected**: Read repair corrects stale data encountered during reads (lazy, on-access). Anti-entropy is a background process that compares Merkle trees and syncs differences (proactive). Both improve eventual consistency.

### 10. What's the difference between at-least-once and exactly-once delivery?

**Expected**: At-least-once guarantees delivery but allows duplicates (receiver must deduplicate). Exactly-once requires deduplication via idempotent operations or transaction logs. Exactly-once is expensive — Kafka achieves it through idempotent producers + transactional semantics.

**Testing angle**: Deliver the same message twice; verify deduplication or at-least-once behavior.

### 11. How does the Raft leader election work?

**Expected**: Nodes are followers, candidates, or leaders. Followers expect periodic heartbeats. On timeout, follower becomes candidate, votes for itself, requests votes from others. Receives majority → becomes leader. Randomized election timeouts prevent split votes.

### 12. What is a split-brain scenario? Name two prevention mechanisms.

**Expected**: Two nodes believe they're both active leaders. Prevention: (1) Quorum-based fencing — only the node with majority can act. (2) Lease/stale read — leader must renew lease, nodes reject expired leases.

### 13. Explain the balloon effect / fencing in distributed systems.

**Expected**: When a partition heals, nodes that missed updates must "catch up." If they receive too many writes at once, they can become overwhelmed (balloon effect). Fencing tokens prevent stale leaders from making writes — every write request includes a monotonically increasing token; older tokens are rejected.

### 14. How do you test that a distributed system is "eventually consistent"?

**Expected**: 
1. Write a known value to one node
2. Partition the network
3. Write another value to a different node
4. Heal the partition
5. Poll all nodes until values converge
6. Assert that all nodes have the same value within a timeout

Add property-based checks: no lost updates, monotonic reads, read-your-writes.

### 15. Design a test for multi-leader replication conflict resolution.

**Expected**: 
1. Set up two leaders in different DCs
2. Write conflicting keys simultaneously on both (use a barrier for timing)
3. Allow replication to propagate
4. Assert that conflict resolution triggers (LWW, CRDT merge, or custom handler)
5. Verify no data loss and resolution is deterministic

### 16. What is the testing challenge with clock skew?

**Expected**: Last-write-wins (LWW) depends on timestamps. If clocks are skewed, an older write can win over a newer one. Solution: use vector clocks or hybrid logical clocks (HLC). Test by artificially skewing clocks and verifying correctness.

### 17. How do you test failure detection in a gossip-based system?

**Expected**: 
1. Kill a node (SIGKILL, not graceful)
2. Wait for gossip rounds
3. Verify remaining nodes mark it as dead
4. Measure detection time
5. Test false positives: temporarily drop packets to a live node, verify it isn't falsely declared dead

### 18. What is the write amplification problem in LSM trees?

**Expected**: LSM trees constantly compact sorted runs, rewriting data. Write amplification = total bytes written / new bytes ingested. High amplification reduces SSD lifespan. Test by measuring disk I/O during sustained writes.

### 19. How would you verify a distributed transaction protocol (2PC)?

**Expected**: 
1. Send prepare to all participants
2. Kill coordinator after all prepared but before commit
3. Verify participants remain in prepared state (blocking)
4. Recover coordinator; verify transaction commits
5. Kill a participant during commit; verify coordinator retries
6. Test timeout behavior and transaction log recovery

### 20. Compare Dynamo-style (AP) vs Bigtable-style (CP) systems from a testing perspective.

| Aspect | Dynamo-style (Cassandra) | Bigtable-style (HBase) |
|--------|--------------------------|----------------------|
| Guarantee | AP, eventual consistency | CP, strong consistency |
| Test focus | Conflict resolution, convergence, read repair | Partition handling, availability tradeoff, failover |
| Key tests | Concurrent writes, cluster heal, hinted handoff | ZK failover, region server crash, WAL replay |

---

## Difficulty Levels

| Topic | Difficulty |
|-------|-----------|
| CAP theorem basics | ★☆☆ |
| Consistent hashing | ★★☆ |
| Hinted handoff / read repair | ★★☆ |
| Vector clocks | ★★★ |
| CRDTs | ★★★ |
| Raft / Paxos | ★★★ |
| Distributed transaction testing | ★★★ |

## Resources

- **Designing Data-Intensive Applications** (Kleppmann) — Chapters 5-9
- **Jepsen** blog — aphyr.com — distributed system failure analyses
- **Couchbase** architecture docs — vBucket, DCP, XDCR
- **Testing Distributed Systems** — curated list at github.com/asatarin/testing-distributed-systems
