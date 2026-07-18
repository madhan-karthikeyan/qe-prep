# Thread Pool

Difficulty: Medium
Estimated Interview Time: 35 min
Prerequisites: Threading, queues, Futures

## Problem Statement

Implement a configurable thread pool that accepts tasks via submit() and returns Future-like wrappers, with graceful shutdown.

## Requirements

- Fixed number of worker threads (configurable)
- submit(callable, *args, **kwargs) -> Future
- Graceful shutdown (waits for pending tasks, rejects new)
- Bounded task queue
- Exception propagation via Future

## Implementation Notes

- Workers run a loop pulling from an internal thread-safe Queue
- Each task is packaged as (fn, args, kwargs, Future) tuple
- shutdown() sets an Event and optionally joins all workers
- Uses concurrent.futures.Future for result/exception handling

## Test Strategy
- Unit: submit/result, multiple tasks, shutdown rejection, exception propagation, wait-for-pending
- Stress: 1000 tasks with 8 workers, verify all complete; 500 tasks with 4 workers, verify all side effects

## Edge Cases

- Submitting after shutdown raises RuntimeError
- Task that raises an exception: exception stored in Future
- Slow tasks complete even after shutdown(wait=True)
- Worker count of 0 is rejected

## Failure Cases

- num_workers < 1 (ValueError)
- Task queue full (blocking put, configurable via max_queue_size)
- Worker thread crash (not handled; would lose that worker)

## Complexity
- Time: O(1) submit, O(n) shutdown
- Space: O(max_queue_size + num_workers)

## Progression Path
Basic → Future results → Graceful shutdown → Dynamic resizing

## Common Interview Follow-ups

- How would you add dynamic resizing (grow/shrink workers)?
- How would you implement a scheduled/delayed task?
- How would you handle a worker thread that crashes?
- How would you implement work-stealing?

## Possible Production Improvements

- Dynamic worker pool that scales based on load
- Work-stealing for better load distribution
- Scheduled/delayed task execution
- Metrics (queue depth, active workers, throughput)
- Thread health monitoring and auto-restart
