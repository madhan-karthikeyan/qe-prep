# Thread Pool
Difficulty: Medium
Estimated Interview Time: 25 min
Prerequisites: goroutines, channels, sync.WaitGroup

## Problem Statement
Implement a goroutine pool (worker pool) that manages a fixed number of workers executing submitted tasks.

## Requirements
- Fixed worker count
- Submit(func()) — non-blocking enqueue
- Wait() — block until all tasks complete
- Stop() — graceful shutdown, no new tasks accepted
- SubmitAndWait convenience method

## Implementation Notes
- Workers are goroutines pulling from a buffered channel
- Submit is non-blocking; returns false if queue is full or pool stopped
- Stop closes the task channel (workers drain and exit)
- Wait uses sync.WaitGroup to track completion
- Stop is safe to call multiple times via sync.Once

## Test Strategy
- Submit and verify all tasks complete
- Results via callback capture
- Graceful shutdown with pending tasks
- No tasks edge case
- Stress test: 1000 tasks, 8 workers, -race

## Edge Cases
- Submit after Stop returns false
- Stop called multiple times
- Zero tasks submitted
- Buffer smaller than number of tasks

## Failure Cases
- N/A (Submit returns bool for failure)

## Complexity (Time + Space)
- Submit: O(1) amortized
- Wait: O(1) signaling
- Space: O(bufferSize + numWorkers goroutines)

## Progression Path (Basic → Intermediate → Advanced → Production)
- Basic: Fixed worker pool with no return values
- Intermediate: Task return values via channels
- Advanced: Dynamic worker scaling, priority scheduling
- Production: Panic recovery, metrics, dynamic scaling

## Common Interview Follow-ups
- How would you return results from tasks?
- How would you handle panics in tasks?
- How would you implement dynamic worker scaling?

## Possible Production Improvements
- Return typed results via channels
- Panic recovery in workers
- Dynamic worker scaling based on queue depth
- Task prioritization
- Metrics (tasks completed, queue depth, worker utilization)
