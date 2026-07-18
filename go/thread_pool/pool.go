package thread_pool

import (
	"sync"
)

// Task is a unit of work that the pool executes.
type Task func()

// Pool is a fixed-size goroutine pool. Tasks are submitted via a buffered
// channel and executed by worker goroutines.
type Pool struct {
	tasks    chan Task
	wg       sync.WaitGroup
	stopOnce sync.Once
	stopped  chan struct{}
}

// New creates a new Pool with the given number of workers and buffer size for
// the task queue.
func New(numWorkers, bufferSize int) *Pool {
	if numWorkers < 1 {
		numWorkers = 1
	}
	if bufferSize < 1 {
		bufferSize = 1
	}
	p := &Pool{
		tasks:   make(chan Task, bufferSize),
		stopped: make(chan struct{}),
	}
	for range numWorkers {
		p.wg.Add(1)
		go p.worker()
	}
	return p
}

// worker pulls tasks from the channel and executes them until the channel is
// closed.
func (p *Pool) worker() {
	defer p.wg.Done()
	for task := range p.tasks {
		task()
	}
}

// Submit enqueues a task for execution. It never blocks; if the queue is full
// or the pool is stopped, the task is discarded and false is returned.
func (p *Pool) Submit(task Task) bool {
	select {
	case <-p.stopped:
		return false
	default:
	}
	select {
	case p.tasks <- task:
		return true
	case <-p.stopped:
		return false
	default:
		return false
	}
}

// Wait blocks until all submitted tasks have completed execution.
func (p *Pool) Wait() {
	p.wg.Wait()
}

// Stop initiates a graceful shutdown. No new tasks are accepted after Stop
// returns. Wait must be called separately to wait for completion.
func (p *Pool) Stop() {
	p.stopOnce.Do(func() {
		close(p.stopped)
		close(p.tasks)
	})
}

// SubmitAndWait is a convenience method that submits a task and blocks until it
// completes. Returns false if the task could not be submitted.
func (p *Pool) SubmitAndWait(task Task) bool {
	done := make(chan struct{})
	ok := p.Submit(func() {
		task()
		close(done)
	})
	if !ok {
		return false
	}
	<-done
	return true
}
