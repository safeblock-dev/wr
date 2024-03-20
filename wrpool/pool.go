package wrpool

import (
	"context"
	"sync/atomic"

	"github.com/safeblock-dev/wr/wrgroup"
)

type Pool struct {
	handle        wrgroup.WaitGroup  // wr.WaitGroup to manage goroutines.
	context       context.Context    // Context for managing the lifecycle of the pool.
	limiter       limiter            // limiter for controlling the number of goroutines.
	tasks         chan func() error  // Channel for tasks to be executed by the pool.
	errorHandler  func(err error)    // Handler for errors occurring during task execution.
	contextCancel context.CancelFunc // Function to cancel the context.
	handleOptions []wrgroup.Option   // Options for the wait group.
	stopped       atomic.Bool        // Flag to indicate whether the pool is stopped.
}

// New creates a new Pool with the provided options.
func New(options ...Option) *Pool {
	var pool = new(Pool)
	pool.handleOptions = make([]wrgroup.Option, 0, 1)
	pool.tasks = make(chan func() error)

	// Apply all options.
	for _, opt := range options {
		opt(pool)
	}

	pool.handle = *wrgroup.New(pool.handleOptions...)

	// Initialize base context (if not already set).
	if pool.context == nil {
		Context(context.Background())(pool)
	}

	return pool
}

// Go submits a task to be run in the pool.
// If all goroutines in the pool are busy, a call to Go() will block until the task can be started.
func (p *Pool) Go(f func() error) {
	if p.context.Err() != nil {
		return
	}

	if p.limiter == nil {
		// No limit on the number of goroutines.
		select {
		case p.tasks <- f:
			// A goroutine was available to handle the task.
		default:
			// No goroutine was available to handle the task.
			// Spawn a new one and send it the task.
			p.handle.Go(func() {
				p.worker(f)
			})
		}
	} else {
		select {
		case p.limiter <- struct{}{}:
			// If we are below our limit, spawn a new worker rather
			// than waiting for one to become available.
			p.handle.Go(func() {
				p.worker(f)
			})
		case p.tasks <- f:
			// A worker is available and has accepted the task.
			return
		case <-p.context.Done():
			// The context was cancelled, return without adding the task.
			return
		}
	}
}

// Wait cleans up spawned goroutines, propagating any panics that were raised by the tasks.
func (p *Pool) Wait() {
	if p.stopped.CompareAndSwap(false, true) {
		defer p.limiter.close()

		close(p.tasks)
		p.contextCancel()
		p.handle.Wait()
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

func (p *Pool) IsStopped() bool {
	return p.stopped.Load()
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
