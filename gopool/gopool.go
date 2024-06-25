package gopool

import (
	"context"
	"sync/atomic"

	"github.com/safeblock-dev/wr/syncgroup"
)

// Pool manages a pool of goroutines that can execute tasks concurrently.
type Pool struct {
	ctx          context.Context      // ctx is the pool's context.
	parentCtx    context.Context      // parentCtx is the parent context of the pool.
	group        *syncgroup.WaitGroup // group is the wait group managing goroutines.
	cancelFunc   context.CancelFunc   // cancelFunc cancels the pool's context.
	limiter      limiter              // limiter controls the number of concurrent goroutines.
	tasks        chan func() error    // tasks is a channel for task functions.
	errorHandler func(err error)      // errorHandler handles errors encountered during task execution.
	panicHandler func(pc any)         // panicHandler handles panics recovered during task execution.
	stopped      atomic.Bool          // stopped indicates if the pool has been stopped.
}

// New creates a new Pool with the provided options.
func New(options ...Option) *Pool {
	pool := &Pool{ //nolint: exhaustruct
		panicHandler: defaultPanicHandler, // Set default panic handler.
		tasks:        make(chan func() error),
	}

	// Apply all options.
	for _, opt := range options {
		opt(pool)
	}

	// Initialize base context (if not already set).
	if pool.ctx == nil {
		Context(context.Background())(pool)
	}

	// Initialize wait group with panic handler.
	pool.group = syncgroup.New(syncgroup.PanicHandler(pool.panicHandler))

	return pool
}

// Go submits a task to be run in the pool. If all goroutines in the pool
// are busy, a call to Go() will block until the task can be started.
// Note: If this function is called after Wait(), it will cause a panic.
func (p *Pool) Go(f func() error) {
	if p.ctx.Err() != nil {
		return // Return if the pool's context is canceled.
	}

	if p.limiter == nil {
		// No limit on the number of goroutines.
		select {
		case p.tasks <- f:
			// A goroutine is available to handle the task.
		default:
			// No goroutine was available to handle the task.
			// Spawn a new one and send it the task.
			p.group.Go(p.worker)
			// We know there is at least one worker running, so wait
			// for it to become available. This ensures we never spawn
			// more workers than the number of tasks.
			p.tasks <- f
		}
	} else {
		select {
		case p.limiter <- struct{}{}:
			// If we are below our limit, spawn a new worker rather
			// than waiting for one to become available.
			p.group.Go(p.worker)
			p.tasks <- f
		case <-p.ctx.Done():
			// Context was cancelled; return without adding the task.
		case p.tasks <- f:
			// A worker is available and has accepted the task.
		}
	}
}

// Wait cleans up spawned goroutines, propagating any panics that were raised by the tasks.
func (p *Pool) Wait() {
	if p.stopped.CompareAndSwap(false, true) {
		p.cancelFunc()
		close(p.tasks)
		p.group.Wait()
		p.limiter.close()
	}
}

// Cancel cancels the pool's context.
func (p *Pool) Cancel() {
	p.cancelFunc()
}

// Reset reactivates the pool, allowing new tasks to be submitted.
func (p *Pool) Reset() {
	p.Wait()
	p.tasks = make(chan func() error)
	if p.limiter != nil {
		p.limiter = make(limiter, p.limiter.limit())
	}
	Context(p.parentCtx)(p)
	p.stopped.Store(false)
}

// worker is the function run by each goroutine in the pool.
// It executes tasks and handles panics.
func (p *Pool) worker() {
	defer p.limiter.release() // Release limiter when worker exits.

	for f := range p.tasks {
		err := f()
		if p.errorHandler != nil && err != nil {
			p.errorHandler(err)
		}
	}
}
