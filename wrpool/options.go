package wrpool

import (
	"context"

	"github.com/safeblock-dev/wr/panics"
	"github.com/safeblock-dev/wr/wrgroup"
)

// Option represents an option that can be passed when instantiating a worker pool to customize it.
type Option func(pool *Pool)

// PanicHandler allows to change the panic handler function of a pool.
func PanicHandler(panicHandler func(recovered panics.Recovered)) Option {
	return func(pool *Pool) {
		pool.handleOptions = append(pool.handleOptions, wrgroup.PanicHandler(panicHandler))
	}
}

// ErrorHandler allows to change the error handler function of a pool.
func ErrorHandler(errorHandler func(err error)) Option {
	return func(pool *Pool) {
		pool.errorHandler = errorHandler
	}
}

// Context configures a parent context on a worker pool to stop all workers when it is cancelled.
func Context(parentCtx context.Context) Option {
	return func(pool *Pool) {
		pool.context, pool.contextCancel = context.WithCancel(parentCtx)
	}
}

func MaxGoroutines(limit int) Option {
	return func(pool *Pool) {
		pool.limiter = make(limiter, limit)
	}
}
