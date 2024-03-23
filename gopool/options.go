package gopool

import (
	"context"

	"github.com/safeblock-dev/wr/panics"
	"github.com/safeblock-dev/wr/syncgroup"
)

// Option represents an option that can be passed when instantiating a Pool to customize it.
type Option func(pool *Pool)

// PanicHandler sets the panic handler function for the pool.
func PanicHandler(panicHandler func(recovered panics.Recovered)) Option {
	return func(pool *Pool) {
		pool.groupOpts = append(pool.groupOpts, syncgroup.PanicHandler(panicHandler))
	}
}

// ErrorHandler sets the error handler function for the pool.
func ErrorHandler(errorHandler func(err error)) Option {
	return func(pool *Pool) {
		pool.errorHandler = errorHandler
	}
}

// Context sets a parent context for the pool to stop all workers when it is cancelled.
func Context(ctx context.Context) Option {
	return func(pool *Pool) {
		pool.ctx, pool.cancelFunc = context.WithCancel(ctx)
	}
}

// MaxGoroutines sets the maximum number of goroutines allowed in the pool.
func MaxGoroutines(limit int) Option {
	return func(pool *Pool) {
		pool.limiter = make(limiter, limit)
	}
}
