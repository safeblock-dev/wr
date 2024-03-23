package wrpool

import (
	"context"
	"sync/atomic"

	"github.com/safeblock-dev/wr/wrgroup"
)

// Pool manages a pool of goroutines that can execute tasks concurrently.
type Pool struct {
	group        wrgroup.WaitGroup  // Manages goroutines.
	ctx          context.Context    // Manages the lifecycle of the pool.
	limiter      limiter            // Controls the number of goroutines.
	tasks        chan func() error  // Channel for tasks to be executed by the pool.
	errorHandler func(err error)    // Handles errors during task execution.
	cancelFunc   context.CancelFunc // Cancels the context.
	groupOpts    []wrgroup.Option   // Options for the wait group.
	stopped      atomic.Bool        // Indicates whether the pool is stopped.
}

// New creates a new Pool with the provided options.
func New(options ...Option) *Pool {
	pool := &Pool{
		groupOpts: make([]wrgroup.Option, 0, 1),
		tasks:     make(chan func() error),
	}

	// Apply all options.
	for _, opt := range options {
		opt(pool)
	}

	pool.group = *wrgroup.New(pool.groupOpts...)

	// Initialize base context (if not already set).
	if pool.ctx == nil {
		Context(context.Background())(pool)
	}

	return pool
}

// Go submits a task to be run in the pool.
// If all goroutines in the pool are busy, it will block until the task can be started.
func (p *Pool) Go(f func() error) bool {
	if p.ctx.Err() != nil {
		return false
	}

	if p.limiter == nil {
		// No limit on the number of goroutines.
		select {
		case p.tasks <- f:
			// A goroutine is available to handle the task.
		default:
			// No goroutine is available; spawn a new one.
			p.group.Go(func() {
				p.worker(f)
			})
		}
	} else {
		select {
		case p.limiter <- struct{}{}:
			// Below the limit, spawn a new worker.
			p.group.Go(func() {
				p.worker(f)
			})
		case p.tasks <- f:
			// A worker is available and has accepted the task.
		case <-p.ctx.Done():
			// Context was cancelled; return without adding the task.
			return false
		}
	}

	return true
}

// Wait cleans up spawned goroutines, propagating any panics that were raised by the tasks.
func (p *Pool) Wait() {
	if p.stopped.CompareAndSwap(false, true) {
		defer p.limiter.close()
		close(p.tasks)
		p.cancelFunc()
		p.group.Wait()
	}
}

// Reset reactivates the pool, allowing new tasks to be submitted.
func (p *Pool) Reset() {
	p.Wait()
	p.stopped.Store(false)
	p.tasks = make(chan func() error)
	if p.limiter != nil {
		p.limiter = make(limiter, cap(p.limiter))
	}
	Context(context.Background())(p)
}

// IsStopped returns true if the pool is stopped.
func (p *Pool) IsStopped() bool {
	return p.stopped.Load()
}

// MaxGoroutines returns the maximum number of goroutines allowed in the pool.
func (p *Pool) MaxGoroutines() int {
	return p.limiter.limit()
}

// worker is the function run by each goroutine in the pool.
// It executes tasks and handles panics.
func (p *Pool) worker(initialFunc func() error) {
	defer p.limiter.release()

	if initialFunc != nil {
		err := initialFunc()
		if p.errorHandler != nil && err != nil {
			p.errorHandler(err)
		}
	}

	for f := range p.tasks {
		err := f()
		if p.errorHandler != nil && err != nil {
			p.errorHandler(err)
		}
	}
}
