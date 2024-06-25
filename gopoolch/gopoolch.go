package gopoolch

import (
	"sync"
	"sync/atomic"

	"github.com/safeblock-dev/werr"
	"github.com/safeblock-dev/wr/gopool"
)

// PoolCh is a wrapper around gopool.Pool with additional error handling.
type PoolCh struct {
	pool    *gopool.Pool // pool is the underlying gopool.Pool instance.
	err     error        // err stores the first error encountered in the pool.
	errCh   chan error   // errCh is a channel for propagating errors from the pool.
	errOnce sync.Once    // errOnce ensures the error is only set once.
	stopped atomic.Bool  // stopped indicates if the pool has been stopped.
}

// New creates a new PoolCh with the provided options.
// It initializes the pool with additional error and panic handling.
func New(options ...gopool.Option) *PoolCh {
	var p PoolCh
	p.errCh = make(chan error, 1)

	// Combine default error handling options with user-provided options
	allOptions := append([]gopool.Option{
		gopool.ErrorHandler(p.errorHandler),
		gopool.PanicHandler(p.panicHandler),
	}, options...)

	// Initialize the pool with the combined options
	p.pool = gopool.New(allOptions...)

	return &p
}

// Go submits a task to the pool for execution.
func (p *PoolCh) Go(f func() error) {
	if !p.stopped.Load() {
		p.pool.Go(f)
	}
}

// Wait waits for all tasks in the pool to complete and closes the error channel.
func (p *PoolCh) Wait() {
	if p.stopped.CompareAndSwap(false, true) {
		p.pool.Wait()
		close(p.errCh)
	}
}

// Reset reactivates the pool, allowing new tasks to be submitted.
func (p *PoolCh) Reset() {
	p.Wait()
	p.pool.Reset()
	p.stopped.Store(false)
	p.err = nil
	p.errCh = make(chan error, 1)
	p.errOnce = sync.Once{}
}

// ErrorChannel returns a channel that can be used to receive errors that occur in the pool.
func (p *PoolCh) ErrorChannel() <-chan error {
	return p.errCh
}

// Error returns the first error that occurred in the pool.
func (p *PoolCh) Error() error {
	return p.err
}

// HasError returns true if an error has occurred in the pool.
func (p *PoolCh) HasError() bool {
	return p.err != nil
}

// panicHandler handles panics that occur in the pool by converting
// them to errors and passing them to the error handler.
func (p *PoolCh) panicHandler(pc interface{}) {
	p.errorHandler(werr.PanicToError(pc))
}

// errorHandler handles errors that occur in the pool by sending
// them to the error channel and setting the error.
func (p *PoolCh) errorHandler(err error) {
	p.errOnce.Do(func() {
		p.err = err
		p.errCh <- err
		p.pool.Cancel()
	})
}
